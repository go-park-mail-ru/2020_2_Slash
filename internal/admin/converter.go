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

func GenreModelToGRPC(genre *models.Genre) *Genre {
	return &Genre{
		ID:   genre.ID,
		Name: genre.Name,
	}
}

func GenreGRPCToModel(genre *Genre) *models.Genre {
	return &models.Genre{
		ID:   genre.GetID(),
		Name: genre.GetName(),
	}
}

func ContentModelToGRPC(content *models.Content) *Content {
	var countries []*Country
	for _, modelCountry := range content.Countries {
		grpcCountry := &Country{
			ID:   modelCountry.ID,
			Name: modelCountry.Name,
		}
		countries = append(countries, grpcCountry)
	}

	var genres []*Genre
	for _, modelGenre := range content.Genres {
		grpcGenre := &Genre{
			ID:   modelGenre.ID,
			Name: modelGenre.Name,
		}
		genres = append(genres, grpcGenre)
	}

	var actors []*Actor
	for _, modelActor := range content.Actors {
		grpcActor := &Actor{
			ID:   modelActor.ID,
			Name: modelActor.Name,
		}
		actors = append(actors, grpcActor)
	}

	var directors []*Director
	for _, modelDirector := range content.Directors {
		grpcDirector := &Director{
			ID:   modelDirector.ID,
			Name: modelDirector.Name,
		}
		directors = append(directors, grpcDirector)
	}

	return &Content{
		ID:               content.ContentID,
		Name:             content.Name,
		OriginalName:     content.OriginalName,
		Description:      content.Description,
		ShortDescription: content.ShortDescription,
		Rating:           int64(content.Rating),
		Year:             int64(content.Year),
		Images:           content.Images,
		Type:             content.Type,
		Countries:        countries,
		Genres:           genres,
		Actors:           actors,
		Directors:        directors,
		IsLiked:          *content.IsLiked,
		IsFavourite:      *content.IsFavourite,
	}
}

func ContentGRPCToModel(content *Content) *models.Content {
	var countries []*models.Country
	for _, grpcCountry := range content.Countries {
		modelCountry := &models.Country{
			ID:   grpcCountry.ID,
			Name: grpcCountry.Name,
		}
		countries = append(countries, modelCountry)
	}

	var genres []*models.Genre
	for _, grpcGenre := range content.Genres {
		modelGenre := &models.Genre{
			ID:   grpcGenre.ID,
			Name: grpcGenre.Name,
		}
		genres = append(genres, modelGenre)
	}

	var actors []*models.Actor
	for _, grpcActor := range content.Actors {
		modelActor := &models.Actor{
			ID:   grpcActor.ID,
			Name: grpcActor.Name,
		}
		actors = append(actors, modelActor)
	}

	var directors []*models.Director
	for _, grpcDirector := range content.Directors {
		modelDirector := &models.Director{
			ID:   grpcDirector.ID,
			Name: grpcDirector.Name,
		}
		directors = append(directors, modelDirector)
	}

	return &models.Content{
		ContentID:        content.ID,
		Name:             content.Name,
		OriginalName:     content.OriginalName,
		Description:      content.Description,
		ShortDescription: content.ShortDescription,
		Rating:           int(content.Rating),
		Year:             int(content.Year),
		Images:           content.Images,
		Type:             content.Type,
		Countries:        countries,
		Genres:           genres,
		Actors:           actors,
		Directors:        directors,
		IsLiked:          &content.IsLiked,
		IsFavourite:      &content.IsFavourite,
	}
}
