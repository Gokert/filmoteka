package auth_repo

import (
	"context"
	"filmoteka/pkg/models"
	"github.com/sirupsen/logrus"
)

type IAuthRepo interface {
	AddSession(ctx context.Context, active models.Session, log *logrus.Logger) (bool, error)
	CheckActiveSession(ctx context.Context, sid string, lg *logrus.Logger) (bool, error)
	GetUserLogin(ctx context.Context, sid string, lg *logrus.Logger) (string, error)
	DeleteSession(ctx context.Context, sid string, lg *logrus.Logger) (bool, error)
}
