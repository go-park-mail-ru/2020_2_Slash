package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite/mocks"
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
	favouriteUseCase := NewFavouriteUsecase(favouriteRep)

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
	favouriteUseCase := NewFavouriteUsecase(favouriteRep)

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
	favouriteUseCase := NewFavouriteUsecase(favouriteRep)

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
	favouriteUseCase := NewFavouriteUsecase(favouriteRep)

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
	favouriteUseCase := NewFavouriteUsecase(favouriteRep)

	var userID uint64 = 3
	pagination := models.Pagination{
		From:  0,
		Count: 2,
	}

	movies := []*models.Movie{
		&models.Movie{
			ID: 2,
		},
		&models.Movie{
			ID: 4,
		},
	}

	tvShows := []*models.TVShow{
		&models.TVShow{
			ID: 1,
		},
		&models.TVShow{
			ID: 3,
		},
	}

	expectReturn := &models.FavouritesResult{
		Movies:  movies,
		TVShows: tvShows,
	}

	favouriteRep.
		EXPECT().
		SelectFavouriteMovies(gomock.Eq(userID), pagination.Count, pagination.From).
		Return(movies, nil)
	favouriteRep.
		EXPECT().
		SelectFavouriteTVShows(gomock.Eq(userID), pagination.Count, pagination.From).
		Return(tvShows, nil)

	res, err := favouriteUseCase.GetUserFavourites(userID, &pagination)
	assert.Equal(t, expectReturn, res)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestFavouriteUsecase_GetUserFavourites_Empty(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteRep := mocks.NewMockFavouriteRepository(ctrl)
	favouriteUseCase := NewFavouriteUsecase(favouriteRep)

	var userID uint64 = 3
	pagination := models.Pagination{
		From:  0,
		Count: 2,
	}

	expectReturn := &models.FavouritesResult{
		Movies:  []*models.Movie{},
		TVShows: []*models.TVShow{},
	}

	favouriteRep.
		EXPECT().
		SelectFavouriteMovies(gomock.Eq(userID), pagination.Count, pagination.From).
		Return(nil, sql.ErrNoRows)

	favouriteRep.
		EXPECT().
		SelectFavouriteTVShows(gomock.Eq(userID), pagination.Count, pagination.From).
		Return(nil, sql.ErrNoRows)

	res, err := favouriteUseCase.GetUserFavourites(userID, &pagination)
	assert.Equal(t, expectReturn, res)
	assert.Equal(t, (*errors.Error)(nil), err)
}
