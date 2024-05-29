package core

import (
	"context"
	"filmoteka/pkg/models"
)

type IProfiles interface {
	CreateUserAccount(ctx context.Context, login string, password string) error
	FindUserAccount(ctx context.Context, login string, password string) (*models.UserItem, bool, error)
	FindUserByLogin(ctx context.Context, login string) (bool, error)
	GetRole(ctx context.Context, userId uint64) (string, error)
}
