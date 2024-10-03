package main

import (
	"context"
	"net"
	"service_user/internal/apperrors"
	"service_user/internal/config"
	"service_user/internal/infrastructure/datastore"
	"service_user/internal/log"
	"service_user/internal/registry"
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

	mongoDB, err := datastore.InintMongoDB(ctx, conf, logger)
	if err != nil {
		logger.Error(err)
		return
	}

	redisDB, err := datastore.InintRedisDB(ctx, conf, logger)
	if err != nil {
		logger.Error(err)
		return
	}

	r := registry.NewRegistry(mongoDB, redisDB, logger, conf)
	appContr := r.NewAppController()

	s := grpc.NewServer()
	service_user.RegisterUserServiceServer(s, appContr.UserController)

	l, err := net.Listen(conf.Protocol, ":"+conf.Port)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Server are listening in protocol %v and port :%v", conf.Protocol, conf.Port)
	defer l.Close()

	if err := s.Serve(l); err != nil {
		logger.Fatal(err)
	}
}
