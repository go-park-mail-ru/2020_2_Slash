package admin

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
)

type AdminMicroservice struct {
	actorsRep    actor.ActorRepository
	directorsRep director.DirectorRepository
	countriesRep country.CountryRepository
	genresRep    genre.GenreRepository
	moviesRep    movie.MovieRepository
	contentRep   content.ContentRepository
	seasonsRep   season.SeasonRepository
	episodesRep  episode.EpisodeRepository
	tvshowsRep   tvshow.TVShowRepository
}

func NewAdminMicroservice(actorsRep actor.ActorRepository,
	directorsRep director.DirectorRepository,
	countriesRep country.CountryRepository,
	genresRep genre.GenreRepository,
	moviesRep movie.MovieRepository,
	contentRep content.ContentRepository,
	seasonsRep season.SeasonRepository,
	episodesRep episode.EpisodeRepository,
	tvshowsRep tvshow.TVShowRepository) AdminPanelServer {
	return &AdminMicroservice{
		actorsRep:    actorsRep,
		directorsRep: directorsRep,
		countriesRep: countriesRep,
		genresRep:    genresRep,
		moviesRep:    moviesRep,
		contentRep:   contentRep,
		seasonsRep:   seasonsRep,
		episodesRep:  episodesRep,
		tvshowsRep:   tvshowsRep,
	}
}
