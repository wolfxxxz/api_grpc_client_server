package datastore

import (
	"context"
	"fmt"
	"service_user/internal/apperrors"
	"service_user/internal/config"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func InintRedisDB(ctx context.Context, conf *config.Config, log *logrus.Logger) (*redis.Client, error) {
	address := fmt.Sprintf("%v:%v", conf.RedisHost, conf.RedisPort)
	dbIndex, err := strconv.Atoi(conf.RedisDB)
	if err != nil {
		appErr := apperrors.RedisInitErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: conf.RedisPassword,
		DB:       dbIndex,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		appErr := apperrors.RedisInitErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	log.Infof("Redis pong success: %v", pong)

	return client, nil
}
