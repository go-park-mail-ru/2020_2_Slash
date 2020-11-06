package usecases

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type CountryUsecase struct {
	countryRepo country.CountryRepository
}

func NewCountryUsecase(repo country.CountryRepository) country.CountryUsecase {
	return &CountryUsecase{
		countryRepo: repo,
	}
}

func (cu *CountryUsecase) Create(country *models.Country) *errors.Error {
	if err := cu.checkByName(country.Name); err == nil {
		return errors.Get(CodeCountryNameAlreadyExists)
	}

	if err := cu.countryRepo.Insert(country); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (cu *CountryUsecase) UpdateByID(countryID uint64, newCountryData *models.Country) (*models.Country, *errors.Error) {
	country, err := cu.GetByID(countryID)
	if err != nil {
		return nil, err
	}

	if newCountryData.Name == "" || country.Name == newCountryData.Name {
		return nil, err
	}

	if err := cu.checkByName(newCountryData.Name); err == nil {
		return nil, errors.Get(CodeCountryNameAlreadyExists)
	}

	country.Name = newCountryData.Name
	if err := cu.countryRepo.Update(country); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return country, nil
}

func (cu *CountryUsecase) DeleteByID(countryID uint64) *errors.Error {
	if err := cu.checkByID(countryID); err != nil {
		return errors.Get(CodeCountryDoesNotExist)
	}

	if err := cu.countryRepo.DeleteByID(countryID); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (cu *CountryUsecase) GetByID(countryID uint64) (*models.Country, *errors.Error) {
	country, err := cu.countryRepo.SelectByID(countryID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeCountryDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return country, nil
}

func (cu *CountryUsecase) GetByName(name string) (*models.Country, *errors.Error) {
	country, err := cu.countryRepo.SelectByName(name)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeCountryDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return country, nil
}

func (cu *CountryUsecase) List() ([]*models.Country, *errors.Error) {
	countries, err := cu.countryRepo.SelectAll()
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return countries, nil
}

func (cu *CountryUsecase) ListByID(countriesID []uint64) ([]*models.Country, *errors.Error) {
	var countries []*models.Country
	for _, countryID := range countriesID {
		country, err := cu.GetByID(countryID)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func (cu *CountryUsecase) checkByID(countryID uint64) *errors.Error {
	_, err := cu.GetByID(countryID)
	return err
}

func (cu *CountryUsecase) checkByName(name string) *errors.Error {
	_, err := cu.GetByName(name)
	return err
}
