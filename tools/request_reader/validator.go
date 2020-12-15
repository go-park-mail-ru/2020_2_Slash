package request_reader

import (
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"net/http"
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

func (rr *RequestReader) ReadUser(request interface{}) *errors.Error {
	if err := rr.cntx.Bind(request); err != nil {
		return errors.New(CodeInternalError, err)
	}

	if err := rr.validator.Struct(request); err != nil {
		// nolint: errcheck, errorlint
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Nickname":
				return errors.New(CodeErrorInNickname, err)
			case "Email":
				return errors.New(CodeErrorInEmail, err)
			case "Password":
				return errors.New(CodeErrorInPassword, err)
			case "RepeatedPassword":
				return errors.New(CodePasswordsDoesNotMatch, err)
			}
		}
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
	if customErr := helpers.CheckImageContentType(image); customErr != nil {
		return nil, customErr
	}
	return image, nil
}

func (rr *RequestReader) ReadNotRequiredImage(field string) (*multipart.FileHeader, *errors.Error) {
	image, err := rr.cntx.FormFile(field)
	switch {
	case err == http.ErrMissingFile:
		return nil, nil
	case err != nil:
		return nil, errors.New(CodeBadRequest, err)
	}

	// Check content type of image
	if customErr := helpers.CheckImageContentType(image); customErr != nil {
		return nil, customErr
	}
	return image, nil
}

func (rr *RequestReader) ReadVideo(field string) (*multipart.FileHeader, *errors.Error) {
	video, err := rr.cntx.FormFile(field)
	if err != nil {
		return nil, errors.New(CodeBadRequest, err)
	}

	// Check content type of video
	if customErr := helpers.CheckVideoContentType(video); customErr != nil {
		return nil, customErr
	}
	return video, nil
}
