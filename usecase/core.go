package usecase

import (
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"filmoteka/repository/auth_repo"
	"filmoteka/repository/psx_repo"
	"github.com/sirupsen/logrus"
)

type Core struct {
	log   *logrus.Logger
	films *psx_repo.PsxRepo
	auth  *auth_repo.AuthRepo
}

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

func (c *Core) GetAll(start uint64, end uint64, order bool) (*[]models.FilmItem, error) {
	films, err := c.films.GetAllFilms(start, end, order)
	if err != nil {
		c.log.Error("GetAll error: ", err)
		return nil, err
	}

	return films, nil
}
