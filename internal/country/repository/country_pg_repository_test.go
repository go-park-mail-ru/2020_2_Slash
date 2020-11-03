package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountryPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoInsertReturnRows(mock, 1, country.Name)
	err = countryPgRep.Insert(country)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_Insert_CountryAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoInsertReturnRows(mock, 1, country.Name)
	err = countryPgRep.Insert(country)
	assert.NoError(t, err)

	mocks.MockCountryRepoInsertReturnErrNoUniq(mock, 2, country.Name)
	err = countryPgRep.Insert(country)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoUpdateReturnResultOk(mock, country.ID, country.Name)
	err = countryPgRep.Update(country)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_Update_NoCountryWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoUpdateReturnResultZero(mock, country.ID, country.Name)
	err = countryPgRep.Update(country)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_DeleteByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoDeleteReturnResultOk(mock, country.ID, country.Name)
	err = countryPgRep.DeleteByID(country.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoSelectByIDReturnRows(mock, country.ID, country.Name)
	dbCountry, err := countryPgRep.SelectByID(country.ID)
	assert.Equal(t, country, dbCountry)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_SelectById_NoCountryWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoSelectByIDReturnErrNoRows(mock, country.ID)
	dbCountry, err := countryPgRep.SelectByID(country.ID)
	assert.Equal(t, dbCountry, (*models.Country)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_SelectByName_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoSelectByNameReturnRows(mock, country.ID, country.Name)
	dbCountry, err := countryPgRep.SelectByName(country.Name)
	assert.Equal(t, country, dbCountry)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_SelectById_NoCountryWithThisName(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	country := &models.Country{
		ID:   1,
		Name: "USA",
	}

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoSelectByNameReturnErrNoRows(mock, country.Name)
	dbCountry, err := countryPgRep.SelectByName(country.Name)
	assert.Equal(t, dbCountry, (*models.Country)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryPgRepository_SelectAll_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

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

	countryPgRep := NewCountryPgRepository(db)

	mocks.MockCountryRepoSelectAllReturnRows(mock, countries)
	dbCountry, err := countryPgRep.SelectAll()
	assert.Equal(t, countries, dbCountry)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
