package usecase

import (
	"context"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"filmoteka/repository/auth_repo"
	"filmoteka/repository/psx_repo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

type Core struct {
	log   *logrus.Logger
	mutex sync.RWMutex
	films *psx_repo.PsxRepo
	auth  auth_repo.IAuthRepo
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetCore(psxCfg *configs.DbPsxConfig, redisCfg *configs.DbRedisCfg, log *logrus.Logger) (*Core, error) {
	filmRepo, err := psx_repo.GetFilmRepo(psxCfg, log)
	if err != nil {
		log.Error("Get GetFilmRepo error: ", err)
		return nil, err
	}

	authRepo, err := auth_repo.GetAuthRepo(redisCfg, log)
	if err != nil {
		log.Error("Get GetAuthRepo error: ", err)
		return nil, err
	}

	core := &Core{
		log:   log,
		films: filmRepo,
		auth:  authRepo,
	}

	return core, nil
}

func (c *Core) GetFilms(request *models.FindFilmRequest) (*[]models.FilmItem, error) {
	films, err := c.films.GetFilms(request)
	if err != nil {
		c.log.Error("GetFilms error: ", err)
		return nil, err
	}

	return films, nil
}

func (c *Core) GetUserName(ctx context.Context, sid string) (string, error) {
	c.mutex.RLock()
	login, err := c.auth.GetUserLogin(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		return "", err
	}

	return login, nil
}

func (c *Core) CreateSession(ctx context.Context, login string) (string, models.Session, error) {
	sid := RandStringRunes(32)

	newSession := models.Session{
		Login:     login,
		SID:       sid,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	c.mutex.Lock()
	sessionAdded, err := c.auth.AddSession(ctx, newSession, c.log)
	c.mutex.Unlock()

	if !sessionAdded && err != nil {
		return "", models.Session{}, err
	}

	if !sessionAdded {
		return "", models.Session{}, nil
	}

	return sid, newSession, nil
}

func (c *Core) FindActiveSession(ctx context.Context, sid string) (bool, error) {
	c.mutex.RLock()
	found, err := c.auth.CheckActiveSession(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		return false, err
	}

	return found, nil
}

func (c *Core) KillSession(ctx context.Context, sid string) error {
	c.mutex.Lock()
	_, err := c.auth.DeleteSession(ctx, sid, c.log)
	c.mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func RandStringRunes(seed int) string {
	symbols := make([]rune, seed)
	for i := range symbols {
		symbols[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(symbols)
}
