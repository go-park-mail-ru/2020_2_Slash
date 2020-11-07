package usecases

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"os"
	"path/filepath"
)

type MovieUsecase struct {
	movieRepo    movie.MovieRepository
	contentUcase content.ContentUsecase
}

func NewMovieUsecase(repo movie.MovieRepository,
	contentUcase content.ContentUsecase) movie.MovieUsecase {
	return &MovieUsecase{
		movieRepo:    repo,
		contentUcase: contentUcase,
	}
}

func (mu *MovieUsecase) Create(movie *models.Movie) *errors.Error {
	if err := mu.checkByContentID(movie.ContentID); err == nil {
		return errors.Get(CodeMovieContentAlreadyExists)
	}

	if err := mu.contentUcase.Create(&movie.Content); err != nil {
		return err
	}

	if err := mu.movieRepo.Insert(movie); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (mu *MovieUsecase) UpdateVideo(movie *models.Movie, newVideoPath string) *errors.Error {
	prevVideoPath := movie.Video
	if newVideoPath == prevVideoPath {
		// Don't need to update
		return nil
	}

	// Update video
	movie.Video = newVideoPath
	if err := mu.movieRepo.Update(movie); err != nil {
		return errors.New(CodeInternalError, err)
	}
	// Don't need to delete prev file,
	// cause video always store with the same filename
	return nil
}

func (mu *MovieUsecase) DeleteByID(movieID uint64) *errors.Error {
	movie, err := mu.GetByID(movieID)
	if err != nil {
		return errors.Get(CodeMovieDoesNotExist)
	}

	// Delete video
	if movie.Video != "" {
		path, err := os.Getwd()
		if err != nil {
			return errors.New(CodeInternalError, err)
		}
		videoPath := filepath.Join(path, movie.Video)
		videoDirPath := filepath.Dir(videoPath)

		if err := os.RemoveAll(videoDirPath); err != nil {
			return errors.New(CodeInternalError, err)
		}
	}

	if err := mu.movieRepo.DeleteByID(movieID); err != nil {
		return errors.New(CodeInternalError, err)
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

func (mu *MovieUsecase) GetWithContentByID(movieID uint64) (*models.Movie, *errors.Error) {
	movie, err := mu.GetByID(movieID)
	if err != nil {
		return nil, err
	}
	content, err := mu.contentUcase.GetByID(movie.ContentID)
	if err != nil {
		return nil, err
	}
	movie.Content = *content
	return movie, nil
}

func (mu *MovieUsecase) GetFullByID(movieID uint64) (*models.Movie, *errors.Error) {
	movie, err := mu.GetByID(movieID)
	if err != nil {
		return nil, err
	}
	content, err := mu.contentUcase.GetFullByID(movie.ContentID)
	if err != nil {
		return nil, err
	}
	movie.Content = *content
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

func (mu *MovieUsecase) ListByParams(params *models.ContentFilter, pgnt *models.Pagination) ([]*models.Movie, *errors.Error) {
	movies, err := mu.movieRepo.SelectByParams(params, pgnt)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(movies) == 0 {
		return []*models.Movie{}, nil
	}
	return movies, nil
}

func (mu *MovieUsecase) ListLatest(pgnt *models.Pagination) ([]*models.Movie, *errors.Error) {
	movies, err := mu.movieRepo.SelectLatest(pgnt)
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
