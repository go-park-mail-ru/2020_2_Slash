package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"github.com/jinzhu/copier"
)

type MovieUsecase struct {
	movieRepo        movie.MovieRepository
	contentUcase     content.ContentUsecase
	adminPanelClient admin.AdminPanelClient
}

func NewMovieUsecase(repo movie.MovieRepository,
	contentUcase content.ContentUsecase, client admin.AdminPanelClient) movie.MovieUsecase {
	return &MovieUsecase{
		movieRepo:        repo,
		contentUcase:     contentUcase,
		adminPanelClient: client,
	}
}

func (mu *MovieUsecase) Create(movie *models.Movie) *errors.Error {
	grpcMovie, err := mu.adminPanelClient.CreateMovie(context.Background(),
		admin.MovieModelToGRPC(movie))
	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	if err := copier.Copy(movie, admin.MovieGRPCToModel(grpcMovie)); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (mu *MovieUsecase) UpdateVideo(movie *models.Movie, newVideoPath string) *errors.Error {
	_, err := mu.adminPanelClient.ChangeVideo(context.Background(),
		&admin.VideoMovie{
			Movie: admin.MovieModelToGRPC(movie),
			Video: newVideoPath,
		})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (mu *MovieUsecase) DeleteByID(movieID uint64) *errors.Error {
	_, err := mu.adminPanelClient.DeleteMovieByID(context.Background(),
		&admin.ID{ID: movieID})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (mu *MovieUsecase) GetByID(movieID uint64) (*models.Movie, *errors.Error) {
	movie, err := mu.movieRepo.SelectByID(movieID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeMovieDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return movie, nil
}

func (mu *MovieUsecase) GetFullByID(movieID uint64, curUserID uint64) (*models.Movie, *errors.Error) {
	movie, err := mu.movieRepo.SelectFullByID(movieID, curUserID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeMovieDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	customErr := mu.contentUcase.FillContent(&movie.Content)
	if customErr != nil {
		return nil, customErr
	}
	return movie, nil
}

func (mu *MovieUsecase) GetByContentID(contentID uint64) (*models.Movie, *errors.Error) {
	movie, err := mu.movieRepo.SelectByContentID(contentID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeMovieDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return movie, nil
}

func (mu *MovieUsecase) ListByParams(params *models.ContentFilter, pgnt *models.Pagination,
	curUserID uint64) ([]*models.Movie, *errors.Error) {

	movies, err := mu.movieRepo.SelectByParams(params, pgnt, curUserID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(movies) == 0 {
		return []*models.Movie{}, nil
	}
	return movies, nil
}

func (mu *MovieUsecase) ListLatest(pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, *errors.Error) {
	movies, err := mu.movieRepo.SelectLatest(pgnt, curUserID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(movies) == 0 {
		return []*models.Movie{}, nil
	}
	return movies, nil
}

func (mu *MovieUsecase) checkByContentID(contentID uint64) *errors.Error {
	_, err := mu.GetByContentID(contentID)
	return err
}

func (mu *MovieUsecase) ListByRating(pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, *errors.Error) {
	movies, err := mu.movieRepo.SelectByRating(pgnt, curUserID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}

	if len(movies) == 0 {
		return []*models.Movie{}, nil
	}

	return movies, nil
}
