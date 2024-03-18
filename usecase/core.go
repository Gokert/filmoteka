package usecase

import (
	"context"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"filmoteka/repository/psx"
	"filmoteka/repository/session"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

type Core struct {
	log      *logrus.Logger
	mutex    sync.RWMutex
	films    psx.IFilmRepo
	profiles psx.IProfileRepo
	sessions session.ISessionRepo
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetCore(psxCfg *configs.DbPsxConfig, redisCfg *configs.DbRedisCfg, log *logrus.Logger) (*Core, error) {
	filmRepo, err := psx.GetFilmRepo(psxCfg, log)
	if err != nil {
		log.Error("Get GetFilmRepo error: ", err)
		return nil, err
	}

	authRepo, err := session.GetAuthRepo(redisCfg, log)
	if err != nil {
		log.Error("Get GetAuthRepo error: ", err)
		return nil, err
	}

	core := &Core{
		log:      log,
		films:    filmRepo,
		profiles: filmRepo,
		sessions: authRepo,
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

func (c *Core) GetUserId(ctx context.Context, sid string) (uint64, error) {
	c.mutex.RLock()
	login, err := c.sessions.GetUserLogin(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		c.log.Errorf("get user login error: %s", err.Error())
		return 0, fmt.Errorf("get user login error: %s", err.Error())
	}

	id, err := c.profiles.GetUserId(login)
	if err != nil {
		c.log.Errorf("get user id error: %s", err.Error())
		return 0, fmt.Errorf("get user id error: %s", err.Error())
	}

	return id, nil
}

func (c *Core) GetRole(userId uint64) (string, error) {
	role, err := c.profiles.GetRole(userId)
	if err != nil {
		c.log.Errorf("GetRole error: %s", err.Error())
		return "", fmt.Errorf("GetRole error: %s", err.Error())
	}

	return role, nil
}

func (c *Core) AddFilm(film *models.FilmRequest, actors []uint64) (uint64, error) {
	filmId, err := c.films.AddFilm(film)
	if err != nil {
		c.log.Error("AddFilm error: ", err)
		return 0, fmt.Errorf("AddFilm error: %w", err)
	}

	err = c.films.AddActorsForFilm(filmId, actors)
	if err != nil {
		c.log.Error("AddActorsForFilm error: ", err.Error())
		return 0, fmt.Errorf("AddActorsForFilm error: %w", err)
	}

	return filmId, nil
}

func (c *Core) AddActor(actor *models.ActorItem) (uint64, error) {
	actorId, err := c.films.AddActor(actor)
	if err != nil {
		c.log.Error("AddActorsForFilm error: ", err.Error())
		return 0, fmt.Errorf("AddActor error: %w", err)
	}

	return actorId, nil
}

func (c *Core) SearchFilms(titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error) {
	films, err := c.films.SearchFilms(titleFilm, nameActor, page, perPage)
	if err != nil {
		c.log.Error("SearchFilms error: ", err.Error())
		return nil, err
	}

	return films, nil
}

func (c *Core) UpdateFilm(film *models.FilmRequest) error {
	err := c.films.UpdateFilm(film)
	if err != nil {
		c.log.Error("ChangeFilm error: ", err.Error())
		return fmt.Errorf("ChangeFilm error: ", err.Error())
	}

	return nil
}

func (c *Core) UpdateActor(actor *models.ActorRequest) error {
	err := c.films.UpdateActor(actor)
	if err != nil {
		c.log.Errorf("ChangeActor error: %s", err.Error())
		return fmt.Errorf("ChangeActor error: %s", err.Error())
	}

	return nil
}

func (c *Core) DeleteFilm(filmId uint64) (bool, error) {
	_, err := c.films.DeleteFilm(filmId)
	if err != nil {
		return false, fmt.Errorf("DeleteFilm error: %w", err)
	}

	return true, nil
}

func (c *Core) FindActors(page uint64, perPage uint64) ([]models.ActorResponse, error) {
	actors, err := c.films.FindActors(page, perPage)
	if err != nil {
		c.log.Error("FindActors error: ", err.Error())
		return nil, fmt.Errorf("FindActors error: ", err.Error())
	}

	return actors, nil
}

func (c *Core) DeleteActor(actorId uint64) error {
	err := c.films.DeleteActor(actorId)
	if err != nil {
		c.log.Errorf("DeleteActor error: %s", err.Error())
		return fmt.Errorf("DeleteActor error: %s", err.Error())
	}

	return nil
}

func (c *Core) GetUserName(ctx context.Context, sid string) (string, error) {
	c.mutex.RLock()
	login, err := c.sessions.GetUserLogin(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		return "", err
	}

	return login, nil
}

func (c *Core) CreateSession(ctx context.Context, login string) (models.Session, error) {
	sid := RandStringRunes(32)

	newSession := models.Session{
		Login:     login,
		SID:       sid,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	c.mutex.Lock()
	sessionAdded, err := c.sessions.AddSession(ctx, newSession, c.log)
	c.mutex.Unlock()

	if !sessionAdded && err != nil {
		return models.Session{}, err
	}

	if !sessionAdded {
		return models.Session{}, nil
	}

	return newSession, nil
}

func (c *Core) FindActiveSession(ctx context.Context, sid string) (bool, error) {
	c.mutex.RLock()
	login, err := c.sessions.CheckActiveSession(ctx, sid, c.log)
	c.mutex.RUnlock()

	c.log.Warning(login)

	if err != nil {
		return false, err
	}

	return login, nil
}

func (c *Core) KillSession(ctx context.Context, sid string) error {
	c.mutex.Lock()
	_, err := c.sessions.DeleteSession(ctx, sid, c.log)
	c.mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (c *Core) CreateUserAccount(login string, password string) error {
	err := c.profiles.CreateUser(login, password)
	if err != nil {
		c.log.Error("create user error: ", err.Error())
		return fmt.Errorf("CreateUserAccount err: %w", err)
	}

	return nil
}

func (c *Core) FindUserAccount(login string, password string) (*models.UserItem, bool, error) {
	user, found, err := c.profiles.GetUser(login, password)
	if err != nil {
		c.log.Error("find user error: ", err.Error())
		return nil, false, fmt.Errorf("FindUserAccount error: %w", err)
	}
	return user, found, nil
}

func (c *Core) FindUserByLogin(login string) (bool, error) {
	found, err := c.profiles.FindUser(login)
	if err != nil {
		c.log.Error("find user error", "err", err.Error())
		return false, fmt.Errorf("FindUserByLogin err: %w", err)
	}

	return found, nil
}

func RandStringRunes(seed int) string {
	symbols := make([]rune, seed)
	for i := range symbols {
		symbols[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(symbols)
}
