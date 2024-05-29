package core

import (
	"context"
	utils "filmoteka/pkg"
	"filmoteka/pkg/models"
	"filmoteka/repository/psx"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Films struct {
	log   *logrus.Logger
	films psx.IFilmRepo
}

func NewCoreFilms(films psx.IFilmRepo, log *logrus.Logger) *Films {
	return &Films{
		log:   log,
		films: films,
	}
}

func (c *Films) GetFilms(ctx context.Context, request *models.FindFilmRequest) (*[]models.FilmItem, error) {
	films, err := c.films.GetFilms(ctx, request)
	if err != nil {
		c.log.Errorf("get films error: %s", err.Error())
		return nil, fmt.Errorf("get films error: %s", err.Error())
	}

	return films, nil
}

func (c *Films) AddFilm(ctx context.Context, film *models.FilmRequest, actors []uint64) (uint64, error) {
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

func (c *Films) SearchFilms(ctx context.Context, titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error) {
	films, err := c.films.SearchFilms(ctx, titleFilm, nameActor, page, perPage)
	if err != nil {
		c.log.Errorf("SearchFilms error: %s", err.Error())
		return nil, fmt.Errorf("SearchFilms error: %s", err.Error())
	}

	return films, nil
}

func (c *Films) UpdateFilm(ctx context.Context, film *models.FilmRequest) error {
	err := c.films.UpdateFilm(ctx, film)
	if err != nil {
		c.log.Errorf("change film error: %s", err.Error())
		return fmt.Errorf("change film error: %s", err.Error())
	}

	return nil
}

func (c *Films) DeleteFilm(ctx context.Context, filmId uint64) (bool, error) {
	_, err := c.films.DeleteFilm(ctx, filmId)
	if err != nil {
		c.log.Errorf("delete film error: %s", err.Error())
		return false, fmt.Errorf("delete film error: %s", err.Error())
	}

	return true, nil
}
