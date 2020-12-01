package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/jinzhu/copier"
)

type DirectorUseCase struct {
	directorRepo     director.DirectorRepository
	adminPanelClient admin.AdminPanelClient
}

func NewDirectorUseCase(repo director.DirectorRepository, client admin.AdminPanelClient) director.DirectorUseCase {
	return &DirectorUseCase{
		directorRepo:     repo,
		adminPanelClient: client,
	}
}

func (du *DirectorUseCase) Create(director *models.Director) *errors.Error {
	grpcDirector, err := du.adminPanelClient.CreateDirector(context.Background(),
		admin.DirectorModelToGRPC(director))
	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	if err := copier.Copy(director, admin.DirectorGRPCToModel(grpcDirector)); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (du *DirectorUseCase) Change(newDirector *models.Director) *errors.Error {
	_, err := du.adminPanelClient.ChangeDirector(context.Background(),
		admin.DirectorModelToGRPC(newDirector))

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (du *DirectorUseCase) DeleteById(id uint64) *errors.Error {
	_, err := du.adminPanelClient.DeleteDirectorByID(context.Background(),
		&admin.ID{ID: id})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
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
