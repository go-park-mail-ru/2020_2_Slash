package usecases

import (
	"database/sql"

	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type DirectorUseCase struct {
	directorRepo director.DirectorRepository
}

func NewDirectorUseCase(repo director.DirectorRepository) director.DirectorUseCase {
	return &DirectorUseCase{
		directorRepo: repo,
	}
}

func (du *DirectorUseCase) Create(director *models.Director) *errors.Error {
	err := du.directorRepo.Insert(director)
	if err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (du *DirectorUseCase) Get(id uint64) (*models.Director, *errors.Error) {
	dbDirector, err := du.directorRepo.SelectById(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(CodeDirectorDoesNotExist)
	} else if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return dbDirector, nil
}

func (du *DirectorUseCase) Change(newDirector *models.Director) *errors.Error {
	if _, customErr := du.Get(newDirector.ID); customErr != nil {
		return customErr
	}

	if err := du.directorRepo.Update(newDirector); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (du *DirectorUseCase) DeleteById(id uint64) *errors.Error {
	if _, customErr := du.Get(id); customErr != nil {
		return customErr
	}

	if err := du.directorRepo.DeleteById(id); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (du *DirectorUseCase) ListByID(directorsID []uint64) ([]*models.Director, *errors.Error) {
	var directors []*models.Director
	for _, directorID := range directorsID {
		director, err := du.Get(directorID)
		if err != nil {
			return nil, err
		}
		directors = append(directors, director)
	}
	return directors, nil
}

func (du *DirectorUseCase) List(pgnt *models.Pagination) ([]*models.Director, *errors.Error) {
	directors, err := du.directorRepo.SelectAll(pgnt)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(directors) == 0 {
		return []*models.Director{}, nil
	}
	return directors, nil
}
