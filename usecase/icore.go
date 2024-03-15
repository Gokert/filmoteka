package usecase

import "filmoteka/pkg/models"

type ICore interface {
	GetAll(start uint64, end uint64, order bool) (*[]models.FilmItem, error)
	GetFilms(request *models.FindFilmRequest) (*[]models.FilmItem, error)
}
