package usecases

import (
	actorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	countryMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
	directorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	genreMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
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

	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase)

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

	countriesID := []uint64{1}
	directorsID := []uint64{1, 2}
	actorsID := []uint64{1, 2}
	genresID := []uint64{1, 2}

	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase)

	contentRep.
		EXPECT().
		SelectByID(gomock.Eq(contentInst.ContentID)).
		Return(contentInst, nil)

	contentRep.
		EXPECT().
		SelectCountriesByID(gomock.Eq(contentInst.ContentID)).
		Return(countriesID, nil)

	countryUseCase.
		EXPECT().
		ListByID(gomock.Eq(countriesID)).
		Return(countries, nil)

	contentRep.
		EXPECT().
		SelectGenresByID(gomock.Eq(contentInst.ContentID)).
		Return(genresID, nil)

	genreUseCase.
		EXPECT().
		ListByID(gomock.Eq(genresID)).
		Return(genres, nil)

	contentRep.
		EXPECT().
		SelectActorsByID(gomock.Eq(contentInst.ContentID)).
		Return(actorsID, nil)

	actorUseCase.
		EXPECT().
		ListByID(gomock.Eq(actorsID)).
		Return(actors, nil)

	contentRep.
		EXPECT().
		SelectDirectorsByID(gomock.Eq(contentInst.ContentID)).
		Return(directorsID, nil)

	directorUseCase.
		EXPECT().
		ListByID(gomock.Eq(directorsID)).
		Return(directors, nil)

	contentRep.
		EXPECT().
		Update(gomock.Eq(contentInst)).
		Return(nil)

	dbContent, err := contentUseCase.UpdateByID(contentInst.ContentID, contentInst)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbContent, contentInst)
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

	newPostersDir := "/images/0"

	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase)

	contentRep.
		EXPECT().
		UpdateImages(gomock.Eq(contentInst)).
		Return(nil)

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

	contentUseCase := NewContentUsecase(contentRep, countryUseCase,
		genreUseCase, actorUseCase, directorUseCase)

	contentRep.
		EXPECT().
		SelectByID(gomock.Eq(contentInst.ContentID)).
		Return(contentInst, nil)

	contentRep.
		EXPECT().
		DeleteByID(gomock.Eq(contentInst.ContentID)).
		Return(nil)

	err := contentUseCase.DeleteByID(contentInst.ContentID)
	assert.Equal(t, err, (*errors.Error)(nil))
}
