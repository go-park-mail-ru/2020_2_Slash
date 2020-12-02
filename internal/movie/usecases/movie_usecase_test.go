package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/admin/mocks"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

	grpcMovie := admin.MovieModelToGRPC(movieInst)
	adminPanelClient.
		EXPECT().
		CreateMovie(context.Background(), grpcMovie).
		Return(grpcMovie, nil)

	err := movieUseCase.Create(movieInst)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestMovieUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

	grpcMovie := admin.MovieModelToGRPC(movieInst)
	adminPanelClient.
		EXPECT().
		CreateMovie(context.Background(), grpcMovie).
		Return(&admin.Movie{}, status.Error(codes.Code(consts.CodeMovieContentAlreadyExists), ""))

	err := movieUseCase.Create(movieInst)
	assert.Equal(t, err, errors.Get(consts.CodeMovieContentAlreadyExists))
}

func TestMovieUseCase_UpdateVideo_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

	newVideoPath := "video/movie.mp4"

	grpcMovie := admin.MovieModelToGRPC(movieInst)
	adminPanelClient.
		EXPECT().
		ChangeVideo(context.Background(), &admin.VideoMovie{
			Video: newVideoPath,
			Movie: grpcMovie,
		}).
		Return(&empty.Empty{}, nil)

	err := movieUseCase.UpdateVideo(movieInst, newVideoPath)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestMovieUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRep := mocks.NewMockMovieRepository(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)
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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

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
		Year:     2001,
		Genre:    1,
		Country:  1,
		Actor:    1,
		Director: 1,
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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	movieUseCase := NewMovieUsecase(movieRep, contentUseCase, adminPanelClient)

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
