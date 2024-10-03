package controller

import (
	context "context"
	"service_user/internal/apperrors"
	"service_user/internal/config"
	"service_user/internal/domain/mappers"
	"service_user/internal/domain/models"
	"service_user/internal/mock"
	"service_user/internal/service_user"
	"service_user/internal/usecase/interactor"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	req := &service_user.CreateUserRequest{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Role:      "user",
	}

	log := logrus.New()
	user, err := mappers.MapCreateUserRequestToUser(req)
	if err != nil {
		log.Error(err)
	}

	UUID, err := uuid.Parse("0ae29b28-b33e-4590-955d-74e205b80e93")
	if err != nil {
		log.Error(err)
	}

	user.ID = &UUID

	testTable := []struct {
		scenario      string
		inputUserReq  *service_user.CreateUserRequest
		user          *models.User
		response      interface{}
		expectedError error
	}{
		{
			"create user positive",
			req,
			user,
			"0ae29b28-b33e-4590-955d-74e205b80e93",
			nil,
		},
		{
			"create user negative",
			req,
			user,
			"Map_Create_User_Request_To_User_Err: Create User Err",
			&apperrors.MapCreateUserRequestToUserErr,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mock.NewMockUserRepository(ctrl)
			userHacheMock := mock.NewMockUserCache(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, userHacheMock)
			config := config.Config{}
			uController := NewUserController(uInteractor, log, &config)
			ctx := context.Background()
			userRepoMock.EXPECT().CreateUser(ctx, gomock.Any()).Return(tc.user.ID.String(), tc.expectedError).AnyTimes()
			resp, err := uController.CreateUser(ctx, tc.inputUserReq)
			if err != nil {
				apperrors.IsAppError(err, tc.expectedError.(*apperrors.AppError))
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, resp.UserId, tc.response)
		})
	}
}

