package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	contentUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/rating"
	"math"
)

type RatingUsecase struct {
	rep rating.RatingRepository
	contentUseCase contentUsecase.ContentUsecase
}

func NewRatingUseCase(rep rating.RatingRepository,
	contentUseCase contentUsecase.ContentUsecase) rating.RatingUsecase {
	return &RatingUsecase{rep: rep, contentUseCase: contentUseCase}
}

func (uc *RatingUsecase) Create(rating *models.Rating) *errors.Error {
	_, customErr := uc.contentUseCase.GetByID(rating.ContentID)
	if customErr != nil {
		return customErr
	}

	isExist, customErr := uc.isExist(rating)
	if customErr != nil {
		return customErr
	}
	if isExist {
		return errors.Get(consts.CodeRatingAlreadyExist)
	}

	err := uc.rep.Insert(rating)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *RatingUsecase) Change(rating *models.Rating) *errors.Error {
	_, customErr := uc.contentUseCase.GetByID(rating.ContentID)
	if customErr != nil {
		return customErr
	}

	dbRating, customErr := uc.GetByUserIDContentID(rating.UserID, rating.ContentID)
	if customErr != nil {
		return customErr
	}

	if dbRating == nil {
		return errors.Get(consts.CodeRatingDoesNotExist)
	} else if dbRating.Likes == rating.Likes {
		return errors.Get(consts.CodeRatingAlreadyExist)
	}

	err := uc.rep.Update(rating)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *RatingUsecase) GetContentRating(contentID uint64) (int, *errors.Error) {
	_, customErr := uc.contentUseCase.GetByID(contentID)
	if customErr != nil {
		return 0, customErr
	}

	numberOfRates, err := uc.rep.SelectRatesCount(contentID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, errors.New(consts.CodeInternalError, err)
	}

	numberOfLikes, err := uc.rep.SelectLikesCount(contentID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, errors.New(consts.CodeInternalError, err)
	}

	rate := 100 * float64(numberOfLikes) / float64(numberOfRates)
	percentage := int(math.Round(rate))
	return percentage, nil
}

func (uc *RatingUsecase) Delete(rating *models.Rating) *errors.Error {
	_, customErr := uc.contentUseCase.GetByID(rating.ContentID)
	if customErr != nil {
		return customErr
	}

	isExist, customErr := uc.isExist(rating)
	if customErr != nil {
		return customErr
	}
	if !isExist {
		return errors.Get(consts.CodeRatingDoesNotExist)
	}

	err := uc.rep.Delete(rating)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *RatingUsecase) isExist(rating *models.Rating) (bool, *errors.Error) {
	rating, err := uc.GetByUserIDContentID(rating.UserID, rating.ContentID)
	if err != nil {
		return false, err
	}
	if rating != nil {
		return true, nil
	}
	return false, nil
}

func (uc *RatingUsecase) GetByUserIDContentID(userID uint64, contentID uint64) (*models.Rating, *errors.Error) {
	rating, err := uc.rep.SelectByUserIDContentID(userID, contentID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	return rating, nil
}
