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

func DirectorModelToGRPC(director *models.Director) *Director {
	return &Director{
		ID:   director.ID,
		Name: director.Name,
	}
}

func DirectorGRPCToModel(director *Director) *models.Director {
	return &models.Director{
		ID:   director.GetID(),
		Name: director.GetName(),
	}
}
