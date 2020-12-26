package usecases

import (
	"database/sql"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var contentInst *models.Content = &models.Content{
	ContentID: 1,
}

var movieInst *models.Movie = &models.Movie{
	ID:      1,
	Video:   "movieInst.mp4",
	Content: *contentInst,
}

func TestMovieUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByContentID(gomock.Eq(movieInst.ContentID)).
		Return(nil, sql.ErrNoRows)

	contentUseCase.
		EXPECT().
		Create(gomock.Eq(contentInst)).
		Return(nil)

	movieRep.
		EXPECT().
		Insert(gomock.Eq(movieInst)).
		Return(nil)

	err := movieUseCase.Create(movieInst)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestMovieUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByContentID(gomock.Eq(movieInst.ContentID)).
		Return(movieInst, nil)

	err := movieUseCase.Create(movieInst)
	assert.Equal(t, err, errors.Get(consts.CodeMovieContentAlreadyExists))
}

func TestMovieUseCase_UpdateVideo_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	newVideoPath := "video/movie.mp4"

	movieRep.
		EXPECT().
		Update(gomock.Eq(movieInst)).
		Return(nil)

	err := movieUseCase.UpdateVideo(movieInst, newVideoPath)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestMovieUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movieInst.ID)).
		Return(movieInst, nil)

	dbMovie, err := movieUseCase.GetByID(movieInst.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovie, movieInst)
}

func TestMovieUseCase_GetByID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movieInst.ID)).
		Return(nil, sql.ErrNoRows)

	dbMovie, err := movieUseCase.GetByID(movieInst.ID)
	assert.Equal(t, err, errors.Get(consts.CodeMovieDoesNotExist))
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
}

func TestMovieUseCase_GetFullByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)
	var userID uint64 = 1

	movieRep.
		EXPECT().
		SelectFullByID(gomock.Eq(movieInst.ID), gomock.Eq(userID)).
		Return(movieInst, nil)

	contentUseCase.
		EXPECT().
		FillContent(gomock.Eq(&movieInst.Content)).
		Return(nil)

	dbMovie, err := movieUseCase.GetFullByID(movieInst.ID, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovie, movieInst)
}

func TestMovieUseCase_GetByContentID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByContentID(gomock.Eq(movieInst.ContentID)).
		Return(nil, sql.ErrNoRows)

	dbMovie, err := movieUseCase.GetByContentID(movieInst.ContentID)
	assert.Equal(t, err, errors.Get(consts.CodeMovieDoesNotExist))
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
}

func TestMovieUseCase_ListByParams_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

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
		Type:             "movie",
	}

	content := []*models.Content{
		contentInst,
	}

	movies := []*models.Movie{
		&models.Movie{
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

	movieRep.
		EXPECT().
		SelectByParams(gomock.Eq(params), gomock.Eq(pgnt), gomock.Eq(userID)).
		Return(movies, nil)

	dbMovies, err := movieUseCase.ListByParams(params, pgnt, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovies, movies)
}

func TestMovieUseCase_ListLatest_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	movies := []*models.Movie{
		&models.Movie{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	movieRep.
		EXPECT().
		SelectLatest(gomock.Eq(pgnt), gomock.Eq(userID)).
		Return(movies, nil)

	dbMovies, err := movieUseCase.ListLatest(pgnt, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovies, movies)
}

func TestMovieUseCase_ListByRating_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	movies := []*models.Movie{
		&models.Movie{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	movieRep.
		EXPECT().
		SelectByRating(gomock.Eq(pgnt), gomock.Eq(userID)).
		Return(movies, nil)

	dbMovies, err := movieUseCase.ListByRating(pgnt, userID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovies, movies)
}
