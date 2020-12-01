package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (am *AdminMicroservice) CreateCountry(ctx context.Context, country *Country) (*Country, error) {
	if err := am.checkCountryByName(country.Name); err == nil {
		return &Country{}, status.Error(codes.Code(consts.CodeCountryNameAlreadyExists), "")
	}

	modelCountry := CountryGRPCToModel(country)
	if err := am.countriesRep.Insert(modelCountry); err != nil {
		return &Country{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	country.ID = modelCountry.ID

	return country, nil
}

func (am *AdminMicroservice) ChangeCountry(ctx context.Context,
	newCountryData *Country) (*empty.Empty, error) {
	modelCountry, err := am.GetCountryByID(newCountryData.GetID())
	if err != nil {
		return &empty.Empty{}, err
	}

	if newCountryData.Name == "" || modelCountry.Name == newCountryData.Name {
		return &empty.Empty{}, nil
	}

	if err := am.checkCountryByName(newCountryData.Name); err == nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeCountryNameAlreadyExists), "")
	}

	modelCountry.Name = newCountryData.Name
	if err := am.countriesRep.Update(modelCountry); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteCountryByID(ctx context.Context, countryID *ID) (*empty.Empty, error) {
	if err := am.checkCountryByID(countryID.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeCountryDoesNotExist), "")
	}

	if err := am.countriesRep.DeleteByID(countryID.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) GetCountryByID(countryID uint64) (*models.Country, error) {
	country, err := am.countriesRep.SelectByID(countryID)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeCountryDoesNotExist), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return country, nil
}

func (am *AdminMicroservice) GetCountryByName(name string) (*models.Country, error) {
	country, err := am.countriesRep.SelectByName(name)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeCountryNameAlreadyExists), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return country, nil
}

func (am *AdminMicroservice) checkCountryByID(countryID uint64) error {
	_, err := am.GetCountryByID(countryID)
	return err
}

func (am *AdminMicroservice) checkCountryByName(name string) error {
	_, err := am.GetCountryByName(name)
	return err
}
