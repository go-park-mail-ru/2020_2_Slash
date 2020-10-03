package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
)

type TestCase struct {
	name    string
	reqBody map[string]interface{}
	resBody interface{}
	status  int
	user    *user.User
}

var url string = "/api/v1"

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

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
				"email":    "o@o.ru",
			},
			resBody: map[string]interface{}{
				"error": "not enough input data",
			},
			status: http.StatusBadRequest,
		},
		TestCase{
			name: "Passwords that don't match",
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
				"email":             "o@o.ru",
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
				"email":             "o@o.ru",
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
				Email:    "o@o.ru",
				Password: "hardpassword",
			},
		},
		TestCase{
			name: "Correct request body with empty nickname",
			reqBody: map[string]interface{}{
				"email":             "oo@o.ru",
				"password":          "hardpassword",
				"repeated_password": "hardpassword",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusCreated,
			user: &user.User{
				ID:       1,
				Nickname: "oo",
				Email:    "oo@o.ru",
				Password: "hardpassword",
			},
		},
		TestCase{
			name: "User already exists",
			reqBody: map[string]interface{}{
				"nickname":          "Oleg",
				"email":             "o@o.ru",
				"password":          "hardpassword",
				"repeated_password": "hardpassword",
			},
			resBody: map[string]interface{}{
				"error": "User with this Email already exists",
			},
			status: http.StatusConflict,
		},
	}

	UserHandler := handlers.NewUserHandler()
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

		// Check responce body
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

type CreateUserTestCase struct {
	name      string
	userInput *handlers.UserInput
	user      *user.User
	err       error
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	cases := []CreateUserTestCase{
		CreateUserTestCase{
			name: "Empty Email",
			userInput: &handlers.UserInput{
				Nickname:         "Oleg",
				Password:         "hardpassword",
				RepeatedPassword: "hardpassword",
			},
			err: errors.New("not enough input data"),
		},
		CreateUserTestCase{
			name: "Empty Password",
			userInput: &handlers.UserInput{
				Nickname: "Oleg",
				Email:    "o@o.ru",
			},
			err: errors.New("not enough input data"),
		},
		CreateUserTestCase{
			name: "Passwords that don't match",
			userInput: &handlers.UserInput{
				Nickname:         "Oleg",
				Email:            "o@o.ru",
				Password:         "hardpassword",
				RepeatedPassword: "otherpassword",
			},
			err: errors.New("passwords don't match"),
		},
		CreateUserTestCase{
			name: "Correct data",
			userInput: &handlers.UserInput{
				Nickname:         "Oleg",
				Email:            "o@o.ru",
				Password:         "hardpassword",
				RepeatedPassword: "hardpassword",
			},
			user: &user.User{
				ID:       0,
				Nickname: "Oleg",
				Email:    "o@o.ru",
				Password: "hardpassword",
			},
		},
		CreateUserTestCase{
			name: "Correct data with empty nickname",
			userInput: &handlers.UserInput{
				Email:            "o@o.ru",
				Password:         "hardpassword",
				RepeatedPassword: "hardpassword",
			},
			user: &user.User{
				ID:       0,
				Nickname: "o",
				Email:    "o@o.ru",
				Password: "hardpassword",
			},
		},
	}

	for _, tc := range cases {
		user, err := handlers.CreateUser(tc.userInput)

		// Check created user
		if !reflect.DeepEqual(tc.user, user) {
			t.Errorf(tc.name + ": exp and act users don't match")
		}

		// Check error
		assert := assert.New(t)
		assert.Equal(tc.err, err, tc.name+": exp and act errors don't match")
	}
}

type GetProfileTestCase struct {
	name    string
	resBody interface{}
	status  int
	cookie  *http.Cookie
}

func TestGetProfileHandler(t *testing.T) {
	t.Parallel()

	method := "GET"
	target := url + "/user/profile"

	// Register user, create session and cookie
	UserHandler := handlers.NewUserHandler()
	newUser := &user.User{
		Nickname: "Oleg",
		Email:    "o@o.ru",
		Password: "hardpassword",
	}
	UserHandler.UserRepo.Register(newUser)
	sess := UserHandler.SessionManager.Create(newUser)
	cookie := handlers.CreateCookie(sess)

	cases := []GetProfileTestCase{
		GetProfileTestCase{
			name: "Empty cookie value",
			resBody: map[string]interface{}{
				"error": "user isn't authorized",
			},
			status: http.StatusUnauthorized,
		},
		GetProfileTestCase{
			name: "Correct request",
			resBody: &user.UserProfile{
				Nickname: "Oleg",
				Email:    "o@o.ru",
				Avatar:   "",
			},
			status: http.StatusOK,
			cookie: cookie,
		},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, target, nil)
		if tc.cookie != nil {
			r.AddCookie(tc.cookie)
		}
		UserHandler.GetUserProfile(w, r)

		// Check status
		assert := assert.New(t)
		assert.Equal(tc.status, w.Code, tc.name+": wrong status code")

		expResBody := new(bytes.Buffer)
		err := json.NewEncoder(expResBody).Encode(tc.resBody)
		if err != nil {
			t.Error(err)
		}

		// Check responce body
		res := w.Result()
		defer res.Body.Close()
		actResBody, _ := ioutil.ReadAll(res.Body)
		assert.Equal(expResBody.String(), string(actResBody)+"\n", tc.name+": exp and act resp bodies don't match")
	}
}
