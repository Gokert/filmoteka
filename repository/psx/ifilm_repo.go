package psx

import (
	"context"
	"filmoteka/pkg/models"
)

type IFilmRepo interface {
	GetFilms(ctx context.Context, request *models.FindFilmRequest) (*[]models.FilmItem, error)
	AddFilm(ctx context.Context, film *models.FilmRequest) (uint64, error)
	AddActor(ctx context.Context, actor *models.ActorItem) (uint64, error)
	AddActorsForFilm(ctx context.Context, filmId uint64, actors []uint64) error
	SearchFilms(ctx context.Context, titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error)
	UpdateFilm(ctx context.Context, film *models.FilmRequest) error
	DeleteFilm(ctx context.Context, filmId uint64) (bool, error)

	FindActors(ctx context.Context, page uint64, perPage uint64) ([]models.ActorResponse, error)
	UpdateActor(ctx context.Context, actor *models.ActorRequest) error
	DeleteActor(ctx context.Context, actorId uint64) error
}
