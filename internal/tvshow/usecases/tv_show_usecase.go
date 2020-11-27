package usecases

import (
	"database/sql"

	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
)

type TVShowUsecase struct {
	tvshowRepo   tvshow.TVShowRepository
	contentUcase content.ContentUsecase
}

func NewTVShowUsecase(repo tvshow.TVShowRepository,
	contentUcase content.ContentUsecase) tvshow.TVShowUsecase {
	return &TVShowUsecase{
		tvshowRepo:   repo,
		contentUcase: contentUcase,
	}
}

func (tu *TVShowUsecase) Create(tvshow *models.TVShow) *errors.Error {
	if err := tu.checkByContentID(tvshow.ContentID); err == nil {
		return errors.Get(CodeTVShowContentAlreadyExists)
	}

	if err := tu.contentUcase.Create(&tvshow.Content); err != nil {
		return err
	}

	if err := tu.tvshowRepo.Insert(tvshow); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (tu *TVShowUsecase) GetByID(tvshowID uint64) (*models.TVShow, *errors.Error) {
	tvshow, err := tu.tvshowRepo.SelectByID(tvshowID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeTVShowDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return tvshow, nil
}

func (tu *TVShowUsecase) GetFullByID(tvshowID uint64, curUserID uint64) (*models.TVShow, *errors.Error) {
	tvshow, err := tu.tvshowRepo.SelectFullByID(tvshowID, curUserID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeTVShowDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	customErr := tu.contentUcase.FillContent(&tvshow.Content)
	if customErr != nil {
		return nil, customErr
	}
	return tvshow, nil
}

func (tu *TVShowUsecase) GetByContentID(contentID uint64) (*models.TVShow, *errors.Error) {
	tvshow, err := tu.tvshowRepo.SelectByContentID(contentID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeTVShowDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return tvshow, nil
}

func (tu *TVShowUsecase) ListByParams(params *models.ContentFilter, pgnt *models.Pagination,
	curUserID uint64) ([]*models.TVShow, *errors.Error) {

	tvshows, err := tu.tvshowRepo.SelectByParams(params, pgnt, curUserID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(tvshows) == 0 {
		return []*models.TVShow{}, nil
	}
	return tvshows, nil
}

func (tu *TVShowUsecase) ListLatest(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, *errors.Error) {
	tvshows, err := tu.tvshowRepo.SelectLatest(pgnt, curUserID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(tvshows) == 0 {
		return []*models.TVShow{}, nil
	}
	return tvshows, nil
}

func (tu *TVShowUsecase) ListByRating(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, *errors.Error) {
	tvshows, err := tu.tvshowRepo.SelectByRating(pgnt, curUserID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}

	if len(tvshows) == 0 {
		return []*models.TVShow{}, nil
	}

	return tvshows, nil
}

func (tu *TVShowUsecase) checkByContentID(contentID uint64) *errors.Error {
	_, err := tu.GetByContentID(contentID)
	return err
}
