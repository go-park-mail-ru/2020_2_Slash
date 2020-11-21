package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/search"
)

type SearchUsecase struct {
	// tvshowsRep tvshow.TvShowRepository
	actorsRep actor.ActorRepository
	moviesRep movie.MovieRepository
}

func NewSearchUsecase(actorsRep actor.ActorRepository,
	moviesRep movie.MovieRepository) search.SearchUsecase {
	return &SearchUsecase{
		actorsRep: actorsRep,
		moviesRep: moviesRep,
	}
}

func (uc SearchUsecase) Search(curUserID uint64, query string, pagination *models.Pagination) ([]*models.Movie, []*models.Actor, *errors.Error) {
	movies, err := uc.moviesRep.SelectWhereNameLike(curUserID, query, pagination.Count, pagination.From)
	if err == sql.ErrNoRows || (err == nil && movies == nil) {
		movies = []*models.Movie{}
	} else if err != nil {
		return nil, nil, errors.New(consts.CodeInternalError, err)
	}

	actors, err := uc.actorsRep.SelectWhereNameLike(query, pagination.Count, pagination.From)
	if err == sql.ErrNoRows || (err == nil && actors == nil) {
		actors = []*models.Actor{}
	} else if err != nil {
		return nil, nil, errors.New(consts.CodeInternalError, err)
	}

	return movies, actors, nil
}
