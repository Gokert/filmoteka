package core

import (
	"context"
	utils "filmoteka/pkg"
	"filmoteka/pkg/models"
	"filmoteka/repository/psx"
	"filmoteka/repository/session"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Profiles struct {
	log      *logrus.Logger
	profiles psx.IProfileRepo
	sessions session.ISessionRepo
}

func NewCoreProfiles(profiles psx.IProfileRepo, sessions session.ISessionRepo, log *logrus.Logger) *Profiles {
	return &Profiles{
		log:      log,
		profiles: profiles,
		sessions: sessions,
	}
}

func (c *Profiles) CreateUserAccount(ctx context.Context, login string, password string) error {
	hashPassword := utils.HashPassword(password)
	err := c.profiles.CreateUser(ctx, login, hashPassword)
	if err != nil {
		c.log.Errorf("create user account error: %s", err.Error())
		return fmt.Errorf("create user account error: %s", err.Error())
	}

	return nil
}

func (c *Profiles) FindUserAccount(ctx context.Context, login string, password string) (*models.UserItem, bool, error) {
	hashPassword := utils.HashPassword(password)
	user, found, err := c.profiles.GetUser(ctx, login, hashPassword)
	if err != nil {
		c.log.Errorf("find user error: %s", err.Error())
		return nil, false, fmt.Errorf("find user account error: %s", err.Error())
	}
	return user, found, nil
}

func (c *Profiles) FindUserByLogin(ctx context.Context, login string) (bool, error) {
	found, err := c.profiles.FindUser(ctx, login)
	if err != nil {
		c.log.Errorf("find user by login error: %s", err.Error())
		return false, fmt.Errorf("find user by login error: %s", err.Error())
	}

	return found, nil
}

func (c *Profiles) GetRole(ctx context.Context, userId uint64) (string, error) {
	role, err := c.profiles.GetRole(ctx, userId)
	if err != nil {
		c.log.Errorf("get role error: %s", err.Error())
		return "", fmt.Errorf("get role error: %s", err.Error())
	}

	return role, nil
}
