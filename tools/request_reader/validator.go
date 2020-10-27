package request_reader

import (
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"mime/multipart"
)

type RequestReader struct {
	cntx      echo.Context
	validator *validator.Validate
}

func NewRequestReader(cntx echo.Context) *RequestReader {
	return &RequestReader{
		cntx:      cntx,
		validator: validator.New(),
	}
}

func (rr *RequestReader) Read(request interface{}) *errors.Error {
	if err := rr.cntx.Bind(request); err != nil {
		return errors.New(CodeInternalError, err)
	}

	if err := rr.validator.Struct(request); err != nil {
		return errors.New(CodeBadRequest, err)
	}
	return nil
}

func (rr *RequestReader) ReadImage(field string) (*multipart.FileHeader, *errors.Error) {
	image, err := rr.cntx.FormFile(field)
	if err != nil {
		return nil, errors.New(CodeBadRequest, err)
	}

	// Check content type of image
	if customErr := helpers.CheckImageContentType(image); err != nil {
		return nil, customErr
	}
	return image, nil
}
