package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/admin/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectorUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep, adminPanelClient)

	director := &models.Director{
		Name: "Sergio Leone",
	}

	grpcDirector := admin.DirectorModelToGRPC(director)
	adminPanelClient.
		EXPECT().
		CreateDirector(context.Background(), grpcDirector).
		Return(grpcDirector, nil)

	err := directorUseCase.Create(director)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestDirectorUseCase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep, adminPanelClient)

	director := &models.Director{
		ID:   3,
		Name: "Sergio Leone",
	}

	adminPanelClient.
		EXPECT().
		DeleteDirectorByID(context.Background(), &admin.ID{ID: director.ID}).
		Return(&empty.Empty{}, nil)

	err := directorUseCase.DeleteById(director.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestDirectorUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep, adminPanelClient)

	director := &models.Director{
		ID:   3,
		Name: "Sergio Leone",
	}

	grpcDirector := admin.DirectorModelToGRPC(director)
	adminPanelClient.
		EXPECT().
		ChangeDirector(context.Background(), grpcDirector).
		Return(&empty.Empty{}, nil)

	err := directorUseCase.Change(director)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestDirectorUseCase_ListByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := mocks.NewMockDirectorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	directorUseCase := NewDirectorUseCase(directorRep, adminPanelClient)

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
