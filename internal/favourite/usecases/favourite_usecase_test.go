package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite/mocks"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFavouriteUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteRep := mocks.NewMockFavouriteRepository(ctrl)
	contentUsecase := contentMocks.NewMockContentUsecase(ctrl)
	favouriteUseCase := NewFavouriteUsecase(favouriteRep, contentUsecase)

	favourite := &models.Favourite{
		UserID:    3,
		ContentID: 3,
		Created:   time.Now(),
	}

	favouriteRep.
		EXPECT().
		Select(gomock.Eq(favourite)).
		Return(sql.ErrNoRows)

	favouriteRep.
		EXPECT().
		Insert(gomock.Eq(favourite)).
		Return(nil)

	err := favouriteUseCase.Create(favourite)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestFavouriteUseCase_Create_AlreadyExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteRep := mocks.NewMockFavouriteRepository(ctrl)
	contentUsecase := contentMocks.NewMockContentUsecase(ctrl)
	favouriteUseCase := NewFavouriteUsecase(favouriteRep, contentUsecase)

	favourite := &models.Favourite{
		UserID:    3,
		ContentID: 3,
		Created:   time.Now(),
	}

	favouriteRep.
		EXPECT().
		Select(gomock.Eq(favourite)).
		Return(nil)

	err := favouriteUseCase.Create(favourite)
	assert.Equal(t, errors.Get(consts.CodeFavouriteAlreadyExist), err)
}

func TestFavouriteUseCase_Delete_DoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteRep := mocks.NewMockFavouriteRepository(ctrl)
	contentUsecase := contentMocks.NewMockContentUsecase(ctrl)
	favouriteUseCase := NewFavouriteUsecase(favouriteRep, contentUsecase)

	favourite := &models.Favourite{
		UserID:    3,
		ContentID: 3,
		Created:   time.Now(),
	}

	favouriteRep.
		EXPECT().
		Select(gomock.Eq(favourite)).
		Return(sql.ErrNoRows)

	err := favouriteUseCase.Delete(favourite)
	assert.Equal(t, errors.Get(consts.CodeFavouriteDoesNotExist), err)
}

func TestFavouriteUseCase_Delete_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteRep := mocks.NewMockFavouriteRepository(ctrl)
	contentUsecase := contentMocks.NewMockContentUsecase(ctrl)
	favouriteUseCase := NewFavouriteUsecase(favouriteRep, contentUsecase)

	favourite := &models.Favourite{
		UserID:    3,
		ContentID: 3,
		Created:   time.Now(),
	}

	favouriteRep.
		EXPECT().
		Select(gomock.Eq(favourite)).
		Return(nil)

	favouriteRep.
		EXPECT().
		Delete(gomock.Eq(favourite)).
		Return(nil)

	err := favouriteUseCase.Delete(favourite)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestFavouriteUsecase_GetUserFavourites_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteRep := mocks.NewMockFavouriteRepository(ctrl)
	contentUsecase := contentMocks.NewMockContentUsecase(ctrl)
	favouriteUseCase := NewFavouriteUsecase(favouriteRep, contentUsecase)

	var userID uint64 = 3

	expectReturn := []*models.Content{
		{
			ContentID: 2,
		},
		{
			ContentID: 4,
		},
	}

	favouriteRep.
		EXPECT().
		SelectFavouriteContent(gomock.Eq(userID)).
		Return(expectReturn, nil)

	res, err := favouriteUseCase.GetUserFavourites(userID)
	assert.Equal(t, expectReturn, res)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestFavouriteUsecase_GetUserFavourites_Empty(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteRep := mocks.NewMockFavouriteRepository(ctrl)
	contentUsecase := contentMocks.NewMockContentUsecase(ctrl)
	favouriteUseCase := NewFavouriteUsecase(favouriteRep, contentUsecase)

	var userID uint64 = 3

	var expectReturn []*models.Content

	favouriteRep.
		EXPECT().
		SelectFavouriteContent(gomock.Eq(userID)).
		Return(nil, sql.ErrNoRows)

	res, err := favouriteUseCase.GetUserFavourites(userID)
	assert.Equal(t, expectReturn, res)
	assert.Equal(t, (*errors.Error)(nil), err)
}
