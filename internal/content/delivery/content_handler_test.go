package delivery

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestContentHandler_UpdateContentPostersHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)

	e := echo.New()
	strId := strconv.Itoa(1)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/content"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	contentHandler := NewContentHandler(contentUseCase)
	handleFunc := contentHandler.UpdatePostersHandler()
	contentHandler.Configure(e, nil)

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
