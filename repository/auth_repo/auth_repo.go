package auth_repo

import (
	"context"
	"filmoteka/configs"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type AuthRepo struct {
	DB *redis.Client
}

func GetAuthRepo(cfg *configs.DbRedisCfg, log *logrus.Logger) (*AuthRepo, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.DbNumber,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Error("Ping redis error: ", err)
		return nil, err
	}

	log.Info("Redis created successful")
	return &AuthRepo{DB: redisClient}, nil
}
