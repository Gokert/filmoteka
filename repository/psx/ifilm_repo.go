package psx

import "filmoteka/pkg/models"

type IFilmRepo interface {
	GetFilms(request *models.FindFilmRequest) (*[]models.FilmItem, error)
	AddFilm(film *models.FilmRequest) (uint64, error)
	AddActor(actor *models.ActorItem) (uint64, error)
	AddActorsForFilm(filmId uint64, actors []uint64) error
	SearchFilms(titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error)
	UpdateFilm(film *models.FilmRequest) error
	DeleteFilm(filmId uint64) (bool, error)

	FindActors(page uint64, perPage uint64) ([]models.ActorResponse, error)
	UpdateActor(actor *models.ActorRequest) error
	DeleteActor(actorId uint64) error
}
