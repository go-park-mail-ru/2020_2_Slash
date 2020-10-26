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
	"reflect"
	"testing"
)

type ChangeProfileTestCase struct {
	name    string
	reqBody map[string]interface{}
	resBody interface{}
	status  int
	cookie  *http.Cookie
	user    *user.User
}

func TestChangeProfileHandler(t *testing.T) {
	t.Parallel()

	method := "GET"
	target := url + "/user/profile"

	// Register user, create session and cookie
	UserHandler := handlers.NewUserHandler()
	newUser := &user.User{
		Nickname: "Oleg",
		Email:    "o.gibadulin@yandex.ru",
		Password: "hardpassword",
	}
	UserHandler.UserRepo.Register(newUser)
	sess := UserHandler.SessionManager.Create(newUser)
	cookie := handlers.CreateCookie(sess)

	anotherUser := &user.User{
		Nickname: "Oleg",
		Email:    "oo.gibadulin@yandex.ru",
		Password: "hardpassword",
	}
	UserHandler.UserRepo.Register(anotherUser)

	userToCompare := &user.User{
		ID:       0,
		Nickname: "Oleg",
		Email:    "o.gibadulin@yandex.ru",
		Password: "hardpassword",
	}

	cases := []ChangeProfileTestCase{
		ChangeProfileTestCase{
			name:    "Empty cookie value",
			reqBody: map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "user isn't authorized",
			},
			status: http.StatusUnauthorized,
			user:   userToCompare,
		},
		ChangeProfileTestCase{
			name: "Try to change with already existed email",
			reqBody: map[string]interface{}{
				"email": "oo.gibadulin@yandex.ru",
			},
			resBody: map[string]interface{}{
				"error": "email already exists",
			},
			status: http.StatusBadRequest,
			cookie: cookie,
			user:   userToCompare,
		},
		ChangeProfileTestCase{
			name: "Try to change with the same email",
			reqBody: map[string]interface{}{
				"email": "o.gibadulin@yandex.ru",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user:   userToCompare,
		},
		ChangeProfileTestCase{
			name:    "Empty requset",
			reqBody: map[string]interface{}{},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user:   userToCompare,
		},
		ChangeProfileTestCase{
			name: "Request with another keys",
			reqBody: map[string]interface{}{
				"nickname 1":  "Alex",
				"email other": "alex.alex@yandex.ru",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user:   userToCompare,
		},
		ChangeProfileTestCase{
			name: "Try to change with wrong email",
			reqBody: map[string]interface{}{
				"email": "a",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user:   userToCompare,
		},
		ChangeProfileTestCase{
			name: "Change one field",
			reqBody: map[string]interface{}{
				"nickname": "Alex",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user: &user.User{
				ID:       0,
				Nickname: "Alex",
				Email:    "o.gibadulin@yandex.ru",
				Password: "hardpassword",
			},
		},
		ChangeProfileTestCase{
			name: "Change correct email",
			reqBody: map[string]interface{}{
				"email": "alex.alex@yandex.ru",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user: &user.User{
				ID:       0,
				Nickname: "Alex",
				Email:    "alex.alex@yandex.ru",
				Password: "hardpassword",
			},
		},
		ChangeProfileTestCase{
			name: "Check if UserRepo store user correctly after prev test",
			reqBody: map[string]interface{}{
				"email": "o.gibadulin@yandex.ru",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user: &user.User{
				ID:       0,
				Nickname: "Alex",
				Email:    "o.gibadulin@yandex.ru",
				Password: "hardpassword",
				Avatar:   "",
			},
		},
		ChangeProfileTestCase{
			name: "Change all fields",
			reqBody: map[string]interface{}{
				"nickname": "Oleg",
				"email":    "oooo.gibadulin@yandex.ru",
			},
			resBody: map[string]interface{}{
				"result": "ok",
			},
			status: http.StatusOK,
			cookie: cookie,
			user: &user.User{
				ID:       0,
				Nickname: "Oleg",
				Email:    "oooo.gibadulin@yandex.ru",
				Password: "hardpassword",
			},
		},
	}

	for _, tc := range cases {
		reqBody := new(bytes.Buffer)
		err := json.NewEncoder(reqBody).Encode(tc.reqBody)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, target, reqBody)
		if tc.cookie != nil {
			r.AddCookie(tc.cookie)
		}
		UserHandler.ChangeUserProfile(w, r)

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
			changedUser, _ := UserHandler.UserRepo.Get(tc.user.ID)
			if !reflect.DeepEqual(tc.user, changedUser) {
				t.Errorf(tc.name + ": exp and act users don't match")
			}
		}
	}
}
