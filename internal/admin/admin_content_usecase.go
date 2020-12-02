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

func (am *AdminMicroservice) ChangePosters(ctx context.Context, contentPostersDir *ContentPostersDir) (*Content, error) {
	content := contentPostersDir.GetContent()
	prevPostersDir := content.GetImages()
	if contentPostersDir.PostersDir == prevPostersDir {
		// Don't need to update
		return contentPostersDir.GetContent(), nil
	}

	// Update images
	content.Images = contentPostersDir.GetPostersDir()
	if err := am.contentRep.UpdateImages(ContentGRPCToModel(content)); err != nil {
		return &Content{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	// Don't need to delete prev directory,
	// cause posters always store into dir with the same name
	return &Content{}, nil
}

func (am *AdminMicroservice) ChangeContent(ctx context.Context, newContentData *Content) (*empty.Empty, error) {
	modelContent, err := am.GetFullByID(newContentData.GetID())
	if err != nil {
		return &empty.Empty{}, err
	}
	modelContent.ReplaceBy(ContentGRPCToModel(newContentData))

	if err := am.contentRep.Update(modelContent); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteContentByID(ctx context.Context, contentID *ID) (*empty.Empty, error) {
	content, err := am.GetContentByID(contentID.GetID())
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeContentDoesNotExist), err.Error())
	}

	// Delete posters dir
	if content.Images != "" {
		path, err := os.Getwd()
		if err != nil {
			return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
		}
		postersDirPath := filepath.Join(path, content.Images)

		if err := os.RemoveAll(postersDirPath); err != nil {
			return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
		}
	}

	if err := am.contentRep.DeleteByID(contentID.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) CreateContent(content *models.Content) error {
	if err := am.contentRep.Insert(content); err != nil {
		return status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return nil
}

func (am *AdminMicroservice) GetFullByID(contentID uint64) (*models.Content, error) {
	content, err := am.GetContentByID(contentID)
	if err != nil {
		return nil, err
	}
	if err := am.FillContent(content); err != nil {
		return nil, err
	}
	return content, nil
}

func (am *AdminMicroservice) GetContentByID(contentID uint64) (*models.Content, error) {
	content, err := am.contentRep.SelectByID(contentID)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeContentDoesNotExist), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return content, nil
}

func (am *AdminMicroservice) FillContent(content *models.Content) error {
	var err error
	if content.Countries, err = am.GetCountriesByID(content.ContentID); err != nil {
		return err
	}
	if content.Genres, err = am.GetGenresByID(content.ContentID); err != nil {
		return err
	}
	if content.Actors, err = am.GetActorsByID(content.ContentID); err != nil {
		return err
	}
	if content.Directors, err = am.GetDirectorsByID(content.ContentID); err != nil {
		return err
	}
	return nil
}

func (am *AdminMicroservice) GetCountriesByID(contentID uint64) ([]*models.Country, error) {
	countriesID, err := am.contentRep.SelectCountriesByID(contentID)
	if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	countries, customErr := am.ListCountriesByID(countriesID)
	if customErr != nil {
		return nil, customErr
	}
	return countries, nil
}

func (am *AdminMicroservice) GetGenresByID(contentID uint64) ([]*models.Genre, error) {
	genresID, err := am.contentRep.SelectGenresByID(contentID)
	if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	genres, customErr := am.ListGenresByID(genresID)
	if customErr != nil {
		return nil, customErr
	}
	return genres, nil
}

func (am *AdminMicroservice) GetActorsByID(contentID uint64) ([]*models.Actor, error) {
	actorsID, err := am.contentRep.SelectActorsByID(contentID)
	if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	actors, customErr := am.ListActorsByID(actorsID)
	if customErr != nil {
		return nil, customErr
	}
	return actors, nil
}

func (am *AdminMicroservice) GetDirectorsByID(contentID uint64) ([]*models.Director, error) {
	directorsID, err := am.contentRep.SelectDirectorsByID(contentID)
	if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	directors, customErr := am.ListDirectorsByID(directorsID)
	if customErr != nil {
		return nil, customErr
	}
	return directors, nil
}
