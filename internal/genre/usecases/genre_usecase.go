package usecases

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type GenreUsecase struct {
	genreRepo genre.GenreRepository
}

func NewGenreUsecase(repo genre.GenreRepository) genre.GenreUsecase {
	return &GenreUsecase{
		genreRepo: repo,
	}
}

func (gu *GenreUsecase) Create(genre *models.Genre) *errors.Error {
	if err := gu.checkByName(genre.Name); err == nil {
		return errors.Get(CodeGenreNameAlreadyExists)
	}

	if err := gu.genreRepo.Insert(genre); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (gu *GenreUsecase) UpdateByID(genreID uint64, newGenreData *models.Genre) (*models.Genre, *errors.Error) {
	genre, err := gu.GetByID(genreID)
	if err != nil {
		return nil, err
	}

	if newGenreData.Name == "" || genre.Name == newGenreData.Name {
		return nil, err
	}

	if err := gu.checkByName(newGenreData.Name); err == nil {
		return nil, errors.Get(CodeGenreNameAlreadyExists)
	}

	genre.Name = newGenreData.Name
	if err := gu.genreRepo.Update(genre); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return genre, nil
}

func (gu *GenreUsecase) DeleteByID(genreID uint64) *errors.Error {
	if err := gu.checkByID(genreID); err != nil {
		return errors.Get(CodeGenreDoesNotExist)
	}

	if err := gu.genreRepo.DeleteByID(genreID); err != nil {
		return errors.New(CodeInternalError, err)
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
