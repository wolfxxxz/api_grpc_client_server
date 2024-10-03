package mappers

import (
	"encoding/json"
	"service_user/internal/apperrors"
	"service_user/internal/domain/models"
	"service_user/internal/service_user"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCreateUserRequestToUser(createUserRequest *service_user.CreateUserRequest) (*models.User, error) {
	if createUserRequest.FirstName == "" {
		return nil, &apperrors.MapCreateUserRequestToUserErr
	}

	user := MapToUser(createUserRequest.Email, createUserRequest.UserName, createUserRequest.FirstName,
		createUserRequest.LastName, createUserRequest.Password, createUserRequest.Role)
	return user, nil
}

func MapUsersToGetUsersByPaginationResponse(users []*models.User, page, perPage string) *service_user.GetUsersByPaginationResponse {
	var totalUsers int32 = int32(len(users))
	var usersResponse []*service_user.User
	for _, user := range users {
		var userResponse service_user.User
		userResponse.Email = user.Email
		userResponse.FirstName = user.FirstName
		userResponse.Id = user.ID.String()
		userResponse.LastName = user.LastName
		userResponse.UserName = user.UserName
		userResponse.Rating = user.Rating()
		usersResponse = append(usersResponse, &userResponse)
	}

	return &service_user.GetUsersByPaginationResponse{Users: usersResponse, Page: page, PerPage: perPage, TotalUsers: totalUsers}
}

func MapUserToGetUserByIdResponse(user *models.User) *service_user.GetUserByIDResponse {
	rating := user.Rating()
	return &service_user.GetUserByIDResponse{Id: user.ID.String(), Email: user.Email, UserName: user.UserName,
		LastName: user.LastName, FirstName: user.FirstName, Rating: rating}
}

func MapUserToGetUserByEmailResponse(user *models.User) *service_user.GetUserByEmailResponse {
	rating := user.Rating()
	return &service_user.GetUserByEmailResponse{Id: user.ID.String(), Email: user.Email, UserName: user.UserName,
		LastName: user.LastName, FirstName: user.FirstName, Rating: rating}
}

func MapUpdateUserRequestToUser(requestUser *service_user.UpdateUserByIDRequest) (*models.User, error) {
	now := primitive.NewDateTimeFromTime(time.Now())
	userUUID, err := uuid.Parse(requestUser.Id)
	if err != nil {
		return nil, &apperrors.CoverstionToUUIDFailedUpdateUserByID
	}

	return &models.User{ID: &userUUID, Email: requestUser.Email, UserName: requestUser.UserName,
		FirstName: requestUser.FirstName, LastName: requestUser.LastName, Updated: &now}, nil
}

func MapToUser(email, userName, firstName, lastName, password, role string) *models.User {
	var votes []*models.Vote
	id := uuid.New()
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.User{ID: &id, Email: email, UserName: userName,
		FirstName: firstName, LastName: lastName, Role: role,
		Password: password, Created: &now, Votes: votes}
}

func MapCachedUserToGetUsersByIdResponse(cachedUser string) (*service_user.GetUserByIDResponse, error) {
	var getUsersByIdResp service_user.GetUserByIDResponse
	err := json.Unmarshal([]byte(cachedUser), &getUsersByIdResp)
	if err != nil {
		return nil, apperrors.UnmarshalError.AppendMessage(err)
	}

	return &getUsersByIdResp, nil
}

func MapCachedUserToGetUsersByEmailResponse(cachedUser string) (*service_user.GetUserByEmailResponse, error) {
	var getUsersByEmailResp service_user.GetUserByEmailResponse
	err := json.Unmarshal([]byte(cachedUser), &getUsersByEmailResp)
	if err != nil {
		return nil, apperrors.UnmarshalError.AppendMessage(err)
	}

	return &getUsersByEmailResp, nil
}

func MapCachedUserByPageAndPerPageToGetUsersByPaginationResponse(cachedUsersByPageAndPerPage string) (*service_user.GetUsersByPaginationResponse, error) {
	var getUsersByPaginationResponse service_user.GetUsersByPaginationResponse
	err := json.Unmarshal([]byte(cachedUsersByPageAndPerPage), &getUsersByPaginationResponse)
	if err != nil {
		return nil, apperrors.UnmarshalError.AppendMessage(err)
	}

	return &getUsersByPaginationResponse, nil
}
