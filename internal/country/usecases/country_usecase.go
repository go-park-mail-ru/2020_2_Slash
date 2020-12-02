package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	"github.com/jinzhu/copier"

	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type CountryUsecase struct {
	countryRepo      country.CountryRepository
	adminPanelClient admin.AdminPanelClient
}

func NewCountryUsecase(repo country.CountryRepository,
	client admin.AdminPanelClient) country.CountryUsecase {
	return &CountryUsecase{
		countryRepo:      repo,
		adminPanelClient: client,
	}
}

func (cu *CountryUsecase) Create(country *models.Country) *errors.Error {
	grpcCountry, err := cu.adminPanelClient.CreateCountry(context.Background(),
		admin.CountryModelToGRPC(country))
	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	if err := copier.Copy(country, admin.CountryGRPCToModel(grpcCountry)); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (cu *CountryUsecase) Update(newCountryData *models.Country) *errors.Error {
	_, err := cu.adminPanelClient.ChangeCountry(context.Background(),
		admin.CountryModelToGRPC(newCountryData))

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (cu *CountryUsecase) DeleteByID(countryID uint64) *errors.Error {
	_, err := cu.adminPanelClient.DeleteCountryByID(context.Background(),
		&admin.ID{ID: countryID})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
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
	if len(countries) == 0 {
		return []*models.Country{}, nil
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
