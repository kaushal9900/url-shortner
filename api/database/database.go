package database

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/kaushal9900/url-shortner/configs"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.EnvConfigs.DBAddress,
		Password: configs.EnvConfigs.DBPassword,
		DB:       dbNo,
	})
	return rdb
}
