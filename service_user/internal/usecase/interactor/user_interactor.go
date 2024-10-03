package interactor

import (
	"context"
	"encoding/json"
	"fmt"
	"service_user/internal/apperrors"
	"service_user/internal/domain/mappers"
	"service_user/internal/service_user"

	"service_user/internal/usecase/cache"
	"service_user/internal/usecase/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

const timeCacheExpiration = 60

type userInteractor struct {
	UserRepository repository.UserRepository
	UserCache      cache.UserCache
}

type UserInteractor interface {
	CreateUser(ctx context.Context, createUserRequest *service_user.CreateUserRequest) (string, error)
	GetUsersByPageAndPerPage(ctx context.Context, page, perPage string) (*service_user.GetUsersByPaginationResponse, error)
	GetUserByID(ctx context.Context, id string) (*service_user.GetUserByIDResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*service_user.GetUserByEmailResponse, error)
	UpdateUserByID(ctx context.Context, updateUserRequest *service_user.UpdateUserByIDRequest) (string, error)
	DropUserByID(ctx context.Context, dropUserReq *service_user.DropUserByIDRequest) error
}

func NewUserInteractor(r repository.UserRepository, h cache.UserCache) UserInteractor {
	return &userInteractor{r, h}
}

func (us *userInteractor) GetUserByEmail(ctx context.Context, email string) (*service_user.GetUserByEmailResponse, error) {
	cachedUser, err := us.UserCache.Get(ctx, email)
	if err != nil && apperrors.IsAppError(err, &apperrors.RedisKeyDoesntExist) {
		return us.cashedUserByEmail(ctx, email)
	}

	if err != nil {
		return nil, err
	}

	getUsersByIdResp, err := mappers.MapCachedUserToGetUsersByEmailResponse(cachedUser)
	if err != nil {
		log.Errorf("Cannot map user to GetUsersByEmailResponse. ID: %s; Err: %+v", email, err)
		return nil, err
	}

	return getUsersByIdResp, nil
}

func (us *userInteractor) CreateUser(ctx context.Context, createUserRequest *service_user.CreateUserRequest) (string, error) {
	user, err := mappers.MapCreateUserRequestToUser(createUserRequest)
	if err != nil {
		return "", err
	}

	hashPass, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}

	user.Password = hashPass
	return us.UserRepository.CreateUser(ctx, user)
}

func (us *userInteractor) GetUsersByPageAndPerPage(ctx context.Context, pageIn, perPageIn string) (*service_user.GetUsersByPaginationResponse, error) {
	key := fmt.Sprintf("%v_%v", pageIn, perPageIn)
	cachedUsersByPageAndPerPage, err := us.UserCache.Get(ctx, key)
	if err != nil && apperrors.IsAppError(err, &apperrors.RedisKeyDoesntExist) {
		return us.cacheUsersByPageAndPerPage(ctx, pageIn, perPageIn)
	}

	if err != nil {
		return nil, err
	}

	getUsersByPaginationResponse, err := mappers.MapCachedUserByPageAndPerPageToGetUsersByPaginationResponse(cachedUsersByPageAndPerPage)
	if err != nil {
		log.Errorf("Cannot map users to GetUsersByPaginationResponse. Key: %s; Err: %+v", key, err)
		return nil, err
	}

	return getUsersByPaginationResponse, nil
}