func TestGetUserById(t *testing.T) {
	req := &service_user.GetUserByIDRequest{
		Id: "a700ce3c-e29a-4523-af3b-ff75d8ab40a5",
	}

	reqCreateUser := &service_user.CreateUserRequest{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Role:      "user",
	}

	log := logrus.New()
	user, err := mappers.MapCreateUserRequestToUser(reqCreateUser)
	if err != nil {
		log.Error(err)
	}

	UUID, err := uuid.Parse("a700ce3c-e29a-4523-af3b-ff75d8ab40a5")
	user.ID = &UUID
	if err != nil {
		log.Error(err)
	}

	user.ID = &UUID

	testTable := []struct {
		scenario         string
		inputUserRequest *service_user.GetUserByIDRequest
		expectedUser     *models.User
		response         interface{}
		expectedCacheErr error
		expectedError    error
	}{
		{
			"user not found by id",
			req,
			user,
			"USER_REPO_ERR: Failed GetUserByID",
			&apperrors.RedisKeyDoesntExist,
			&apperrors.MongoGetUserByIDErr,
		},
		{
			"get one user by id positive",
			req,
			user,
			*mappers.MapUserToGetUserByIdResponse(user),
			&apperrors.RedisKeyDoesntExist,
			nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mock.NewMockUserRepository(ctrl)
			userHacheMock := mock.NewMockUserCache(ctrl)

			uInteractor := interactor.NewUserInteractor(userRepoMock, userHacheMock)
			config := config.Config{}
			log := logrus.New()
			uController := NewUserController(uInteractor, log, &config)
			ctx := context.Background()
			userHacheMock.EXPECT().Get(ctx, gomock.Any()).Return("", tc.expectedCacheErr).AnyTimes()
			userHacheMock.EXPECT().SetWithExpiration(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			userRepoMock.EXPECT().GetUserByID(ctx, gomock.Any()).Return(tc.expectedUser, tc.expectedError).AnyTimes()

			getUserByIDResp, err := uController.GetUserById(ctx, tc.inputUserRequest)
			if err != nil {
				apperrors.IsAppError(err, tc.expectedError.(*apperrors.AppError))
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, getUserByIDResp.Id, tc.expectedUser.ID.String())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	req := &service_user.GetUserByEmailRequest{
		Email: "test@example.com",
	}

	reqCreateUser := &service_user.CreateUserRequest{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Role:      "user",
	}

	log := logrus.New()
	user, err := mappers.MapCreateUserRequestToUser(reqCreateUser)
	if err != nil {
		log.Error(err)
	}

	UUID, err := uuid.Parse("a700ce3c-e29a-4523-af3b-ff75d8ab40a5")
	user.ID = &UUID
	if err != nil {
		log.Error(err)
	}

	user.ID = &UUID

	testTable := []struct {
		scenario         string
		inputUserRequest *service_user.GetUserByEmailRequest
		expectedUser     *models.User
		response         interface{}
		expectedCacheErr error
		expectedError    error
	}{
		{
			"user not found by id",
			req,
			user,
			"USER_REPO_ERR: Failed GetUserByID",
			&apperrors.RedisKeyDoesntExist,
			&apperrors.MongoGetUserByIDErr,
		},
		{
			"get one user by id positive",
			req,
			user,
			*mappers.MapUserToGetUserByIdResponse(user),
			&apperrors.RedisKeyDoesntExist,
			nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mock.NewMockUserRepository(ctrl)
			userHacheMock := mock.NewMockUserCache(ctrl)

			uInteractor := interactor.NewUserInteractor(userRepoMock, userHacheMock)
			config := config.Config{}
			log := logrus.New()
			uController := NewUserController(uInteractor, log, &config)
			ctx := context.Background()
			userHacheMock.EXPECT().Get(ctx, gomock.Any()).Return("", tc.expectedCacheErr).AnyTimes()
			userHacheMock.EXPECT().SetWithExpiration(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			userRepoMock.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(tc.expectedUser, tc.expectedError).AnyTimes()

			getUserByIDResp, err := uController.GetUserByEmail(ctx, tc.inputUserRequest)
			if err != nil {
				apperrors.IsAppError(err, tc.expectedError.(*apperrors.AppError))
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, getUserByIDResp.Id, tc.expectedUser.ID.String())
		})
	}
}

func TestGetUsersByPagination(t *testing.T) {
	user := mappers.MapToUser("hello@hello.com", "John_Snow", "John", "Targaryen", "password", "user")
	userTwo := mappers.MapToUser("hello@hello.com", "John_Snow", "John", "Targaryen", "password", "user")
	users := []*models.User{user, userTwo}
	req := &service_user.GetUsersByPaginationRequest{
		Page:    "1",
		PerPage: "10",
	}
	testTable := []struct {
		scenario                      string
		expectedUser                  []*models.User
		expectedGetUsers              int32
		inputGetUsersByPageAndPerPage *service_user.GetUsersByPaginationRequest
		response                      interface{}
		cacheErr                      error
		expectedError                 error
	}{
		{
			"get users by pagination success",
			users,
			2,
			req,
			mappers.MapUsersToGetUsersByPaginationResponse(users, "130", "100500"),
			&apperrors.RedisKeyDoesntExist,
			nil,
		},
		{
			"get users by pagination success",
			users,
			2,
			req,
			mappers.MapUsersToGetUsersByPaginationResponse(users, "1", "10"),
			&apperrors.RedisKeyDoesntExist,
			nil,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			pageNum, err := strconv.Atoi(tc.inputGetUsersByPageAndPerPage.Page)
			if err != nil {
				t.Error(err)
				return
			}

			perPageNum, err := strconv.Atoi(tc.inputGetUsersByPageAndPerPage.PerPage)
			if err != nil {
				t.Error(err)
				return
			}

			userRepoMock := mock.NewMockUserRepository(ctrl)
			userRepoMock.EXPECT().GetUsersByPageAndPerPage(ctx, pageNum,
				perPageNum).Return(tc.expectedUser, tc.expectedError).AnyTimes()
			userHacheMock := mock.NewMockUserCache(ctrl)
			userHacheMock.EXPECT().Get(ctx, gomock.Any()).Return("", tc.cacheErr).AnyTimes()
			userHacheMock.EXPECT().SetWithExpiration(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			uInteractor := interactor.NewUserInteractor(userRepoMock, userHacheMock)
			log := logrus.New()
			config := config.Config{}
			uController := NewUserController(uInteractor, log, &config)

			getUsersByPPResp, err := uController.GetUsersByPagination(ctx, tc.inputGetUsersByPageAndPerPage)
			if err != nil {
				apperrors.IsAppError(err, tc.expectedError.(*apperrors.AppError))
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, tc.expectedGetUsers, getUsersByPPResp.TotalUsers)
		})
	}
}

func TestDropUserByID(t *testing.T) {
	log := logrus.New()
	req := &service_user.DropUserByIDRequest{
		Id: "1",
	}
	reqPositive := &service_user.DropUserByIDRequest{
		Id: "a700ce3c-e29a-4523-af3b-ff75d8ab40a5",
	}

	testTable := []struct {
		scenario      string
		inputUserReq  *service_user.DropUserByIDRequest
		response      string
		expectedError error
	}{
		{
			"invalid UUID length: 1",
			req,
			"invalid UUID length: 1",
			apperrors.CoverstionToUUIDFailedDropUserByID.AppendMessage("invalid UUID length: 1"),
		},
		{
			"delete user positive",
			reqPositive,
			"the user has been deleted",
			nil,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mock.NewMockUserRepository(ctrl)
			userHacheMock := mock.NewMockUserCache(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, userHacheMock)
			config := config.Config{}
			uController := NewUserController(uInteractor, log, &config)
			userRepoMock.EXPECT().DropUserByID(ctx, gomock.Any()).Return(tc.expectedError).AnyTimes()

			respDropUser, err := uController.DropUserById(ctx, tc.inputUserReq)

			if err != nil {
				log.Error(err)
				apperrors.IsAppError(err, tc.expectedError.(*apperrors.AppError))
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, tc.response, respDropUser.Result)

		})
	}
}

func TestUpdateUserByID(t *testing.T) {
	log := logrus.New()
	reqFailedID := &service_user.UpdateUserByIDRequest{
		Id:        "a700ce3c-e29a-4523-af3b-ff75d8ab40a",
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}
	req := &service_user.UpdateUserByIDRequest{
		Id:        "a700ce3c-e29a-4523-af3b-ff75d8ab40a5",
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	reqCreateUser := &service_user.CreateUserRequest{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Role:      "user",
	}

	user, err := mappers.MapCreateUserRequestToUser(reqCreateUser)
	if err != nil {
		log.Error(err)
	}

	UUID, err := uuid.Parse(req.Id)
	if err != nil {
		log.Error(err)
	}

	user.ID = &UUID

	testTable := []struct {
		scenario      string
		inputUserReq  *service_user.UpdateUserByIDRequest
		user          *models.User
		expectedError error
	}{
		{
			"update user UUID Body is broken",
			reqFailedID,
			user,
			&apperrors.CoverstionToUUIDFailedUpdateUserByID,
		},
		{
			"update user positive",
			req,
			user,
			nil,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mock.NewMockUserRepository(ctrl)
			userRepoMock.EXPECT().UpdateUserByID(ctx, gomock.Any()).Return(tc.user.Email, tc.expectedError).AnyTimes()
			userHacheMock := mock.NewMockUserCache(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, userHacheMock)
			cfg := config.Config{SecretKey: "secret"}
			uController := NewUserController(uInteractor, log, &cfg)
			respUpdateUser, err := uController.UpdateUserByID(ctx, tc.inputUserReq)
			if err != nil {
				apperrors.IsAppError(err, tc.expectedError.(*apperrors.AppError))
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, tc.user.Email, respUpdateUser.Email)
		})
	}
}
