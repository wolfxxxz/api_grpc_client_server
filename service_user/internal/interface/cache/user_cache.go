package cache

import (
	"context"
	"service_user/internal/apperrors"
	"service_user/internal/usecase/cache"

	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type userCache struct {
	log     *logrus.Logger
	redisDB *redis.Client
}

func NewUserCache(log *logrus.Logger, redisDB *redis.Client) cache.UserCache {
	return &userCache{log: log, redisDB: redisDB}
}

func (uc *userCache) Get(ctx context.Context, key string) (string, error) {
	cachedData, err := uc.redisDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			appErr := apperrors.RedisKeyDoesntExist.AppendMessage(err)
			uc.log.Info(appErr)
			return "", appErr
		}

		appErr := apperrors.RedisGetUsersByPageAndPerPageErr.AppendMessage(err)
		uc.log.Error(appErr)
		return "", appErr
	}

	uc.log.Info("GET data from REDIS - CACHE is working")
	return cachedData, nil
}

func (uc *userCache) SetWithExpiration(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	err := uc.redisDB.Set(ctx, key, data, expiration).Err()
	if err != nil {
		appErr := apperrors.RedisGetUserByIDErr.AppendMessage(err)
		uc.log.Error("redisDB.Set ", appErr)
		return appErr
	}

	uc.log.Info("SET data into REDIS - CACHE is working")
	return nil
}
