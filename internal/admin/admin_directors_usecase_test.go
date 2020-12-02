package admin

import (
	"context"
	directorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectorUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := directorMocks.NewMockDirectorRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		directorsRep: directorRep,
	}

	director := &Director{
		Name: "Jamie Fox",
	}

	directorRep.
		EXPECT().
		Insert(DirectorGRPCToModel(director)).
		Return(nil)

	_, err := adminMicroservice.CreateDirector(context.Background(), director)
	assert.Equal(t, err, (error)(nil))
}

func TestDirectorUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := directorMocks.NewMockDirectorRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		directorsRep: directorRep,
	}

	director := &models.Director{
		ID:   3,
		Name: "Margo Robbie",
	}

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(director.ID)).
		Return(director, nil)

	directorRep.
		EXPECT().
		DeleteById(gomock.Eq(director.ID)).
		Return(nil)

	_, err := adminMicroservice.DeleteDirectorByID(context.Background(),
		&ID{ID: director.ID})
	assert.Equal(t, err, (error)(nil))
}

func TestDirectorUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorRep := directorMocks.NewMockDirectorRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		directorsRep: directorRep,
	}
	director := &models.Director{
		ID:   3,
		Name: "Margo Robbie",
	}

	directorRep.
		EXPECT().
		SelectById(gomock.Eq(director.ID)).
		Return(director, nil)

	directorRep.
		EXPECT().
		Update(gomock.Eq(director)).
		Return(nil)

	_, err := adminMicroservice.ChangeDirector(context.Background(), DirectorModelToGRPC(director))
	assert.Equal(t, err, (error)(nil))
}

