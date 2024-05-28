package usecase

import (
	"context"
	"filmoteka/configs"
	utils "filmoteka/pkg"
	"filmoteka/pkg/models"
	"filmoteka/repository/psx"
	"filmoteka/repository/session"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Core struct {
	log      *logrus.Logger
	mutex    sync.RWMutex
	films    psx.IFilmRepo
	actors   psx.IActorRepo
	profiles psx.IProfileRepo
	sessions session.ISessionRepo
}

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
		actors:   filmRepo,
		profiles: filmRepo,
		sessions: authRepo,
	}

	return core, nil
}

func (c *Core) GetFilms(ctx context.Context, request *models.FindFilmRequest) (*[]models.FilmItem, error) {
	films, err := c.films.GetFilms(ctx, request)
	if err != nil {
		c.log.Errorf("get films error: %s", err.Error())
		return nil, fmt.Errorf("get films error: %s", err.Error())
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

	id, err := c.profiles.GetUserId(ctx, login)
	if err != nil {
		c.log.Errorf("get user id error: %s", err.Error())
		return 0, fmt.Errorf("get user id error: %s", err.Error())
	}

	return id, nil
}

func (c *Core) GetRole(ctx context.Context, userId uint64) (string, error) {
	role, err := c.profiles.GetRole(ctx, userId)
	if err != nil {
		c.log.Errorf("get role error: %s", err.Error())
		return "", fmt.Errorf("get role error: %s", err.Error())
	}

	return role, nil
}

func (c *Core) AddFilm(ctx context.Context, film *models.FilmRequest, actors []uint64) (uint64, error) {
	if film.Rating < utils.FilmRatingBegin || film.Rating > utils.FilmRatingEnd {
		c.log.Error(utils.RatingSizeError)
		return 0, fmt.Errorf(utils.RatingSizeError)
	}

	err := utils.ValidateStringSize(film.Title, utils.FilmTitleBegin, utils.FilmTitleEnd, utils.TitleSizeError, c.log)
	if err != nil {
		return 0, err
	}

	err = utils.ValidateStringSize(film.Info, utils.FilmDescriptionBegin, utils.FilmDescriptionEnd, utils.DescriptionSizeError, c.log)
	if err != nil {
		return 0, err
	}

	filmId, err := c.films.AddFilm(ctx, film)
	if err != nil {
		c.log.Error("add film error: ", err)
		return 0, fmt.Errorf("add film error: %w", err)
	}

	err = c.films.AddActorsForFilm(ctx, filmId, actors)
	if err != nil {
		c.log.Error("AddActorsForFilm error: ", err.Error())
		return 0, fmt.Errorf("AddActorsForFilm error: %w", err)
	}

	return filmId, nil
}

func (c *Core) AddActor(ctx context.Context, actor *models.ActorItem) (uint64, error) {
	actorId, err := c.actors.AddActor(ctx, actor)
	if err != nil {
		c.log.Errorf("add actor error: %s", err.Error())
		return 0, fmt.Errorf("add actor error: %s", err.Error())
	}

	return actorId, nil
}

func (c *Core) SearchFilms(ctx context.Context, titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error) {
	films, err := c.films.SearchFilms(ctx, titleFilm, nameActor, page, perPage)
	if err != nil {
		c.log.Errorf("SearchFilms error: %s", err.Error())
		return nil, fmt.Errorf("SearchFilms error: %s", err.Error())
	}

	return films, nil
}

func (c *Core) UpdateFilm(ctx context.Context, film *models.FilmRequest) error {
	err := c.films.UpdateFilm(ctx, film)
	if err != nil {
		c.log.Errorf("change film error: %s", err.Error())
		return fmt.Errorf("change film error: %s", err.Error())
	}

	return nil
}

func (c *Core) UpdateActor(ctx context.Context, actor *models.ActorRequest) error {
	err := c.actors.UpdateActor(ctx, actor)
	if err != nil {
		c.log.Errorf("change actor error: %s", err.Error())
		return fmt.Errorf("change actor error: %s", err.Error())
	}

	return nil
}

func (c *Core) DeleteFilm(ctx context.Context, filmId uint64) (bool, error) {
	_, err := c.films.DeleteFilm(ctx, filmId)
	if err != nil {
		c.log.Errorf("delete film error: %s", err.Error())
		return false, fmt.Errorf("delete film error: %s", err.Error())
	}

	return true, nil
}

func (c *Core) FindActors(ctx context.Context, page uint64, perPage uint64) ([]models.ActorResponse, error) {
	actors, err := c.actors.FindActors(ctx, page, perPage)
	if err != nil {
		c.log.Errorf("find actors error: %s", err.Error())
		return nil, fmt.Errorf("find actors error: %s", err.Error())
	}

	return actors, nil
}

func (c *Core) DeleteActor(ctx context.Context, actorId uint64) error {
	err := c.actors.DeleteActor(ctx, actorId)
	if err != nil {
		c.log.Errorf("delete actor error: %s", err.Error())
		return fmt.Errorf("delete actor error: %s", err.Error())
	}

	return nil
}

func (c *Core) GetUserName(ctx context.Context, sid string) (string, error) {
	c.mutex.RLock()
	login, err := c.sessions.GetUserLogin(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		c.log.Errorf("get user name error: %s", err.Error())
		return "", fmt.Errorf("get user name error: %s", err.Error())
	}

	return login, nil
}

func (c *Core) CreateSession(ctx context.Context, login string) (models.Session, error) {
	sid := utils.RandStringRunes(32)

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

	if err != nil {
		c.log.Errorf("find active session error: %s", err.Error())
		return false, fmt.Errorf("find active session error: %s", err.Error())
	}

	return login, nil
}

func (c *Core) KillSession(ctx context.Context, sid string) error {
	c.mutex.Lock()
	_, err := c.sessions.DeleteSession(ctx, sid, c.log)
	c.mutex.Unlock()

	if err != nil {
		c.log.Errorf("delete session error: %s", err.Error())
		return fmt.Errorf("delete sessionerror: %s", err.Error())
	}

	return nil
}

func (c *Core) CreateUserAccount(ctx context.Context, login string, password string) error {
	hashPassword := utils.HashPassword(password)
	err := c.profiles.CreateUser(ctx, login, hashPassword)
	if err != nil {
		c.log.Errorf("create user account error: %s", err.Error())
		return fmt.Errorf("create user account error: %s", err.Error())
	}

	return nil
}

func (c *Core) FindUserAccount(ctx context.Context, login string, password string) (*models.UserItem, bool, error) {
	hashPassword := utils.HashPassword(password)
	user, found, err := c.profiles.GetUser(ctx, login, hashPassword)
	if err != nil {
		c.log.Errorf("find user error: %s", err.Error())
		return nil, false, fmt.Errorf("find user account error: %s", err.Error())
	}
	return user, found, nil
}

func (c *Core) FindUserByLogin(ctx context.Context, login string) (bool, error) {
	found, err := c.profiles.FindUser(ctx, login)
	if err != nil {
		c.log.Errorf("find user by login error: %s", err.Error())
		return false, fmt.Errorf("find user by login error: %s", err.Error())
	}

	return found, nil
}
