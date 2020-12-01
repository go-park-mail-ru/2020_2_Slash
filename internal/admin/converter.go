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

func CountryModelToGRPC(country *models.Country) *Country {
	return &Country{
		ID:   country.ID,
		Name: country.Name,
	}
}

func CountryGRPCToModel(country *Country) *models.Country {
	return &models.Country{
		ID:   country.GetID(),
		Name: country.GetName(),
	}
}
