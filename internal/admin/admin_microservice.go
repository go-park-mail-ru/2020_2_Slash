package admin

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
	"github.com/golang/protobuf/ptypes/empty"
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

func (am *AdminMicroservice) CreateDirector(ctx context.Context, d *Director) (*Director, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangeDirector(ctx context.Context, d *Director) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) DeleteDirectorByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) CreateCountry(ctx context.Context, c *Country) (*Country, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangeCountry(ctx context.Context, c *Country) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) DeleteCountryByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) CreateGenre(ctx context.Context, g *Genre) (*Genre, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangeGenre(ctx context.Context, g *Genre) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) DeleteGenreByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) CreateMovie(ctx context.Context, m *Movie) (*Movie, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangeVideo(ctx context.Context, videoMovie *VideoMovie) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) DeleteMovieByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) CreateTVShow(ctx context.Context, show *TVShow) (*TVShow, error) {
	panic("implement me")
}

func (am *AdminMicroservice) DeleteTVShowByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangePosters(ctx context.Context, dir *ContentPostersDir) (*Content, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangeContent(ctx context.Context, c *Content) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) CreateSeason(ctx context.Context, s *Season) (*Season, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangeSeason(ctx context.Context, s *Season) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) DeleteSeasonsByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) CreateEpisode(ctx context.Context, e *Episode) (*Episode, error) {
	panic("implement me")
}

func (am *AdminMicroservice) ChangeEpisode(ctx context.Context, e *Episode) (*empty.Empty, error) {
	panic("implement me")
}

func (am *AdminMicroservice) DeleteEpisodeByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	panic("implement me")
}
