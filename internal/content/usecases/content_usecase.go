package usecases

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type ContentUsecase struct {
	contentRepo   content.ContentRepository
	countryUcase  country.CountryUsecase
	genreUcase    genre.GenreUsecase
	actorUcase    actor.ActorUseCase
	directorUcase director.DirectorUseCase
}

func NewContentUsecase(repo content.ContentRepository, countryUcase country.CountryUsecase,
	genreUcase genre.GenreUsecase, actorUcase actor.ActorUseCase,
	directorUcase director.DirectorUseCase) content.ContentUsecase {
	return &ContentUsecase{
		contentRepo:   repo,
		countryUcase:  countryUcase,
		genreUcase:    genreUcase,
		actorUcase:    actorUcase,
		directorUcase: directorUcase,
	}
}

func (cu *ContentUsecase) Create(content *models.Content) *errors.Error {
	if err := cu.contentRepo.Insert(content); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (cu *ContentUsecase) UpdateByID(contentID uint64, newContentData *models.Content) (*models.Content, *errors.Error) {
	content, err := cu.GetFullByID(contentID)
	if err != nil {
		return nil, err
	}
	content.ReplaceBy(newContentData)

	if err := cu.contentRepo.Update(content); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return content, nil
}

func (cu *ContentUsecase) UpdatePosters(content *models.Content, newPostersDir string) *errors.Error {
	prevPostersDir := content.Images
	if newPostersDir == prevPostersDir {
		// Don't need to update
		return nil
	}

	// Update images
	content.Images = newPostersDir
	if err := cu.contentRepo.UpdateImages(content); err != nil {
		return errors.New(CodeInternalError, err)
	}
	// Don't need to delete prev directory,
	// cause posters always store into dir with the same name
	return nil
}

func (cu *ContentUsecase) DeleteByID(contentID uint64) *errors.Error {
	content, err := cu.GetByID(contentID)
	if err != nil {
		return errors.Get(CodeContentDoesNotExist)
	}

	// Delete posters dir
	if content.Images != "" {
		path, err := os.Getwd()
		if err != nil {
			return errors.New(CodeInternalError, err)
		}
		postersDirPath := filepath.Join(path, content.Images)

		if err := os.RemoveAll(postersDirPath); err != nil {
			return errors.New(CodeInternalError, err)
		}
	}

	if err := cu.contentRepo.DeleteByID(contentID); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (cu *ContentUsecase) GetByID(contentID uint64) (*models.Content, *errors.Error) {
	content, err := cu.contentRepo.SelectByID(contentID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeContentDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return content, nil
}

func (cu *ContentUsecase) GetFullByID(contentID uint64) (*models.Content, *errors.Error) {
	content, err := cu.GetByID(contentID)
	if err != nil {
		return nil, err
	}
	if err := cu.FillContent(content); err != nil {
		return nil, err
	}
	return content, nil
}

func (cu *ContentUsecase) FillContent(content *models.Content) *errors.Error {
	var err *errors.Error
	if content.Countries, err = cu.GetCountriesByID(content.ContentID); err != nil {
		return err
	}
	if content.Genres, err = cu.GetGenresByID(content.ContentID); err != nil {
		return err
	}
	if content.Actors, err = cu.GetActorsByID(content.ContentID); err != nil {
		return err
	}
	if content.Directors, err = cu.GetDirectorsByID(content.ContentID); err != nil {
		return err
	}
	return nil
}

func (cu *ContentUsecase) GetCountriesByID(contentID uint64) ([]*models.Country, *errors.Error) {
	countriesID, err := cu.contentRepo.SelectCountriesByID(contentID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	countries, customErr := cu.countryUcase.ListByID(countriesID)
	if customErr != nil {
		return nil, customErr
	}
	return countries, nil
}

func (cu *ContentUsecase) GetGenresByID(contentID uint64) ([]*models.Genre, *errors.Error) {
	genresID, err := cu.contentRepo.SelectGenresByID(contentID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	genres, customErr := cu.genreUcase.ListByID(genresID)
	if customErr != nil {
		return nil, customErr
	}
	return genres, nil
}

func (cu *ContentUsecase) GetActorsByID(contentID uint64) ([]*models.Actor, *errors.Error) {
	actorsID, err := cu.contentRepo.SelectActorsByID(contentID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	actors, customErr := cu.actorUcase.ListByID(actorsID)
	if customErr != nil {
		return nil, customErr
	}
	return actors, nil
}

func (cu *ContentUsecase) GetDirectorsByID(contentID uint64) ([]*models.Director, *errors.Error) {
	directorsID, err := cu.contentRepo.SelectDirectorsByID(contentID)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	directors, customErr := cu.directorUcase.ListByID(directorsID)
	if customErr != nil {
		return nil, customErr
	}
	return directors, nil
}
