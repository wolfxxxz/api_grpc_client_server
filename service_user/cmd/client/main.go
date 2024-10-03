package main

import (
	"context"
	"service_user/internal/apperrors"
	"service_user/internal/config"
	"service_user/internal/log"
	"service_user/internal/service_user"

	"google.golang.org/grpc"
)

func main() {
	logger, err := log.NewLogAndSetLevel("info")
	if err != nil {
		logger.Error(err)
		return
	}

	conf := config.NewConfig()
	err = conf.ParseConfig("config/.env", logger)
	if err != nil {
		logger.Error(apperrors.EnvConfigLoadError.AppendMessage(err))
		return
	}

	if err = log.SetLevel(logger, conf.LogLevel); err != nil {
		logger.Error(err)
		return
	}

	ctx := context.Background()
	conn, err := grpc.Dial(":"+conf.Port, grpc.WithInsecure())
	if err != nil {
		logger.Error(err)
		return
	}

	client := service_user.NewUserServiceClient(conn)
	respCreateUserID, err := createUser(ctx, client)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info(respCreateUserID)

	respDropUserByID, err := dropUserByID(ctx, client, respCreateUserID)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info(respDropUserByID)
}

func createUser(ctx context.Context, client service_user.UserServiceClient) (string, error) {
	reqCreateUser := &service_user.CreateUserRequest{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Role:      "user",
	}

	responce, err := client.CreateUser(ctx, reqCreateUser)
	if err != nil {

		return "", err
	}

	return responce.GetUserId(), nil
}

func dropUserByID(ctx context.Context, client service_user.UserServiceClient, userId string) (string, error) {
	reqCreateUser := &service_user.DropUserByIDRequest{
		Id: userId,
	}

	responce, err := client.DropUserById(ctx, reqCreateUser)
	if err != nil {

		return "", err
	}

	return responce.String(), nil
}
