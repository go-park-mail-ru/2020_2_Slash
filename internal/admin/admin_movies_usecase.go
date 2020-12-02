package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"path/filepath"
)

func (am *AdminMicroservice) CreateMovie(ctx context.Context, movie *Movie) (*Movie, error) {
	if err := am.checkByContentID(movie.GetContent().GetID()); err == nil {
		return &Movie{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	modelContent := ContentGRPCToModel(movie.Content)
	if err := am.CreateContent(modelContent); err != nil {
		return &Movie{}, err
	}
	movie.Content.ID = modelContent.ContentID

	modelMovie := MovieGRPCToModel(movie)
	if err := am.moviesRep.Insert(modelMovie); err != nil {
		return &Movie{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	movie.ID = modelMovie.ID

	return movie, nil
}

func (am *AdminMicroservice) ChangeVideo(ctx context.Context, videoMovie *VideoMovie) (*empty.Empty, error) {
	prevVideoPath := videoMovie.Movie.GetVideo()
	newVideoPath := videoMovie.GetVideo()
	if newVideoPath == prevVideoPath {
		// Don't need to update
		return &empty.Empty{}, nil
	}

	// Update video
	videoMovie.Movie.Video = newVideoPath
	if err := am.moviesRep.Update(MovieGRPCToModel(videoMovie.Movie)); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	// Don't need to delete prev file,
	// cause video always store with the same filename
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteMovieByID(ctx context.Context, movieID *ID) (*empty.Empty, error) {
	movie, err := am.GetMovieByID(movieID.GetID())
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeMovieDoesNotExist), "")
	}

	// Delete video
	if movie.Video != "" {
		path, err := os.Getwd()
		if err != nil {
			return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
		}
		videoPath := filepath.Join(path, movie.Video)
		videoDirPath := filepath.Dir(videoPath)

		if err := os.RemoveAll(videoDirPath); err != nil {
			return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
		}
	}

	if err := am.moviesRep.DeleteByID(movieID.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) GetMovieByID(movieID uint64) (*models.Movie, error) {
	movie, err := am.moviesRep.SelectByID(movieID)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeMovieDoesNotExist), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return movie, nil
}
