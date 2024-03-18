package psx

import "filmoteka/pkg/models"

type IProfileRepo interface {
	GetUser(login string, password string) (*models.UserItem, bool, error)
	FindUser(login string) (bool, error)
	CreateUser(login string, password string) error
	GetUserId(login string) (uint64, error)
	GetRole(userId uint64) (string, error)
}
