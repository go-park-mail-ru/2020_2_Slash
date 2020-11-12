package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/rating/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRatingUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(rating.UserID, rating.ContentID).
		Return(nil, sql.ErrNoRows)

	ratingRep.
		EXPECT().
		Insert(rating).
		Return(nil)

	err := ratingUseCase.Create(rating)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestRatingUseCase_Create_ContentDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, errors.Get(consts.CodeContentDoesNotExist))

	err := ratingUseCase.Create(rating)
	assert.Equal(t, errors.Get(consts.CodeContentDoesNotExist), err)
}

func TestRatingUseCase_Create_AlreadyExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(rating.UserID, rating.ContentID).
		Return(rating, nil)

	err := ratingUseCase.Create(rating)
	assert.Equal(t, errors.Get(consts.CodeRatingAlreadyExist), err)
}

func TestRatingUsecase_Change_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	existedRating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}
	updateRating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     false,
	}

	contentUseCase.
		EXPECT().
		GetByID(existedRating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(existedRating.UserID, existedRating.ContentID).
		Return(existedRating, nil)

	ratingRep.
		EXPECT().
		Update(updateRating).
		Return(nil)

	err := ratingUseCase.Change(updateRating)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestRatingUsecase_Change_ContentDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, errors.Get(consts.CodeContentDoesNotExist))

	err := ratingUseCase.Change(rating)
	assert.Equal(t, errors.Get(consts.CodeContentDoesNotExist), err)
}

func TestRatingUsecase_Change_LikesEqualToDb(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	existedRating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(existedRating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(existedRating.UserID, existedRating.ContentID).
		Return(existedRating, nil)

	err := ratingUseCase.Change(existedRating)
	assert.Equal(t, errors.Get(consts.CodeRatingAlreadyExist), err)
}

func TestRatingUsecase_Change_RatingDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(rating.UserID, rating.ContentID).
		Return(nil, nil)

	err := ratingUseCase.Change(rating)
	assert.Equal(t, errors.Get(consts.CodeRatingDoesNotExist), err)
}

func TestRatingUsecase_Delete_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(rating.UserID, rating.ContentID).
		Return(rating, nil)

	ratingRep.
		EXPECT().
		Delete(rating).
		Return(nil)

	err := ratingUseCase.Delete(rating)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestRatingUsecase_Delete_RatingDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(rating.UserID, rating.ContentID).
		Return(nil, nil)

	err := ratingUseCase.Delete(rating)
	assert.Equal(t, errors.Get(consts.CodeRatingDoesNotExist), err)
}

func TestRatingUsecase_Delete_ContentDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, errors.Get(consts.CodeContentDoesNotExist))

	err := ratingUseCase.Delete(rating)
	assert.Equal(t, errors.Get(consts.CodeContentDoesNotExist), err)
}

func TestRatingUsecase_GetByUserIDContentID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(rating.UserID, rating.ContentID).
		Return(rating, nil)

	dbRating, err := ratingUseCase.GetByUserIDContentID(rating.UserID, rating.ContentID)
	assert.Equal(t, rating, dbRating)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestRatingUsecase_GetByUserIDContentID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingRep.
		EXPECT().
		SelectByUserIDContentID(rating.UserID, rating.ContentID).
		Return(nil, sql.ErrNoRows)

	dbRating, err := ratingUseCase.GetByUserIDContentID(rating.UserID, rating.ContentID)
	assert.Equal(t, (*models.Rating)(nil), dbRating)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestRatingUsecase_GetContentRating_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectLikesCount(rating.ContentID).
		Return(73, nil)

	ratingRep.
		EXPECT().
		SelectRatesCount(rating.ContentID).
		Return(100, nil)

	contentRating, err := ratingUseCase.GetContentRating(rating.ContentID)
	assert.Equal(t, 73, contentRating)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestRatingUsecase_GetContentRating_Round(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, nil)

	ratingRep.
		EXPECT().
		SelectLikesCount(rating.ContentID).
		Return(39, nil)

	ratingRep.
		EXPECT().
		SelectRatesCount(rating.ContentID).
		Return(55, nil)

	contentRating, err := ratingUseCase.GetContentRating(rating.ContentID)
	assert.Equal(t, 71, contentRating)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestRatingUsecase_GetContentRating_ContentDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingRep := mocks.NewMockRatingRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	ratingUseCase := NewRatingUseCase(ratingRep, contentUseCase)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	contentUseCase.
		EXPECT().
		GetByID(rating.ContentID).
		Return(nil, errors.Get(consts.CodeContentDoesNotExist))

	_, err := ratingUseCase.GetContentRating(rating.ContentID)
	assert.Equal(t, errors.Get(consts.CodeContentDoesNotExist), err)
}
