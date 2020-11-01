package country

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type CountryRepository interface {
	Insert(country *models.Country) error
	Update(country *models.Country) error
	DeleteByID(countryID uint64) error
	SelectByID(countryID uint64) (*models.Country, error)
	SelectByName(name string) (*models.Country, error)
	SelectAll() ([]*models.Country, error)
}
