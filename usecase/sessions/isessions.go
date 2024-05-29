package core

import (
	"context"
	"filmoteka/pkg/models"
)

type ISessions interface {
	CreateSession(ctx context.Context, login string) (models.Session, error)
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	KillSession(ctx context.Context, sid string) error
	GetUserName(ctx context.Context, sid string) (string, error)
	GetUserId(ctx context.Context, sid string) (uint64, error)
}