func (us *userInteractor) cacheUsersByPageAndPerPage(ctx context.Context, pageIn, perPageIn string) (*service_user.GetUsersByPaginationResponse, error) {
	page, err := strconv.Atoi(pageIn)
	if err != nil {
		msg := fmt.Sprintf("Cannot parse page to int. Err: %v", err)
		return nil, apperrors.ParsingStringToIntFailed.AppendMessage(msg)
	}

	perPage, err := strconv.Atoi(perPageIn)
	if err != nil {
		msg := fmt.Sprintf("Cannot parse per page to int. Err: %v", err)
		return nil, apperrors.ParsingStringToIntFailed.AppendMessage(msg)
	}

	users, err := us.UserRepository.GetUsersByPageAndPerPage(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	respUsers := mappers.MapUsersToGetUsersByPaginationResponse(users, pageIn, perPageIn)
	data, err := json.Marshal(respUsers)
	if err != nil {
		appErr := apperrors.UnmarshalError.AppendMessage(err)
		return nil, appErr
	}

	key := fmt.Sprintf("%v_%v", pageIn, perPageIn)
	expiration := time.Duration(timeCacheExpiration) * time.Second
	err = us.UserCache.SetWithExpiration(ctx, key, data, expiration)
	if err != nil {
		return nil, err
	}

	return respUsers, nil
}

func (us *userInteractor) GetUserByID(ctx context.Context, id string) (*service_user.GetUserByIDResponse, error) {
	cachedUser, err := us.UserCache.Get(ctx, id)
	if err != nil && apperrors.IsAppError(err, &apperrors.RedisKeyDoesntExist) {
		return us.cashedUserByID(ctx, id)
	}

	if err != nil {
		return nil, err
	}

	getUsersByIdResp, err := mappers.MapCachedUserToGetUsersByIdResponse(cachedUser)
	if err != nil {
		log.Errorf("Cannot map user to GetUsersByIdResponse. ID: %s; Err: %+v", id, err)
		return nil, err
	}

	return getUsersByIdResp, nil
}

func (us *userInteractor) cashedUserByEmail(ctx context.Context, email string) (*service_user.GetUserByEmailResponse, error) {
	user, err := us.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	getUsersByEmailResp := mappers.MapUserToGetUserByEmailResponse(user)
	data, err := json.Marshal(getUsersByEmailResp)
	if err != nil {
		appErr := apperrors.UnmarshalError.AppendMessage(err)
		return nil, appErr
	}

	expiration := time.Duration(timeCacheExpiration) * time.Second
	err = us.UserCache.SetWithExpiration(ctx, email, data, expiration)
	if err != nil {
		return nil, err
	}

	return getUsersByEmailResp, nil
}

func (us *userInteractor) cashedUserByID(ctx context.Context, id string) (*service_user.GetUserByIDResponse, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		msg := fmt.Sprintf("Cannot convert UserID to UUID. Err: %v", err)
		return nil, apperrors.CoverstionToUUIDFailed.AppendMessage(msg)
	}

	user, err := us.UserRepository.GetUserByID(ctx, &userUUID)
	if err != nil {
		return nil, err
	}

	log.Info("RATING", user)
	getUsersByIdResp := mappers.MapUserToGetUserByIdResponse(user)
	data, err := json.Marshal(getUsersByIdResp)
	if err != nil {
		appErr := apperrors.UnmarshalError.AppendMessage(err)
		return nil, appErr
	}

	expiration := time.Duration(timeCacheExpiration) * time.Second
	err = us.UserCache.SetWithExpiration(ctx, id, data, expiration)
	if err != nil {
		return nil, err
	}

	return getUsersByIdResp, nil
}

func (us *userInteractor) UpdateUserByID(ctx context.Context, updateUserRequest *service_user.UpdateUserByIDRequest) (string, error) {
	user, err := mappers.MapUpdateUserRequestToUser(updateUserRequest)
	if err != nil {
		return "", err
	}

	return us.UserRepository.UpdateUserByID(ctx, user)
}

func (us *userInteractor) DropUserByID(ctx context.Context, dropUserByIdReq *service_user.DropUserByIDRequest) error {
	userUUID, err := uuid.Parse(dropUserByIdReq.Id)
	if err != nil {
		return apperrors.CoverstionToUUIDFailedDropUserByID.AppendMessage(err)
	}

	return us.UserRepository.DropUserByID(ctx, &userUUID)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", apperrors.HashPasswordErr.AppendMessage(err)
	}

	return string(bytes), nil
}
