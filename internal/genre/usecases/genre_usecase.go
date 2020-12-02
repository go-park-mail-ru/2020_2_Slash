package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/jinzhu/copier"
)

type GenreUsecase struct {
	genreRepo        genre.GenreRepository
	adminPanelClient admin.AdminPanelClient
}

func NewGenreUsecase(repo genre.GenreRepository,
	client admin.AdminPanelClient) genre.GenreUsecase {
	return &GenreUsecase{
		genreRepo:        repo,
		adminPanelClient: client,
	}
}

func (gu *GenreUsecase) Create(genre *models.Genre) *errors.Error {
	grpcGenre, err := gu.adminPanelClient.CreateGenre(context.Background(),
		admin.GenreModelToGRPC(genre))
	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	if err := copier.Copy(genre, admin.GenreGRPCToModel(grpcGenre)); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (gu *GenreUsecase) Update(newGenreData *models.Genre) *errors.Error {
	_, err := gu.adminPanelClient.ChangeGenre(context.Background(),
		admin.GenreModelToGRPC(newGenreData))

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (gu *GenreUsecase) DeleteByID(genreID uint64) *errors.Error {
	_, err := gu.adminPanelClient.DeleteGenreByID(context.Background(),
		&admin.ID{ID: genreID})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (gu *GenreUsecase) GetByID(genreID uint64) (*models.Genre, *errors.Error) {
	genre, err := gu.genreRepo.SelectByID(genreID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeGenreDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return genre, nil
}

func (gu *GenreUsecase) GetByName(name string) (*models.Genre, *errors.Error) {
	genre, err := gu.genreRepo.SelectByName(name)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeGenreDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return genre, nil
}

func (gu *GenreUsecase) List() ([]*models.Genre, *errors.Error) {
	genres, err := gu.genreRepo.SelectAll()
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(genres) == 0 {
		return []*models.Genre{}, nil
	}
	return genres, nil
}

func (gu *GenreUsecase) ListByID(genresID []uint64) ([]*models.Genre, *errors.Error) {
	var genres []*models.Genre
	for _, genreID := range genresID {
		genre, err := gu.GetByID(genreID)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (gu *GenreUsecase) checkByID(genreID uint64) *errors.Error {
	_, err := gu.GetByID(genreID)
	return err
}

func (gu *GenreUsecase) checkByName(name string) *errors.Error {
	_, err := gu.GetByName(name)
	return err
}
