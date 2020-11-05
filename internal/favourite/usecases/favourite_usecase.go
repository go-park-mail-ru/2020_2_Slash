package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type FavouriteUsecase struct {
	favouriteRepo  favourite.FavouriteRepository
	contentUseCase content.ContentUsecase
}

func NewFavouriteUsecase(repo favourite.FavouriteRepository, contentUseCase content.ContentUsecase) favourite.FavouriteUsecase {
	return &FavouriteUsecase{
		favouriteRepo: repo,
		contentUseCase : contentUseCase,
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

func (uc *FavouriteUsecase) GetUserFavourites(userID uint64) ([]*models.Content, *errors.Error) {
	favourites, err := uc.favouriteRepo.SelectFavouriteContent(userID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	if len(favourites) == 0 {
		return []*models.Content{}, nil
	}

	return favourites, nil
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
