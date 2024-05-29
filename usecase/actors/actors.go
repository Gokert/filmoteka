package core

import (
	"context"
	"filmoteka/pkg/models"
	"filmoteka/repository/psx"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Actors struct {
	log    *logrus.Logger
	actors psx.IActorRepo
}

func NewCoreActors(actors psx.IActorRepo, log *logrus.Logger) *Actors {
	return &Actors{
		log:    log,
		actors: actors,
	}
}

func (c *Actors) AddActor(ctx context.Context, actor *models.ActorItem) (uint64, error) {
	actorId, err := c.actors.AddActor(ctx, actor)
	if err != nil {
		c.log.Errorf("add actor error: %s", err.Error())
		return 0, fmt.Errorf("add actor error: %s", err.Error())
	}

	return actorId, nil
}

func (c *Actors) FindActors(ctx context.Context, page uint64, perPage uint64) ([]models.ActorResponse, error) {
	actors, err := c.actors.FindActors(ctx, page, perPage)
	if err != nil {
		c.log.Errorf("find actors error: %s", err.Error())
		return nil, fmt.Errorf("find actors error: %s", err.Error())
	}

	return actors, nil
}

func (c *Actors) UpdateActor(ctx context.Context, actor *models.ActorRequest) error {
	err := c.actors.UpdateActor(ctx, actor)
	if err != nil {
		c.log.Errorf("change actor error: %s", err.Error())
		return fmt.Errorf("change actor error: %s", err.Error())
	}

	return nil
}

func (c *Actors) DeleteActor(ctx context.Context, actorId uint64) error {
	err := c.actors.DeleteActor(ctx, actorId)
	if err != nil {
		c.log.Errorf("delete actor error: %s", err.Error())
		return fmt.Errorf("delete actor error: %s", err.Error())
	}

	return nil
}
