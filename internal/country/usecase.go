package country

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type CountryUsecase interface {
	Create(country *models.Country) *errors.Error
	UpdateByID(countryID uint64, newCountryData *models.Country) (*models.Country, *errors.Error)
	DeleteByID(countryID uint64) *errors.Error
	GetByID(countryID uint64) (*models.Country, *errors.Error)
	GetByName(name string) (*models.Country, *errors.Error)
	List() ([]*models.Country, *errors.Error)
}
