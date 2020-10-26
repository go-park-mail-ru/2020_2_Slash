package tests

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	method := "POST"
	target := url + "/user/register"

	cases := []TestCase{
		TestCase{
			name:    "Empty request body",
			reqBody: map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "not enough input data",
			},
			status: http.StatusBadRequest,
		},
		TestCase{
			name: "Empty Email",
			reqBody: map[string]interface{}{
				"nickname":          "Oleg",
				"password":          "hardpassword",
				"repeated_password": "hardpassword",
			},
			resBody: map[string]interface{}{
				"error": "not enough input data",
			},
			status: http.StatusBadRequest,
		},
		TestCase{
			name: "Empty Password",
			reqBody: map[string]interface{}{
				"nickname": "Oleg",
				"email":    "o.gibadulin@yandex.ru",
			},
			resBody: map[string]interface{}{
				"error": "not enough input data",
			},
			status: http.StatusBadRequest,
		},
		TestCase{
			name: "Invalid email",
			reqBody: map[string]interface{}{
				"nickname":          "Oleg",
				"email":             "o",
				"password":          "hardpassword",
				"repeated_password": "hardpassword",
			},
			resBody: map[string]interface{}{
				"error": "email is invalid",
			},
			status: http.StatusBadRequest,
		},
		TestCase{
			name: "Passwords that don't match",
			reqBody: map[string]interface{}{
				"nickname":          "Oleg",
				"email":             "o.gibadulin@yandex.ru",
				"password":          "hardpassword",
				"repeated_password": "otherpassword",
			},
			resBody: map[string]interface{}{
				"error": "passwords don't match",
			},
			status: http.StatusBadRequest,
		},
		TestCase{
			name: "Correct request body",
			reqBody: map[string]interface{}{
				"nickname":          "Oleg",
				"email":             "o.gibadulin@yandex.ru",
				"password":          "hardpassword",
				"repeated_password": "hardpassword",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusCreated,
			user: &user.User{
				ID:       0,
				Nickname: "Oleg",
				Email:    "o.gibadulin@yandex.ru",
				Password: "hardpassword",
			},
		},
		TestCase{
			name: "Correct request body with empty nickname",
			reqBody: map[string]interface{}{
				"email":             "oo.gibadulin@yandex.ru",
				"password":          "hardpassword",
				"repeated_password": "hardpassword",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusCreated,
			user: &user.User{
				ID:       1,
				Nickname: "oo.gibadulin",
				Email:    "oo.gibadulin@yandex.ru",
				Password: "hardpassword",
			},
		},
		TestCase{
			name: "User already exists",
			reqBody: map[string]interface{}{
				"nickname":          "Oleg",
				"email":             "o.gibadulin@yandex.ru",
				"password":          "hardpassword",
				"repeated_password": "hardpassword",
			},
			resBody: map[string]interface{}{
				"error": "User with this Email already exists",
			},
			status: http.StatusConflict,
		},
	}

	db, mock, err := sqlmock.New()
	mockCheckAndInsertUser(mock, 0, "Oleg", "o.gibadulin@yandex.ru", "hardpassword", "")
	mockCheckAndInsertSession(mock, 0, 1)
	mockUserRepoGetRows(mock, 0, "Oleg", "o.gibadulin@yandex.ru", "hardpassword", "")

	mockCheckAndInsertUser(mock, 1, "oo.gibadulin", "oo.gibadulin@yandex.ru", "hardpassword", "")
	mockCheckAndInsertSession(mock, 1, 2)
	mockUserRepoGetRows(mock, 1, "oo.gibadulin", "oo.gibadulin@yandex.ru", "hardpassword", "")

	mockAlreadyExistTest(mock, 0, "Oleg", "o.gibadulin@yandex.ru", "hardpassword", "")
	err = mock.ExpectationsWereMet()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	UserHandler := handlers.NewUserHandler(db)
	for _, tc := range cases {
		reqBody := new(bytes.Buffer)
		err := json.NewEncoder(reqBody).Encode(tc.reqBody)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, target, reqBody)
		UserHandler.Register(w, r)

		// Check status
		assert := assert.New(t)
		assert.Equal(tc.status, w.Code, tc.name+": wrong status code")

		expResBody := new(bytes.Buffer)
		err = json.NewEncoder(expResBody).Encode(tc.resBody)
		if err != nil {
			t.Error(err)
		}

		// Check response body
		res := w.Result()
		defer res.Body.Close()
		actResBody, _ := ioutil.ReadAll(res.Body)
		assert.Equal(expResBody.String(), string(actResBody)+"\n", tc.name+": exp and act resp bodies don't match")

		// Check user
		if tc.user != nil {
			createdUser, _ := UserHandler.UserRepo.Get(tc.user.ID)
			if !reflect.DeepEqual(tc.user, createdUser) {
				t.Errorf(tc.name + ": exp and act users don't match")
			}
		}
	}
}
