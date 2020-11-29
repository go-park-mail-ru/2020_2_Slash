package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type FavouriteUsecase struct {
	favouriteRepo  favourite.FavouriteRepository
}

func NewFavouriteUsecase(repo favourite.FavouriteRepository) favourite.FavouriteUsecase {
	return &FavouriteUsecase{
		favouriteRepo: repo,
	}
}

func (uc *FavouriteUsecase) Create(favourite *models.Favourite) *errors.Error {
	isExist, customErr := uc.IsExist(favourite)
	if customErr != nil {
		return customErr
	}
	if isExist {
		return errors.Get(consts.CodeFavouriteAlreadyExist)
	}

	err := uc.favouriteRepo.Insert(favourite)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *FavouriteUsecase) GetUserFavourites(userID uint64,
	pagination *models.Pagination) (*models.FavouritesResult, *errors.Error) {
	favouriteMovies, err := uc.favouriteRepo.
		SelectFavouriteMovies(userID, pagination.Count, pagination.From)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	if len(favouriteMovies) == 0 {
		favouriteMovies = []*models.Movie{}
	}

	favouriteTVShows, err := uc.favouriteRepo.
		SelectFavouriteTVShows(userID, pagination.Count, pagination.From)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	if len(favouriteTVShows) == 0 {
		favouriteTVShows = []*models.TVShow{}
	}

	result := &models.FavouritesResult{
		Movies:  favouriteMovies,
		TVShows: favouriteTVShows,
	}

	return result, nil
}

func (uc *FavouriteUsecase) Delete(favourite *models.Favourite) *errors.Error {
	isExist, customErr := uc.IsExist(favourite)
	if customErr != nil {
		return customErr
	}
	if !isExist {
		return errors.Get(consts.CodeFavouriteDoesNotExist)
	}

	err := uc.favouriteRepo.Delete(favourite)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *FavouriteUsecase) IsExist(favourite *models.Favourite) (bool, *errors.Error) {
	err := uc.favouriteRepo.Select(favourite)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.New(consts.CodeInternalError, err)
	}
	return true, nil
}
