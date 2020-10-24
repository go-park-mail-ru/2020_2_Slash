package tests

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetProfileTestCase struct {
	name     string
	resBody  interface{}
	status   int
	cookie   *http.Cookie
	mockFunc func()
}

func TestGetProfileHandler(t *testing.T) {
	t.Parallel()

	method := "GET"
	target := url + "/user/profile"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	// Register user, create session and cookie
	UserHandler := handlers.NewUserHandler(db)
	newUser := &user.User{
		Nickname: "Oleg",
		Email:    "o@o.ru",
		Password: "hardpassword",
	}
	mockCheckAndInsertUser(mock, 1, newUser.Nickname, newUser.Email, newUser.Password, newUser.Avatar)
	UserHandler.UserRepo.Register(newUser)

	mockInsertSessionReturnRows(mock, newUser.ID, 1)
	sess, err := UserHandler.SessionManager.Create(newUser)
	if err != nil {
		t.Fatal(err)
	}

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

	mockGetProfileSuccess(mock, cookie.Value, 1, cookie.Expires, 1, "Oleg", "o@o.ru", "hardpassword", "")
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
