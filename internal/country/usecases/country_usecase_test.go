package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountryUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByName(gomock.Eq(country.Name)).
		Return(nil, sql.ErrNoRows)

	countryRep.
		EXPECT().
		Insert(gomock.Eq(country)).
		Return(nil)

	err := countryUseCase.Create(country)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestCountryUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByName(gomock.Eq(country.Name)).
		Return(country, nil)

	err := countryUseCase.Create(country)
	assert.Equal(t, err, errors.Get(consts.CodeCountryNameAlreadyExists))
}

func TestCountryUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	newCountryData := &models.Country{
		ID:   1,
		Name: "GB",
	}

	countryRep.
		EXPECT().
		SelectByID(gomock.Eq(country.ID)).
		Return(country, nil)

	countryRep.
		EXPECT().
		SelectByName(gomock.Eq(newCountryData.Name)).
		Return(nil, sql.ErrNoRows)

	countryRep.
		EXPECT().
		Update(gomock.Eq(country)).
		Return(nil)

	dbCountry, err := countryUseCase.UpdateByID(country.ID, newCountryData)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbCountry, newCountryData)
}

func TestCountryUseCase_Update_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	newCountryData := &models.Country{
		ID:   1,
		Name: "GB",
	}

	countryRep.
		EXPECT().
		SelectByID(gomock.Eq(country.ID)).
		Return(country, nil)

	countryRep.
		EXPECT().
		SelectByName(gomock.Eq(newCountryData.Name)).
		Return(newCountryData, nil)

	dbCountry, err := countryUseCase.UpdateByID(country.ID, newCountryData)
	assert.Equal(t, err, errors.Get(consts.CodeCountryNameAlreadyExists))
	assert.Equal(t, dbCountry, (*models.Country)(nil))
}

func TestCountryUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByID(gomock.Eq(country.ID)).
		Return(country, nil)

	countryRep.
		EXPECT().
		DeleteByID(gomock.Eq(country.ID)).
		Return(nil)

	err := countryUseCase.DeleteByID(country.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestCountryUseCase_Delete_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByID(gomock.Eq(country.ID)).
		Return(nil, sql.ErrNoRows)

	err := countryUseCase.DeleteByID(country.ID)
	assert.Equal(t, err, errors.Get(consts.CodeCountryDoesNotExist))
}

func TestCountryUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByID(gomock.Eq(country.ID)).
		Return(country, nil)

	dbCountry, err := countryUseCase.GetByID(country.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbCountry, country)
}

func TestCountryUseCase_GetByID_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByID(gomock.Eq(country.ID)).
		Return(nil, sql.ErrNoRows)

	dbCountry, err := countryUseCase.GetByID(country.ID)
	assert.Equal(t, err, errors.Get(consts.CodeCountryDoesNotExist))
	assert.Equal(t, dbCountry, (*models.Country)(nil))
}

func TestCountryUseCase_GetByName_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByName(gomock.Eq(country.Name)).
		Return(country, nil)

	dbCountry, err := countryUseCase.GetByName(country.Name)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbCountry, country)
}

func TestCountryUseCase_GetByName_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryRep.
		EXPECT().
		SelectByName(gomock.Eq(country.Name)).
		Return(nil, sql.ErrNoRows)

	dbCountry, err := countryUseCase.GetByName(country.Name)
	assert.Equal(t, err, errors.Get(consts.CodeCountryDoesNotExist))
	assert.Equal(t, dbCountry, (*models.Country)(nil))
}

func TestCountryUseCase_List_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	countries := []*models.Country{
		&models.Country{
			ID:   1,
			Name: "USA",
		},
		&models.Country{
			ID:   2,
			Name: "GB",
		},
	}

	countryRep.
		EXPECT().
		SelectAll().
		Return(countries, nil)

	dbCountries, err := countryUseCase.List()
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbCountries, countries)
}

func TestCountryUseCase_ListByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	countryUseCase := NewCountryUsecase(countryRep)

	countries := []*models.Country{
		&models.Country{
			ID:   1,
			Name: "USA",
		},
		&models.Country{
			ID:   2,
			Name: "GB",
		},
	}

	countriesID := []uint64{1, 2}

	countryRep.
		EXPECT().
		SelectByID(countriesID[0]).
		Return(countries[0], nil)

	countryRep.
		EXPECT().
		SelectByID(countriesID[1]).
		Return(countries[1], nil)

	dbCountries, err := countryUseCase.ListByID(countriesID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbCountries, countries)
}
