package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
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
	adminMicroservice := &AdminMicroservice{
		countriesRep: countryRep,
	}

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

	_, err := adminMicroservice.CreateCountry(context.Background(), CountryModelToGRPC(country))
	assert.Equal(t, err, (error)(nil))
}

func TestCountryUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryRep := mocks.NewMockCountryRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		countriesRep: countryRep,
	}

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

	_, err := adminMicroservice.ChangeCountry(context.Background(), CountryModelToGRPC(newCountryData))
	assert.Equal(t, err, (error)(nil))
}
