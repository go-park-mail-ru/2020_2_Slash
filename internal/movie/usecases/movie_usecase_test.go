package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var content_inst *models.Content = &models.Content{
	ContentID: 1,
}

var movie_inst *models.Movie = &models.Movie{
	ID:      1,
	Video:   "movie_inst.mp4",
	Content: *content_inst,
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
		SelectByContentID(gomock.Eq(movie_inst.ContentID)).
		Return(nil, sql.ErrNoRows)

	contentUseCase.
		EXPECT().
		Create(gomock.Eq(content_inst)).
		Return(nil)

	movieRep.
		EXPECT().
		Insert(gomock.Eq(movie_inst)).
		Return(nil)

	err := movieUseCase.Create(movie_inst)
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
		SelectByContentID(gomock.Eq(movie_inst.ContentID)).
		Return(movie_inst, nil)

	err := movieUseCase.Create(movie_inst)
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
		Update(gomock.Eq(movie_inst)).
		Return(nil)

	err := movieUseCase.UpdateVideo(movie_inst, newVideoPath)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestMovieUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(movie_inst, nil)

	movieRep.
		EXPECT().
		DeleteByID(gomock.Eq(movie_inst.ID)).
		Return(nil)

	err := movieUseCase.DeleteByID(movie_inst.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestMovieUseCase_Delete_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(nil, sql.ErrNoRows)

	err := movieUseCase.DeleteByID(movie_inst.ID)
	assert.Equal(t, err, errors.Get(consts.CodeMovieDoesNotExist))
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
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(movie_inst, nil)

	dbMovie, err := movieUseCase.GetByID(movie_inst.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovie, movie_inst)
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
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(nil, sql.ErrNoRows)

	dbMovie, err := movieUseCase.GetByID(movie_inst.ID)
	assert.Equal(t, err, errors.Get(consts.CodeMovieDoesNotExist))
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
}

func TestMovieUseCase_GetWithContentByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(movie_inst, nil)

	contentUseCase.
		EXPECT().
		GetByID(gomock.Eq(movie_inst.ContentID)).
		Return(content_inst, nil)

	dbMovie, err := movieUseCase.GetWithContentByID(movie_inst.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovie, movie_inst)
}

func TestMovieUseCase_GetWithContentByID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(nil, sql.ErrNoRows)

	dbMovie, err := movieUseCase.GetWithContentByID(movie_inst.ID)
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

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(movie_inst, nil)

	contentUseCase.
		EXPECT().
		GetFullByID(gomock.Eq(movie_inst.ContentID)).
		Return(content_inst, nil)

	dbMovie, err := movieUseCase.GetFullByID(movie_inst.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovie, movie_inst)
}

func TestMovieUseCase_GetFullByID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByID(gomock.Eq(movie_inst.ID)).
		Return(nil, sql.ErrNoRows)

	dbMovie, err := movieUseCase.GetFullByID(movie_inst.ID)
	assert.Equal(t, err, errors.Get(consts.CodeMovieDoesNotExist))
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
}

func TestMovieUseCase_GetByContentID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase)

	movieRep.
		EXPECT().
		SelectByContentID(gomock.Eq(movie_inst.ContentID)).
		Return(movie_inst, nil)

	dbMovie, err := movieUseCase.GetByContentID(movie_inst.ContentID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbMovie, movie_inst)
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
		SelectByContentID(gomock.Eq(movie_inst.ContentID)).
		Return(nil, sql.ErrNoRows)

	dbMovie, err := movieUseCase.GetByContentID(movie_inst.ContentID)
	assert.Equal(t, err, errors.Get(consts.CodeMovieDoesNotExist))
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
}
