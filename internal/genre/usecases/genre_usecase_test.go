package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/admin/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGenreUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		Name: "comedy",
	}

	grpcGenre := admin.GenreModelToGRPC(genre)
	adminPanelClient.
		EXPECT().
		CreateGenre(context.Background(), grpcGenre).
		Return(grpcGenre, nil)

	err := genreUseCase.Create(genre)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestGenreUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		Name: "comedy",
	}

	grpcGenre := admin.GenreModelToGRPC(genre)
	adminPanelClient.
		EXPECT().
		CreateGenre(context.Background(), grpcGenre).
		Return(nil, status.Error(codes.Code(consts.CodeGenreNameAlreadyExists), ""))

	err := genreUseCase.Create(genre)
	assert.Equal(t, err, errors.Get(consts.CodeGenreNameAlreadyExists))
}

func TestGenreUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	newGenreData := &models.Genre{
		ID:   1,
		Name: "Drama",
	}

	grpcGenre := admin.GenreModelToGRPC(newGenreData)
	adminPanelClient.
		EXPECT().
		ChangeGenre(context.Background(), grpcGenre).
		Return(&empty.Empty{}, nil)

	err := genreUseCase.Update(newGenreData)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestGenreUseCase_Update_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	newGenreData := &models.Genre{
		ID:   1,
		Name: "GB",
	}

	grpcGenre := admin.GenreModelToGRPC(newGenreData)
	adminPanelClient.
		EXPECT().
		ChangeGenre(context.Background(), grpcGenre).
		Return(&empty.Empty{}, status.Error(codes.Code(consts.CodeGenreNameAlreadyExists), ""))

	err := genreUseCase.Update(newGenreData)
	assert.Equal(t, err, errors.Get(consts.CodeGenreNameAlreadyExists))
}

func TestGenreUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	adminPanelClient.
		EXPECT().
		DeleteGenreByID(context.Background(), &admin.ID{ID: genre.ID}).
		Return(&empty.Empty{}, nil)

	err := genreUseCase.DeleteByID(genre.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestGenreUseCase_Delete_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	adminPanelClient.
		EXPECT().
		DeleteGenreByID(context.Background(), &admin.ID{ID: genre.ID}).
		Return(nil, status.Error(codes.Code(consts.CodeGenreDoesNotExist), ""))

	err := genreUseCase.DeleteByID(genre.ID)
	assert.Equal(t, err, errors.Get(consts.CodeGenreDoesNotExist))
}

func TestGenreUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(genre, nil)

	dbGenre, err := genreUseCase.GetByID(genre.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenre, genre)
}

func TestGenreUseCase_GetByID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(nil, sql.ErrNoRows)

	dbGenre, err := genreUseCase.GetByID(genre.ID)
	assert.Equal(t, err, errors.Get(consts.CodeGenreDoesNotExist))
	assert.Equal(t, dbGenre, (*models.Genre)(nil))
}

func TestGenreUseCase_GetByName_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(genre.Name)).
		Return(genre, nil)

	dbGenre, err := genreUseCase.GetByName(genre.Name)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenre, genre)
}

func TestGenreUseCase_GetByName_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(genre.Name)).
		Return(nil, sql.ErrNoRows)

	dbGenre, err := genreUseCase.GetByName(genre.Name)
	assert.Equal(t, err, errors.Get(consts.CodeGenreDoesNotExist))
	assert.Equal(t, dbGenre, (*models.Genre)(nil))
}

func TestGenreUseCase_List_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genres := []*models.Genre{
		&models.Genre{
			ID:   1,
			Name: "comedy",
		},
		&models.Genre{
			ID:   2,
			Name: "mult",
		},
	}

	genreRep.
		EXPECT().
		SelectAll().
		Return(genres, nil)

	dbGenres, err := genreUseCase.List()
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenres, genres)
}

func TestGenreUseCase_ListByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	genreUseCase := NewGenreUsecase(genreRep, adminPanelClient)

	genres := []*models.Genre{
		&models.Genre{
			ID:   1,
			Name: "comedy",
		},
		&models.Genre{
			ID:   2,
			Name: "mult",
		},
	}

	genresID := []uint64{1, 2}

	genreRep.
		EXPECT().
		SelectByID(genresID[0]).
		Return(genres[0], nil)

	genreRep.
		EXPECT().
		SelectByID(genresID[1]).
		Return(genres[1], nil)

	dbGenres, err := genreUseCase.ListByID(genresID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbGenres, genres)
}
