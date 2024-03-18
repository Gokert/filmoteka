package usecase

import (
	"context"
	"filmoteka/pkg/models"
)

type ICore interface {
	GetFilms(request *models.FindFilmRequest) (*[]models.FilmItem, error)
	AddFilm(film *models.FilmRequest, actors []uint64) (uint64, error)
	SearchFilms(titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error)
	UpdateFilm(film *models.FilmRequest) error
	DeleteFilm(filmId uint64) (bool, error)

	AddActor(actor *models.ActorItem) (uint64, error)
	FindActors(page uint64, perPage uint64) ([]models.ActorResponse, error)
	UpdateActor(actor *models.ActorRequest) error
	DeleteActor(actorId uint64) error

	GetUserName(ctx context.Context, sid string) (string, error)
	CreateSession(ctx context.Context, login string) (models.Session, error)
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	KillSession(ctx context.Context, sid string) error
	GetUserId(ctx context.Context, sid string) (uint64, error)

	CreateUserAccount(login string, password string) error
	FindUserAccount(login string, password string) (*models.UserItem, bool, error)
	FindUserByLogin(login string) (bool, error)
	GetRole(userId uint64) (string, error)
}
