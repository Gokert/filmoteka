package core

import (
	"context"
	"filmoteka/pkg/models"
)

type IFilms interface {
	GetFilms(ctx context.Context, request *models.FindFilmRequest) (*[]models.FilmItem, error)
	AddFilm(ctx context.Context, film *models.FilmRequest, actors []uint64) (uint64, error)
	SearchFilms(ctx context.Context, titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error)
	UpdateFilm(ctx context.Context, film *models.FilmRequest) error
	DeleteFilm(ctx context.Context, filmId uint64) (bool, error)
}
