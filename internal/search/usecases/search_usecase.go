package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/search"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
)

type SearchUsecase struct {
	tvshowsRep tvshow.TVShowRepository
	actorsRep  actor.ActorRepository
	moviesRep  movie.MovieRepository
}

func NewSearchUsecase(actorsRep actor.ActorRepository,
	moviesRep movie.MovieRepository,
	tvshowsRep tvshow.TVShowRepository) search.SearchUsecase {
	return &SearchUsecase{
		tvshowsRep: tvshowsRep,
		actorsRep:  actorsRep,
		moviesRep:  moviesRep,
	}
}

func (uc SearchUsecase) Search(curUserID uint64, query string,
	pagination *models.Pagination) (*models.SearchResult, *errors.Error) {
	movies, err := uc.moviesRep.SelectWhereNameLike(curUserID, query, pagination.Count, pagination.From)
	if err == sql.ErrNoRows || (err == nil && movies == nil) {
		movies = []*models.Movie{}
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	actors, err := uc.actorsRep.SelectWhereNameLike(query, pagination.Count, pagination.From)
	if err == sql.ErrNoRows || (err == nil && actors == nil) {
		actors = []*models.Actor{}
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	tvShows, err := uc.tvshowsRep.SelectWhereNameLike(query, pagination, curUserID)
	if err == sql.ErrNoRows || (err == nil && tvShows == nil) {
		tvShows = []*models.TVShow{}
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	result := &models.SearchResult{
		TVShows: tvShows,
		Movies:  movies,
		Actors:  actors,
	}

	return result, nil
}
