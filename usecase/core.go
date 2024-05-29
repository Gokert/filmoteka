package usecase

import (
	"filmoteka/configs"
	"filmoteka/repository/psx"
	"filmoteka/repository/session"
	core_actor "filmoteka/usecase/actors"
	core_films "filmoteka/usecase/films"
	core_profiles "filmoteka/usecase/profiles"
	core_sessions "filmoteka/usecase/sessions"
	"github.com/sirupsen/logrus"
)

type Core struct {
	log      *logrus.Logger
	Films    core_films.IFilms
	Actors   core_actor.IActors
	Profiles core_profiles.IProfiles
	Sessions core_sessions.ISessions
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

	return &Core{
		log:      log,
		Films:    core_films.NewCoreFilms(filmRepo, log),
		Actors:   core_actor.NewCoreActors(filmRepo, log),
		Profiles: core_profiles.NewCoreProfiles(filmRepo, authRepo, log),
		Sessions: core_sessions.NewCoreSessions(filmRepo, authRepo, log),
	}, nil
}
