package usecase

import (
	"context"
	"filmoteka/pkg/models"
)

type ICore interface {
	GetFilms(request *models.FindFilmRequest) (*[]models.FilmItem, error)
	GetUserName(ctx context.Context, sid string) (string, error)
	CreateSession(ctx context.Context, login string) (string, models.Session, error)
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	KillSession(ctx context.Context, sid string) error
}
