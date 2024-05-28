package psx

import (
	"context"
	"filmoteka/pkg/models"
)

type IProfileRepo interface {
	GetUser(ctx context.Context, login string, password []byte) (*models.UserItem, bool, error)
	FindUser(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login string, password []byte) error
	GetUserId(ctx context.Context, login string) (uint64, error)
	GetRole(ctx context.Context, userId uint64) (string, error)
}
