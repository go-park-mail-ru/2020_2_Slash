package usecases

import (
	"context"
	"database/sql"
	"github.com/golang/protobuf/ptypes/empty"

	//"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/admin/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	seasonMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/season/mocks"
	//tvShowMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/mocks"
	"github.com/golang/mock/gomock"
	//"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

var testEpisode = &models.Episode{
	ID:          2,
	Name:        "Огурчик Рик",
	Number:      3,
	Video:       "/videos/rickandmorty_22/3/3",
	Description: "Рик превращает себя в огурчик.",
	Poster:      "/images/rickandmorty_22/3/3",
	SeasonID:    3,
}

var testSeason = &models.Season{
	ID:             3,
	Number:         3,
	EpisodesNumber: 3,
	TVShowID:       1,
	Episodes:       []*models.Episode{testEpisode},
}

var testContent = &models.Content{
	ContentID:        1,
	Name:             "Рик и Морти",
	OriginalName:     "Rick and morty",
	Description:      "",
	ShortDescription: "",
	Rating:           0,
	Year:             0,
	Images:           "",
	Type:             "",
	Countries:        nil,
	Genres:           nil,
	Actors:           nil,
	Directors:        nil,
	IsLiked:          nil,
	IsFavourite:      nil,
}

var updTestEpisode = &models.Episode{
	ID:          2,
	Name:        "Огурчик Рик UPD",
	Number:      3,
	Video:       "/videos/rickandmorty_22/3/3",
	Description: "Рик превращает себя в огурчик.",
	Poster:      "/images/rickandmorty_22/3/3",
	SeasonID:    3,
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

func TestEpisodeUsecase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	grpcEpisode := admin.EpisodeModelToGRPC(testEpisode)
	adminPanelClient.
		EXPECT().
		CreateEpisode(context.Background(), grpcEpisode).
		Return(grpcEpisode, nil)

	customErr := seasonUsecase.Create(testEpisode)
	assert.Nil(t, customErr)
}

func TestEpisodeUsecase_Create_Conflict(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	grpcEpisode := admin.EpisodeModelToGRPC(testEpisode)
	adminPanelClient.
		EXPECT().
		CreateEpisode(context.Background(), grpcEpisode).
		Return(&admin.Episode{}, status.Error(codes.Code(consts.CodeEpisodeAlreadyExist), ""))

	customErr := seasonUsecase.Create(testEpisode)
	assert.Equal(t, errors.Get(consts.CodeEpisodeAlreadyExist), customErr)
}

func TestEpisodeUsecase_Change_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	grpcEpisode := admin.EpisodeModelToGRPC(updTestEpisode)
	adminPanelClient.
		EXPECT().
		ChangeEpisode(context.Background(), grpcEpisode).
		Return(&empty.Empty{}, nil)
	customErr := seasonUsecase.Change(updTestEpisode)
	assert.Nil(t, customErr)
}

func TestEpisodeUsecase_Change_NoEpisode(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	grpcEpisode := admin.EpisodeModelToGRPC(updTestEpisode)
	adminPanelClient.
		EXPECT().
		ChangeEpisode(context.Background(), grpcEpisode).
		Return(&empty.Empty{}, status.Error(codes.Code(consts.CodeEpisodeDoesNotExist), ""))

	customErr := seasonUsecase.Change(updTestEpisode)
	assert.Equal(t, errors.Get(consts.CodeEpisodeDoesNotExist), customErr)
}

func TestEpisodeUsecase_Change_Equal(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	grpcEpisode := admin.EpisodeModelToGRPC(updTestEpisode)
	adminPanelClient.
		EXPECT().
		ChangeEpisode(context.Background(), grpcEpisode).
		Return(&empty.Empty{}, nil)

	customErr := seasonUsecase.Change(updTestEpisode)
	assert.Nil(t, customErr)
}

func TestEpisodeUsecase_Change_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	grpcEpisode := admin.EpisodeModelToGRPC(updTestEpisode)
	adminPanelClient.
		EXPECT().
		ChangeEpisode(context.Background(), grpcEpisode).
		Return(&empty.Empty{}, status.Error(codes.Code(consts.CodeEpisodeAlreadyExist), ""))

	customErr := seasonUsecase.Change(updTestEpisode)
	assert.Equal(t, errors.Get(consts.CodeEpisodeAlreadyExist), customErr)
}

func TestEpisodeUsecase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	grpcEpisode := admin.EpisodeModelToGRPC(updTestEpisode)
	adminPanelClient.
		EXPECT().
		DeleteEpisodeByID(context.Background(), &admin.ID{ID: grpcEpisode.ID}).
		Return(&empty.Empty{}, nil)

	customErr := seasonUsecase.DeleteByID(testEpisode.ID)
	assert.Nil(t, customErr)
}

func TestEpisodeUsecase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectByID(testEpisode.ID).
		Return(testEpisode, nil)

	seasonDB, customErr := seasonUsecase.GetByID(testEpisode.ID)
	assert.Equal(t, testEpisode, seasonDB)
	assert.Nil(t, customErr)
}

func TestEpisodeUsecase_Get_NoRows(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectByID(testEpisode.ID).
		Return(nil, sql.ErrNoRows)

	seasonDB, customErr := seasonUsecase.GetByID(testEpisode.ID)
	assert.Equal(t, errors.Get(consts.CodeEpisodeDoesNotExist), customErr)
	assert.Nil(t, seasonDB)
}

func TestEpisodeUsecase_GetContentByEID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectContentByID(testEpisode.ID).
		Return(testContent, nil)

	content, customErr := seasonUsecase.GetContentByEID(testEpisode.ID)
	assert.Equal(t, testContent, content)
	assert.Nil(t, customErr)
}

func TestEpisodeUsecase_GetSeasonNumber(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	rep.
		EXPECT().
		SelectSeasonNumberByID(testEpisode.ID).
		Return(testEpisode.Number, nil)

	seasonNumber, customErr := seasonUsecase.GetSeasonNumber(testEpisode.ID)
	assert.Equal(t, testEpisode.Number, seasonNumber)
	assert.Nil(t, customErr)
}

var testPoster = "/images/rickandmorty_22/3/3.png"

func TestEpisodeUsecase_UpdatePoster(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	testEpisode.Poster = testPoster
	grpcEpisode := admin.EpisodeModelToGRPC(testEpisode)
	adminPanelClient.
		EXPECT().
		UpdatePoster(context.Background(), &admin.EpisodePostersDir{
			Episode:    grpcEpisode,
			PostersDir: testPoster,
		}).
		Return(&empty.Empty{}, nil)

	customErr := seasonUsecase.UpdatePoster(testEpisode, testPoster)
	assert.Nil(t, customErr)
}

var testVideo = "video"

func TestEpisodeUsecase_UpdateVideo(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockEpisodeRepository(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	seasonUsecase := NewEpisodeUsecase(rep, seasonUseCase, adminPanelClient)
	defer ctrl.Finish()

	testEpisode.Video = testVideo
	grpcEpisode := admin.EpisodeModelToGRPC(testEpisode)
	adminPanelClient.
		EXPECT().
		UpdateVideo(context.Background(), &admin.EpisodeVideo{
			Episode: grpcEpisode,
			Video:   testVideo,
		}).
		Return(&empty.Empty{}, nil)

	customErr := seasonUsecase.UpdateVideo(testEpisode, testVideo)
	assert.Nil(t, customErr)
}