package tests

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoginTestCase struct {
	name    string
	reqBody map[string]interface{}
	resBody map[string]interface{}
	status  uint64
}

const TestUserID = 0
const TestUserNick = "test"
const TestUserEmail = "test@mail.ru"
const TestUserPassword = "love"
const TestUserAvatar = ""

func buildOkTestCase() LoginTestCase {
	return LoginTestCase{
		name: "ok",
		reqBody: map[string]interface{}{
			"email":    TestUserEmail,
			"password": TestUserPassword,
		},
		resBody: map[string]interface{}{
			"id":       TestUserID,
			"nickname": TestUserNick,
			"avatar": TestUserAvatar,
		},
		status: http.StatusOK,
	}
}

func bodyToBytesBuffer(body map[string]interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func TestLogin(t *testing.T) {
	testCases := []LoginTestCase{
		buildOkTestCase(),
		{
			name: "wrong email",
			reqBody: map[string]interface{}{
				"email":    TestUserEmail[1:],
				"password": TestUserPassword,
			},
			resBody: map[string]interface{}{
				"error": handlers.WrongEmailMsg,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "wrong password",
			reqBody: map[string]interface{}{
				"email":    TestUserEmail,
				"password": TestUserPassword[1:],
			},
			resBody: map[string]interface{}{
				"error": handlers.WrongPasswordMsg,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			reqBody: map[string]interface{}{
				"email":    "t.2",
				"password": "xyz",
			},
			resBody: map[string]interface{}{
				"error": handlers.InvalidEmailMsg,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "empty email",
			reqBody: map[string]interface{}{
				"email":    "",
				"password": "xyz",
			},
			resBody: map[string]interface{}{
				"error": handlers.EmptyEmailMsg,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "empty password",
			reqBody: map[string]interface{}{
				"email":    "test@mail.ru",
				"password": "",
			},
			resBody: map[string]interface{}{
				"error": handlers.EmptyPasswordMsg,
			},
			status: http.StatusBadRequest,
		},
	}

	t.Parallel()

	h := handlers.NewUserHandler()
	h.UserRepo.Register(&user.User{
		ID:       TestUserID,
		Nickname: TestUserNick,
		Email:    TestUserEmail,
		Password: TestUserPassword,
	})

	for _, testCase := range testCases {
		reqBody := new(bytes.Buffer)
		err := json.NewEncoder(reqBody).Encode(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		r := httptest.NewRequest("POST", "/api/v1/user/login", reqBody)
		w := httptest.NewRecorder()

		h.Login(w, r)

		assert.Equal(t, testCase.status, uint64(w.Code), "codes doesn't match")
		if testCase.status == http.StatusOK {
			assert.NotEmpty(t, w.Result().Cookies(), "should be session")
		}

		expResBody, err := bodyToBytesBuffer(testCase.resBody)
		if err != nil {
			t.Error(err)
		}
		bytes, _ := ioutil.ReadAll(w.Body)
		assert.JSONEq(t, expResBody.String(), string(bytes), "bodies doesn't match")
	}
}
