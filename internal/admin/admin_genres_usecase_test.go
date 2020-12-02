package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenreUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		genresRep: genreRep,
	}

	genre := &models.Genre{
		Name: "USA",
	}

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(genre.Name)).
		Return(nil, sql.ErrNoRows)

	genreRep.
		EXPECT().
		Insert(gomock.Eq(genre)).
		Return(nil)

	_, err := adminMicroservice.CreateGenre(context.Background(), GenreModelToGRPC(genre))
	assert.Equal(t, err, (error)(nil))
}

func TestGenreUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreRep := mocks.NewMockGenreRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		genresRep: genreRep,
	}

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	newGenreData := &models.Genre{
		ID:   1,
		Name: "GB",
	}

	genreRep.
		EXPECT().
		SelectByID(gomock.Eq(genre.ID)).
		Return(genre, nil)

	genreRep.
		EXPECT().
		SelectByName(gomock.Eq(newGenreData.Name)).
		Return(nil, sql.ErrNoRows)

	genreRep.
		EXPECT().
		Update(gomock.Eq(genre)).
		Return(nil)

	_, err := adminMicroservice.ChangeGenre(context.Background(), GenreModelToGRPC(newGenreData))
	assert.Equal(t, err, (error)(nil))
}

