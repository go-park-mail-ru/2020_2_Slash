package admin

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func ActorModelToGRPC(actor *models.Actor) *Actor {
	return &Actor{
		ID:   actor.ID,
		Name: actor.Name,
	}
}

func ActorGRPCToModel(actor *Actor) *models.Actor {
	return &models.Actor{
		ID:   actor.GetID(),
		Name: actor.GetName(),
	}
}
