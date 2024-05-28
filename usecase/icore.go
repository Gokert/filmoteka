package usecase

import (
	"context"
	"filmoteka/pkg/models"
)

type ICore interface {
	GetFilms(ctx context.Context, request *models.FindFilmRequest) (*[]models.FilmItem, error)
	AddFilm(ctx context.Context, film *models.FilmRequest, actors []uint64) (uint64, error)
	SearchFilms(ctx context.Context, titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error)
	UpdateFilm(ctx context.Context, film *models.FilmRequest) error
	DeleteFilm(ctx context.Context, filmId uint64) (bool, error)

	AddActor(ctx context.Context, actor *models.ActorItem) (uint64, error)
	FindActors(ctx context.Context, page uint64, perPage uint64) ([]models.ActorResponse, error)
	UpdateActor(ctx context.Context, actor *models.ActorRequest) error
	DeleteActor(ctx context.Context, actorId uint64) error

	GetUserName(ctx context.Context, sid string) (string, error)
	CreateSession(ctx context.Context, login string) (models.Session, error)
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	KillSession(ctx context.Context, sid string) error
	GetUserId(ctx context.Context, sid string) (uint64, error)

	CreateUserAccount(ctx context.Context, login string, password string) error
	FindUserAccount(ctx context.Context, login string, password string) (*models.UserItem, bool, error)
	FindUserByLogin(ctx context.Context, login string) (bool, error)
	GetRole(ctx context.Context, userId uint64) (string, error)
}
