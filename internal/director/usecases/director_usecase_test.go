package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectorUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep)

	director := &models.Director{
		Name: "Sergio Leone",
	}

	directorRep.
		EXPECT().
		Insert(gomock.Eq(director)).
		Return(nil)

	err := directorUseCase.Create(director)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestDirectorUseCase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep)

	director := &models.Director{
		ID:   3,
		Name: "Sergio Leone",
	}

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(director.ID)).
		Return(director, nil)

	dbDirector, err := directorUseCase.Get(director.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbDirector, director)
}

func TestDirectorUseCase_Get_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep)

	director := &models.Director{
		ID:   3,
		Name: "Sergio Leone",
	}

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(director.ID)).
		Return(nil, sql.ErrNoRows)

	dbDirector, err := directorUseCase.Get(director.ID)
	assert.Equal(t, err, errors.Get(consts.CodeDirectorDoesNotExist))
	assert.Equal(t, dbDirector, (*models.Director)(nil))
}

func TestDirectorUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep)

	director := &models.Director{
		ID:   3,
		Name: "Sergio Leone",
	}

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(director.ID)).
		Return(director, nil)

	directorRep.
		EXPECT().
		DeleteById(gomock.Eq(director.ID)).
		Return(nil)

	err := directorUseCase.DeleteById(director.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestDirectorUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep)

	director := &models.Director{
		ID:   3,
		Name: "Sergio Leone",
	}

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(director.ID)).
		Return(director, nil)

	directorRep.
		EXPECT().
		Update(gomock.Eq(director)).
		Return(nil)

	err := directorUseCase.Change(director)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestDirectorUseCase_ListByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep)

	directors := []*models.Director{
		&models.Director{
			ID:   1,
			Name: "Sergio Leone",
		},
		&models.Director{
			ID:   2,
			Name: "No Sergio Leone",
		},
	}

	directorsID := []uint64{1, 2}

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(directorsID[0])).
		Return(directors[0], nil)

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(directorsID[1])).
		Return(directors[1], nil)

	dbDirectors, err := directorUseCase.ListByID(directorsID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbDirectors, directors)
}
