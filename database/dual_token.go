package database

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const (
	TOKEN_PREFIX = "dual_token_"
	TOKEN_EXPIRE = 7 * 24 * time.Hour
)

func SetToken(refreshToken, authToken string) {
	client := InitRedisClient()
	err := client.Set(TOKEN_PREFIX+refreshToken, authToken, TOKEN_EXPIRE).Err()
	if err != nil {
		zap.L().Error("set token failed", zap.String("refresh_token", refreshToken), zap.String("auth_token", authToken), zap.Error(err))
	}

}

func GetToken(refreshToken string) string {
	client := InitRedisClient()
	result, err := client.Get(TOKEN_PREFIX + refreshToken).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zap.L().Error("get token failed", zap.String("refresh_token", refreshToken), zap.Error(err))
		}

	}
	return result
}
