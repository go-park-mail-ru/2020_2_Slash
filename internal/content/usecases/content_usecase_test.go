package usecases

import (
	"context"
	actorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/admin/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	countryMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
	directorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	genreMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"testing"
)

var countries = []*models.Country{
	&models.Country{
		ID:   1,
		Name: "США",
	},
}

var genres = []*models.Genre{
	&models.Genre{
		Name: "Мультфильм",
	},
	&models.Genre{
		Name: "Комедия",
	},
}

var actors = []*models.Actor{
	&models.Actor{
		Name: "Майк Майерс",
	},
	&models.Actor{
		Name: "Эдди Мёрфи",
	},
}

var directors = []*models.Director{
	&models.Director{
		Name: "Эндрю Адамсон",
	},
	&models.Director{
		Name: "Вики Дженсон",
	},
}

var contentInst *models.Content = &models.Content{
	Name:             "Шрек",
	OriginalName:     "Shrek",
	Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
	ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
	Year:             2001,
	Countries:        countries,
	Genres:           genres,
	Actors:           actors,
	Directors:        directors,
	Type:             "movie",
}

func TestContentUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentRep := mocks.NewMockContentRepository(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)

	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase, adminPanelClient)

	contentRep.
		EXPECT().
		Insert(gomock.Eq(contentInst)).
		Return(nil)

	err := contentUseCase.Create(contentInst)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestContentUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentRep := mocks.NewMockContentRepository(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase, adminPanelClient)

	adminPanelClient.
		EXPECT().
		ChangeContent(context.Background(), admin.ContentModelToGRPC(contentInst)).
		Return(&empty.Empty{}, nil)

	err := contentUseCase.Update(contentInst)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestContentUseCase_UpdatePosters_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentRep := mocks.NewMockContentRepository(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)

	newPostersDir := "/images/0"

	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase, adminPanelClient)

	grpcContent := admin.ContentModelToGRPC(contentInst)
	adminPanelClient.
		EXPECT().
		ChangePosters(context.Background(), &admin.ContentPostersDir{
			Content:    grpcContent,
			PostersDir: newPostersDir,
		}).
		Return(grpcContent, nil)

	err := contentUseCase.UpdatePosters(contentInst, newPostersDir)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestContentUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentRep := mocks.NewMockContentRepository(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)

	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase, adminPanelClient)

	adminPanelClient.
		EXPECT().
		DeleteContentByID(context.Background(), &admin.ID{ID: contentInst.ContentID}).
		Return(&empty.Empty{}, nil)

	err := contentUseCase.DeleteByID(contentInst.ContentID)
	assert.Equal(t, err, (*errors.Error)(nil))
}
