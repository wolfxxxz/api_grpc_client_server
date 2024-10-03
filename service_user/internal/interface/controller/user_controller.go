package controller

import (
	context "context"
	"service_user/internal/apperrors"
	"service_user/internal/config"
	"service_user/internal/service_user"
	"service_user/internal/usecase/interactor"

	"github.com/sirupsen/logrus"
)

type UserController struct {
	userInteractor interactor.UserInteractor
	log            *logrus.Logger
	config         *config.Config
}

func NewUserController(us interactor.UserInteractor, log *logrus.Logger, config *config.Config) *UserController {
	return &UserController{us, log, config}
}

func (uc *UserController) GetUserById(c context.Context,
	req *service_user.GetUserByIDRequest) (*service_user.GetUserByIDResponse, error) {

	return uc.userInteractor.GetUserByID(c, req.Id)
}

func (uc *UserController) GetUserByEmail(c context.Context,
	req *service_user.GetUserByEmailRequest) (*service_user.GetUserByEmailResponse, error) {

	return uc.userInteractor.GetUserByEmail(c, req.Email)
}

func (uc *UserController) UpdateUserByID(c context.Context,
	req *service_user.UpdateUserByIDRequest) (*service_user.UpdateUserByIDResponse, error) {

	userEmailResponse, err := uc.userInteractor.UpdateUserByID(c, req)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		uc.log.Error(appErr)
		return nil, appErr
	}

	return &service_user.UpdateUserByIDResponse{Email: userEmailResponse}, nil
}

func (uc *UserController) DropUserById(c context.Context,
	req *service_user.DropUserByIDRequest) (*service_user.DropUserByIDResponse, error) {

	err := uc.userInteractor.DropUserByID(c, req)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		uc.log.Error(appErr)
		return nil, appErr
	}

	return &service_user.DropUserByIDResponse{Result: "the user has been deleted"}, nil
}

func (uc *UserController) CreateUser(c context.Context, req *service_user.CreateUserRequest) (*service_user.CreateUserResponse, error) {
	userId, err := uc.userInteractor.CreateUser(c, req)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		uc.log.Error(appErr)
		return nil, appErr
	}

	return &service_user.CreateUserResponse{UserId: userId}, nil
}

func (uc *UserController) GetUsersByPagination(c context.Context,
	req *service_user.GetUsersByPaginationRequest) (*service_user.GetUsersByPaginationResponse, error) {

	return uc.userInteractor.GetUsersByPageAndPerPage(c, req.Page, req.PerPage)
}
