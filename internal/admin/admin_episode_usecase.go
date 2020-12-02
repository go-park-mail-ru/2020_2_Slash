package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (am *AdminMicroservice) CreateEpisode(ctx context.Context, episode *Episode) (*Episode, error) {
	_, err := am.GetSeason(episode.SeasonID)
	if err != nil {
		return &Episode{}, err
	}

	modelEpisode := EpisodeGRPCToModel(episode)
	isConflict, err := am.isEpisodeConflict(modelEpisode)
	if err != nil {
		return &Episode{}, err
	}
	if isConflict {
		return &Episode{}, status.Error(codes.Code(consts.CodeEpisodeAlreadyExist), "")
	}

	err = am.episodesRep.Insert(modelEpisode)
	if err != nil {
		return &Episode{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	episode.ID = modelEpisode.ID

	return episode, nil
}

func (am *AdminMicroservice) ChangeEpisode(ctx context.Context, episode *Episode) (*empty.Empty, error) {
	_, err := am.GetSeason(episode.SeasonID)
	if err != nil {
		return &empty.Empty{}, err
	}

	modelEpisode := EpisodeGRPCToModel(episode)
	isExist, err := am.isEpisodeExist(modelEpisode)
	if err != nil {
		return &empty.Empty{}, err
	}
	if !isExist {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeEpisodeDoesNotExist), err.Error())
	}

	err = am.episodesRep.Update(modelEpisode)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteEpisodeByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	_, err := am.GetEpisodeByID(id.GetID())
	if err != nil {
		return &empty.Empty{}, err
	}

	err = am.episodesRep.DeleteByID(id.GetID())
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) UpdatePoster(ctx context.Context, episodePosterDir *EpisodePostersDir) (*empty.Empty, error) {
	prevPosterPath := episodePosterDir.Episode.GetPoster()
	newPosterPath := episodePosterDir.GetPostersDir()
	if newPosterPath == prevPosterPath {
		return &empty.Empty{}, nil
	}

	episodePosterDir.Episode.Poster = newPosterPath
	if err := am.episodesRep.UpdatePoster(EpisodeGRPCToModel(episodePosterDir.Episode)); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) UpdateVideo(ctx context.Context, episodeVideo *EpisodeVideo) (*empty.Empty, error) {
	prevVideoPath := episodeVideo.Episode.GetVideo()
	newVideoPath := episodeVideo.GetVideo()
	if newVideoPath == prevVideoPath {
		return &empty.Empty{}, nil
	}

	episodeVideo.Episode.Video = newVideoPath
	if err := am.episodesRep.UpdateVideo(EpisodeGRPCToModel(episodeVideo.Episode)); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) GetEpisodeByID(id uint64) (*models.Episode, error) {
	dbEpisode, err := am.episodesRep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.Code(consts.CodeEpisodeDoesNotExist), "")
	} else if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return dbEpisode, nil
}

func (am *AdminMicroservice) isEpisodeConflict(episode *models.Episode) (bool, error) {
	_, err := am.episodesRep.SelectByNumberAndSeason(episode.Number, episode.SeasonID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return true, nil
}

func (am *AdminMicroservice) isEpisodeExist(episode *models.Episode) (bool, error) {
	_, err := am.episodesRep.SelectByID(episode.ID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return true, nil
}

func (am *AdminMicroservice) GetContentByEID(eid uint64) (*models.Content, *errors.Error) {
	content, err := am.episodesRep.SelectContentByID(eid)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return content, nil
}

func (am *AdminMicroservice) GetSeasonNumber(eid uint64) (int, error) {
	seasonNumber, err := am.episodesRep.SelectSeasonNumberByID(eid)
	if err != nil {
		return 0, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return seasonNumber, nil
}
