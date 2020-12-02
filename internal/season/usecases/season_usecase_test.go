package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/admin/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season/mocks"
	tvShowMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/mocks"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

var testSeason = &models.Season{
	ID:             1,
	Number:         3,
	EpisodesNumber: 8,
	TVShowID:       1,
	Episodes:       nil,
}

var testTvShow = &models.TVShow{
	ID:      1,
	Seasons: 3,
	Content: models.Content{},
}

var updTestSeason = &models.Season{
	ID:             1,
	Number:         4,
	EpisodesNumber: 8,
	TVShowID:       1,
	Episodes:       nil,
}

var existedSeason = &models.Season{
	ID:             2,
	Number:         4,
	EpisodesNumber: 9,
	TVShowID:       1,
	Episodes:       nil,
}

var testEpisodes = []*models.Episode{
	&models.Episode{
		ID:          1,
		Name:        "Рикбег из Рикшенка",
		Number:      1,
		Video:       "/videos/rickandmorty_22/3/1",
		Description: "Саммер решает спасти Рика из тюрьмы.",
		Poster:      "/images/rickandmorty_22/3/1",
		SeasonID:    3,
	},
	&models.Episode{
		ID:          2,
		Name:        "Рикман с камнем",
		Number:      2,
		Video:       "/videos/rickandmorty_22/3/2",
		Description: "Рик, Морти и Саммер охотятся за новым источником энергии в постакалиптической версии Земли.",
		Poster:      "/images/rickandmorty_22/3/2",
		SeasonID:    3,
	},
	&models.Episode{
		ID:          2,
		Name:        "Огурчик Рик",
		Number:      3,
		Video:       "/videos/rickandmorty_22/3/3",
		Description: "Рик превращает себя в огурчик.",
		Poster:      "/images/rickandmorty_22/3/3",
		SeasonID:    3,
	},
}

func TestSeasonUsecase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	grpcSeason := admin.SeasonModelToGRPC(testSeason)
	adminPanelClient.
		EXPECT().
		CreateSeason(context.Background(), grpcSeason).
		Return(grpcSeason, nil)

	customErr := seasonUsecase.Create(testSeason)
	assert.Nil(t, customErr)
}

func TestSeasonUsecase_Create_Conflict(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	grpcSeason := admin.SeasonModelToGRPC(testSeason)
	adminPanelClient.
		EXPECT().
		CreateSeason(context.Background(), grpcSeason).
		Return(&admin.Season{}, status.Error(codes.Code(consts.CodeSeasonAlreadyExist), ""))

	customErr := seasonUsecase.Create(testSeason)
	assert.Equal(t, errors.Get(consts.CodeSeasonAlreadyExist), customErr)
}

func TestSeasonUsecase_Change_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	grpcSeason := admin.SeasonModelToGRPC(updTestSeason)
	adminPanelClient.
		EXPECT().
		ChangeSeason(context.Background(), grpcSeason).
		Return(&empty.Empty{}, nil)
	customErr := seasonUsecase.Change(updTestSeason)
	assert.Nil(t, customErr)
}

func TestSeasonUsecase_Change_NoSeason(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	grpcSeason := admin.SeasonModelToGRPC(updTestSeason)
	adminPanelClient.
		EXPECT().
		ChangeSeason(context.Background(), grpcSeason).
		Return(&empty.Empty{}, status.Error(codes.Code(consts.CodeSeasonDoesNotExist), ""))

	customErr := seasonUsecase.Change(updTestSeason)
	assert.Equal(t, errors.Get(consts.CodeSeasonDoesNotExist), customErr)
}

func TestSeasonUsecase_Change_Equal(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	grpcSeason := admin.SeasonModelToGRPC(updTestSeason)
	adminPanelClient.
		EXPECT().
		ChangeSeason(context.Background(), grpcSeason).
		Return(&empty.Empty{}, nil)

	customErr := seasonUsecase.Change(updTestSeason)
	assert.Nil(t, customErr)
}

func TestSeasonUsecase_Change_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	grpcSeason := admin.SeasonModelToGRPC(updTestSeason)
	adminPanelClient.
		EXPECT().
		ChangeSeason(context.Background(), grpcSeason).
		Return(&empty.Empty{}, status.Error(codes.Code(consts.CodeSeasonAlreadyExist), ""))

	customErr := seasonUsecase.Change(updTestSeason)
	assert.Equal(t, errors.Get(consts.CodeSeasonAlreadyExist), customErr)
}

func TestSeasonUsecase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	grpcSeason := admin.SeasonModelToGRPC(updTestSeason)
	adminPanelClient.
		EXPECT().
		DeleteSeasonsByID(context.Background(), &admin.ID{ID: grpcSeason.ID}).
		Return(&empty.Empty{}, nil)

	customErr := seasonUsecase.Delete(testSeason.ID)
	assert.Nil(t, customErr)
}

func TestSeasonUsecase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectByID(testSeason.ID).
		Return(testSeason, nil)

	seasonDB, customErr := seasonUsecase.Get(testSeason.ID)
	assert.Equal(t, testSeason, seasonDB)
	assert.Nil(t, customErr)
}

func TestSeasonUsecase_Get_NoRows(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectByID(testSeason.ID).
		Return(nil, sql.ErrNoRows)

	seasonDB, customErr := seasonUsecase.Get(testSeason.ID)
	assert.Equal(t, errors.Get(consts.CodeSeasonDoesNotExist), customErr)
	assert.Nil(t, seasonDB)
}

func TestSeasonUsecase_GetEpisodes_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectEpisodes(testSeason.ID).
		Return(testEpisodes, nil)

	episodes, customErr := seasonUsecase.GetEpisodes(testSeason.ID)
	assert.Equal(t, testEpisodes, episodes)
	assert.Nil(t, customErr)
}

func TestSeasonUsecase_GetEpisodes_Nil(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectEpisodes(testSeason.ID).
		Return(nil, nil)

	episodes, customErr := seasonUsecase.GetEpisodes(testSeason.ID)
	assert.Equal(t, []*models.Episode{}, episodes)
	assert.Nil(t, customErr)
}

func TestSeasonUsecase_GetEpisodes_NoRows(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockSeasonRepository(ctrl)
	tvShowUsecase := tvShowMocks.NewMockTVShowUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewSeasonUsecase(rep, tvShowUsecase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectEpisodes(testSeason.ID).
		Return(nil, sql.ErrNoRows)

	episodes, customErr := seasonUsecase.GetEpisodes(testSeason.ID)
	assert.Equal(t, []*models.Episode{}, episodes)
	assert.Nil(t, customErr)
}
