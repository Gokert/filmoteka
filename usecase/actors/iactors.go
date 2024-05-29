package core

import (
	"context"
	"filmoteka/pkg/models"
)

type IActors interface {
	AddActor(ctx context.Context, actor *models.ActorItem) (uint64, error)
	FindActors(ctx context.Context, page uint64, perPage uint64) ([]models.ActorResponse, error)
	UpdateActor(ctx context.Context, actor *models.ActorRequest) error
	DeleteActor(ctx context.Context, actorId uint64) error
}
