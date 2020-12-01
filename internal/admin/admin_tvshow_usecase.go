package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (am *AdminMicroservice) CreateTVShow(ctx context.Context, tvshow *TVShow) (*TVShow, error) {
	if err := am.checkByContentID(tvshow.GetContent().GetID()); err == nil {
		return &TVShow{}, status.Error(codes.Code(consts.CodeTVShowContentAlreadyExists), "")
	}

	modelContent := ContentGRPCToModel(tvshow.GetContent())
	if err := am.CreateContent(modelContent); err != nil {
		return &TVShow{}, err
	}
	tvshow.Content.ID = modelContent.ContentID

	tvshowModel := TVShowGRPCToModel(tvshow)
	if err := am.tvshowsRep.Insert(tvshowModel); err != nil {
		return &TVShow{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	tvshow.ID = tvshowModel.ID

	return tvshow, nil
}

func (am *AdminMicroservice) checkByContentID(contentID uint64) error {
	_, err := am.GetByContentID(contentID)
	return err
}

func (am *AdminMicroservice) GetByContentID(contentID uint64) (*models.TVShow, error) {
	tvshow, err := am.tvshowsRep.SelectByContentID(contentID)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeTVShowDoesNotExist), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return tvshow, nil
}
