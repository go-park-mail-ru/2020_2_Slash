package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	sessMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/session/mocks"
	userMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserHandler_CreateUserHandler_EmailAlreadyExists(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := userMocks.NewMockUserUsecase(ctrl)
	sessUseCase := sessMocks.NewMockSessionUsecase(ctrl)

	type Request struct {
		Nickname         string `json:"nickname" validate:"gte=3,lte=32"`
		Email            string `json:"email" validate:"required,email,lte=64"`
		Password         string `json:"password" validate:"required,gte=6,lte=32"`
		RepeatedPassword string `json:"repeated_password" validate:"eqfield=Password"`
	}

	var reqInst = &Request{
		Nickname:         "Jhon",
		Email:            "jhon@gmail.com",
		Password:         "hardpassword",
		RepeatedPassword: "hardpassword",
	}

	var userInst = &models.User{
		Nickname: "Jhon",
		Email:    "jhon@gmail.com",
		Password: "hardpassword",
		Role:     consts.User,
	}

	userJSON, err := converter.AnyToBytesBuffer(reqInst)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", strings.NewReader(userJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	userHandler := NewUserHandler(userUseCase, sessUseCase)
	handleFunc := userHandler.RegisterUserHandler()
	userHandler.Configure(e, nil)

	userUseCase.
		EXPECT().
		Create(userInst).
		Return(errors.Get(consts.CodeEmailAlreadyExists))

	response := &response.Response{Error: errors.Get(consts.CodeEmailAlreadyExists)}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestUserHandler_UpdateUserHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := userMocks.NewMockUserUsecase(ctrl)
	sessUseCase := sessMocks.NewMockSessionUsecase(ctrl)

	type Request struct {
		Nickname string `json:"nickname" validate:"omitempty,gte=3,lte=32"`
		Email    string `json:"email" validate:"omitempty,email,lte=64"`
		Password string `json:"password" validate:"omitempty,gte=6,lte=32"`
	}

	var reqInst = &Request{
		Nickname: "Jhon",
		Email:    "jhonJhon@gmail.com",
		Password: "hardpassword",
	}

	var userInst = &models.User{
		Nickname: "Jhon",
		Email:    "jhonJhon@gmail.com",
		Password: "hardpassword",
		Role:     consts.User,
	}

	userJSON, err := converter.AnyToBytesBuffer(reqInst)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/profile", strings.NewReader(userJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userInst.ID)

	userHandler := NewUserHandler(userUseCase, sessUseCase)
	handleFunc := userHandler.UpdateUserProfileHandler()
	userHandler.Configure(e, nil)

	userUseCase.
		EXPECT().
		UpdateProfile(userInst).
		Return(userInst, nil)

	response := &response.Response{Body: &response.Body{"user": userInst}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestUserHandler_GetUserProfileHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := userMocks.NewMockUserUsecase(ctrl)
	sessUseCase := sessMocks.NewMockSessionUsecase(ctrl)

	var userInst = &models.User{
		Nickname: "Jhon",
		Email:    "jhonJhon@gmail.com",
		Password: "hardpassword",
		Role:     consts.User,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/profile", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userInst.ID)

	userHandler := NewUserHandler(userUseCase, sessUseCase)
	handleFunc := userHandler.GetUserProfileHandler()
	userHandler.Configure(e, nil)

	userUseCase.
		EXPECT().
		GetByID(userInst.ID).
		Return(userInst, nil)

	response := &response.Response{Body: &response.Body{"user": userInst}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestUserHandler_UpdateAvatarHandler_Fail(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := userMocks.NewMockUserUsecase(ctrl)
	sessUseCase := sessMocks.NewMockSessionUsecase(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/avatar", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userHandler := NewUserHandler(userUseCase, sessUseCase)
	handleFunc := userHandler.UpdateAvatarHandler()
	userHandler.Configure(e, nil)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":         101,
			"message":      "request Content-Type isn't multipart/form-data",
			"user_message": "Неверный формат запроса",
		},
	}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}
