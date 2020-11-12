package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	userMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func setupSessionHandler(sessionUseCase session.SessionUsecase,
	userUseCase user.UserUsecase, httpMethod string, stringifiedJSON string) (
	echo.Context, *SessionHandler, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(httpMethod, "/api/v1/session",
		strings.NewReader(stringifiedJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sessionHandler := NewSessionHandler(sessionUseCase, userUseCase)
	sessionHandler.Configure(e, nil)
	mw := mwares.NewMiddlewareManager(sessionUseCase, userUseCase)
	sessionHandler.Configure(e, mw)
	return c, sessionHandler, rec
}

func TestSessionHandler_LoginHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mocks.NewMockSessionUsecase(ctrl)
	userUseCase := userMocks.NewMockUserUsecase(ctrl)
	logger.InitLogger("/dev/null", 10)

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	email := "test_user@mail.ru"
	password := "123456"

	request := Request{
		Email:    email,
		Password: password,
	}

	user := &models.User{
		ID:       1,
		Nickname: "test_user",
		Email:    email,
		Password: password,
		Avatar:   "",
		Role:     "user",
	}

	sessionJSON, err := converter.AnyToBytesBuffer(request)
	if err != nil {
		t.Fatal(err)
	}

	c, sessionHandler, rec := setupSessionHandler(sessionUseCase, userUseCase,
		http.MethodPost, sessionJSON.String())
	handleFunc := sessionHandler.LoginHandler()

	userUseCase.
		EXPECT().
		GetByEmail(user.Email).
		Return(user, nil)

	userUseCase.
		EXPECT().
		CheckPassword(user, user.Password).
		Return(nil)

	sessionUseCase.
		EXPECT().
		Create(gomock.Any()).
		Return(nil)

	response := &response.Response{Body: &response.Body{"user": user}}

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

func TestSessionHandler_LogoutHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mocks.NewMockSessionUsecase(ctrl)
	userUseCase := userMocks.NewMockUserUsecase(ctrl)

	session := models.NewSession(3)
	cookie := tools.CreateCookie(session)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/session",
		strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetCookie(cookie)
	c.SetParamNames("sessValue", "userID")
	c.SetParamValues(session.Value, strconv.FormatUint(session.UserID, 10))

	sessionHandler := NewSessionHandler(sessionUseCase, userUseCase)
	mw := mwares.NewMiddlewareManager(sessionUseCase, userUseCase)
	sessionHandler.Configure(e, mw)

	sessionUseCase.
		EXPECT().
		Delete(session.Value).
		Return(nil)

	handleFunc := sessionHandler.LogoutHandler()

	response := &response.Response{Message: "success"}

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
