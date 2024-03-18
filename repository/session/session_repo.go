package session

import (
	"context"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

type SessionRepo struct {
	DB *redis.Client
}

func GetAuthRepo(cfg *configs.DbRedisCfg, log *logrus.Logger) (ISessionRepo, error) {
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
	return &SessionRepo{DB: redisClient}, nil
}

func (repo *SessionRepo) AddSession(ctx context.Context, active models.Session, log *logrus.Logger) (bool, error) {

	repo.DB.Set(ctx, active.SID, active.Login, 24*time.Hour)

	added, err := repo.CheckActiveSession(ctx, active.SID, log)
	if err != nil {
		return false, err
	}

	return added, nil
}

func (repo *SessionRepo) CheckActiveSession(ctx context.Context, sid string, lg *logrus.Logger) (bool, error) {
	_, err := repo.DB.Get(ctx, sid).Result()
	if err == redis.Nil {
		lg.Error("Key " + sid + " not found")
		return false, nil
	}

	if err != nil {
		lg.Error("Get request could not be completed ", err)
		return false, err
	}

	return true, err
}

func (repo *SessionRepo) GetUserLogin(ctx context.Context, sid string, lg *logrus.Logger) (string, error) {
	value, err := repo.DB.Get(ctx, sid).Result()
	if err != nil {
		lg.Error("Error, cannot find session " + sid)
		return "", err
	}

	return value, nil
}

func (repo *SessionRepo) DeleteSession(ctx context.Context, sid string, lg *logrus.Logger) (bool, error) {
	_, err := repo.DB.Del(ctx, sid).Result()
	if err != nil {
		lg.Error("Delete request could not be completed:", err)
		return false, err
	}

	return true, nil
}
