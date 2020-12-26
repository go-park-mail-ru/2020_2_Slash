package usecases

import (
	"database/sql"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var contentInst *models.Content = &models.Content{
	ContentID: 1,
}

var tvshowInst *models.TVShow = &models.TVShow{
	ID:      1,
	Content: *contentInst,
}

func TestTVShowUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	tvshowRep.
		EXPECT().
		SelectByContentID(gomock.Eq(tvshowInst.ContentID)).
		Return(nil, sql.ErrNoRows)

	contentUseCase.
		EXPECT().
		Create(gomock.Eq(contentInst)).
		Return(nil)

	tvshowRep.
		EXPECT().
		Insert(gomock.Eq(tvshowInst)).
		Return(nil)

	err := tvshowUseCase.Create(tvshowInst)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestTVShowUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	tvshowRep.
		EXPECT().
		SelectByContentID(gomock.Eq(tvshowInst.ContentID)).
		Return(tvshowInst, nil)

	err := tvshowUseCase.Create(tvshowInst)
	assert.Equal(t, err, errors.Get(consts.CodeTVShowContentAlreadyExists))
}

func TestTVShowUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	tvshowRep.
		EXPECT().
		SelectByID(gomock.Eq(tvshowInst.ID)).
		Return(tvshowInst, nil)

	dbTVShow, err := tvshowUseCase.GetByID(tvshowInst.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbTVShow, tvshowInst)
}

func TestTVShowUseCase_GetByID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	tvshowRep.
		EXPECT().
		SelectByID(gomock.Eq(tvshowInst.ID)).
		Return(nil, sql.ErrNoRows)

	dbTVShow, err := tvshowUseCase.GetByID(tvshowInst.ID)
	assert.Equal(t, err, errors.Get(consts.CodeTVShowDoesNotExist))
	assert.Equal(t, dbTVShow, (*models.TVShow)(nil))
}

func TestTVShowUseCase_GetFullByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)
	var userID uint64 = 1

	tvshowRep.
		EXPECT().
		SelectFullByID(gomock.Eq(tvshowInst.ID), gomock.Eq(userID)).
		Return(tvshowInst, nil)

	contentUseCase.
		EXPECT().
		FillContent(gomock.Eq(&tvshowInst.Content)).
		Return(nil)

	dbTVShow, err := tvshowUseCase.GetFullByID(tvshowInst.ID, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbTVShow, tvshowInst)
}

func TestTVShowUseCase_GetByContentID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	tvshowRep.
		EXPECT().
		SelectByContentID(gomock.Eq(tvshowInst.ContentID)).
		Return(nil, sql.ErrNoRows)

	dbTVShow, err := tvshowUseCase.GetByContentID(tvshowInst.ContentID)
	assert.Equal(t, err, errors.Get(consts.CodeTVShowDoesNotExist))
	assert.Equal(t, dbTVShow, (*models.TVShow)(nil))
}

func TestTVShowUseCase_ListByParams_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	var contentInst *models.Content = &models.Content{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		Countries:        nil,
		Genres:           nil,
		Actors:           nil,
		Directors:        nil,
		Type:             "tvshow",
	}

	content := []*models.Content{
		contentInst,
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	params := &models.ContentFilter{
		Year:     []int{2001},
		Genre:    []int{1},
		Country:  []int{1},
		Actor:    []int{1},
		Director: []int{1},
	}

	tvshowRep.
		EXPECT().
		SelectByParams(gomock.Eq(params), gomock.Eq(pgnt), gomock.Eq(userID)).
		Return(tvshows, nil)

	dbTVShows, err := tvshowUseCase.ListByParams(params, pgnt, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbTVShows, tvshows)
}

func TestTVShowUseCase_ListLatest_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	tvshowRep.
		EXPECT().
		SelectLatest(gomock.Eq(pgnt), gomock.Eq(userID)).
		Return(tvshows, nil)

	dbTVShows, err := tvshowUseCase.ListLatest(pgnt, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbTVShows, tvshows)
}

func TestTVShowUseCase_ListByRating_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tvshowRep := mocks.NewMockTVShowRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	tvshowUseCase := NewTVShowUsecase(tvshowRep, contentUseCase)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	tvshowRep.
		EXPECT().
		SelectByRating(gomock.Eq(pgnt), gomock.Eq(userID)).
		Return(tvshows, nil)

	dbTVShows, err := tvshowUseCase.ListByRating(pgnt, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbTVShows, tvshows)
}
