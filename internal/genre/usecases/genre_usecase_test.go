package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenreUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(genre.Name)).
		Return(nil, sql.ErrNoRows)

	genreRep.
		EXPECT().
		Insert(gomock.Eq(genre)).
		Return(nil)

	err := genreUseCase.Create(genre)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestGenreUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(genre.Name)).
		Return(genre, nil)

	err := genreUseCase.Create(genre)
	assert.Equal(t, err, errors.Get(consts.CodeGenreNameAlreadyExists))
}

func TestGenreUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	newGenreData := &models.Genre{
		ID:   1,
		Name: "GB",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(genre, nil)

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(newGenreData.Name)).
		Return(nil, sql.ErrNoRows)

	genreRep.
		EXPECT().
		Update(gomock.Eq(genre)).
		Return(nil)

	dbGenre, err := genreUseCase.UpdateByID(genre.ID, newGenreData)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenre, newGenreData)
}

func TestGenreUseCase_Update_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	newGenreData := &models.Genre{
		ID:   1,
		Name: "GB",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(genre, nil)

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(newGenreData.Name)).
		Return(newGenreData, nil)

	dbGenre, err := genreUseCase.UpdateByID(genre.ID, newGenreData)
	assert.Equal(t, err, errors.Get(consts.CodeGenreNameAlreadyExists))
	assert.Equal(t, dbGenre, (*models.Genre)(nil))
}

func TestGenreUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(genre, nil)

	genreRep.
		EXPECT().
		DeleteByID(gomock.Eq(genre.ID)).
		Return(nil)

	err := genreUseCase.DeleteByID(genre.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestGenreUseCase_Delete_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(nil, sql.ErrNoRows)

	err := genreUseCase.DeleteByID(genre.ID)
	assert.Equal(t, err, errors.Get(consts.CodeGenreDoesNotExist))
}

func TestGenreUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(genre, nil)

	dbGenre, err := genreUseCase.GetByID(genre.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenre, genre)
}

func TestGenreUseCase_GetByID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(nil, sql.ErrNoRows)

	dbGenre, err := genreUseCase.GetByID(genre.ID)
	assert.Equal(t, err, errors.Get(consts.CodeGenreDoesNotExist))
	assert.Equal(t, dbGenre, (*models.Genre)(nil))
}

func TestGenreUseCase_GetByName_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(genre.Name)).
		Return(genre, nil)

	dbGenre, err := genreUseCase.GetByName(genre.Name)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenre, genre)
}

func TestGenreUseCase_GetByName_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(genre.Name)).
		Return(nil, sql.ErrNoRows)

	dbGenre, err := genreUseCase.GetByName(genre.Name)
	assert.Equal(t, err, errors.Get(consts.CodeGenreDoesNotExist))
	assert.Equal(t, dbGenre, (*models.Genre)(nil))
}

func TestGenreUseCase_List_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genres := []*models.Genre{
		&models.Genre{
			ID:   1,
			Name: "comedy",
		},
		&models.Genre{
			ID:   2,
			Name: "mult",
		},
	}

	genreRep.
		EXPECT().
		SelectAll().
		Return(genres, nil)

	dbGenres, err := genreUseCase.List()
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenres, genres)
}

func TestGenreUseCase_ListByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	genreUseCase := NewGenreUsecase(genreRep)

	genres := []*models.Genre{
		&models.Genre{
			ID:   1,
			Name: "comedy",
		},
		&models.Genre{
			ID:   2,
			Name: "mult",
		},
	}

	genresID := []uint64{1, 2}

	genreRep.
		EXPECT().
		SelectByID(genresID[0]).
		Return(genres[0], nil)

	genreRep.
		EXPECT().
		SelectByID(genresID[1]).
		Return(genres[1], nil)

	dbGenres, err := genreUseCase.ListByID(genresID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenres, genres)
}
