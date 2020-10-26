package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
			"avatar":   TestUserAvatar,
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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlmock.MonitorPingsOption(true)
	mock.ExpectPing()

	h := handlers.NewUserHandler(db)

	// for positive test
	mockUserRepoSelectByEmailReturnRows(mock, TestUserID, TestUserNick, TestUserEmail, TestUserPassword, TestUserAvatar)
	mockInsertSessionReturnRows(mock, 1, TestUserID)
	// for wrong email test
	mockUserRepoSelectByEmailReturnErrNoRows(mock, TestUserID, TestUserNick, TestUserEmail[1:], TestUserPassword, TestUserAvatar)
	// for wrong password test
	mockUserRepoSelectByEmailReturnRows(mock, TestUserID, TestUserNick, TestUserEmail, TestUserPassword, TestUserAvatar)

	for _, testCase := range testCases {
		reqBody := new(bytes.Buffer)
		err := json.NewEncoder(reqBody).Encode(testCase.reqBody)
		if err != nil {
			t.Error(err)
			return
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
			return
		}
		bytes, _ := ioutil.ReadAll(w.Body)
		assert.JSONEq(t, expResBody.String(), string(bytes), "bodies doesn't match")
	}
}

func TestLogout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlmock.MonitorPingsOption(true)
	mock.ExpectPing()

	testUser := user.User{
		ID:       TestUserID,
		Nickname: TestUserNick,
		Email:    TestUserEmail,
		Password: TestUserPassword,
		Avatar:   TestUserAvatar,
	}
	testCase := buildOkTestCase()
	reqBody, err := bodyToBytesBuffer(testCase.reqBody)
	if err != nil {
		t.Error(err)
	}

	h := handlers.NewUserHandler(db)
	r := httptest.NewRequest("POST", "/api/v1/user/login", reqBody)
	w := httptest.NewRecorder()

	mockUserRepoSelectByEmailReturnRows(mock, TestUserID, TestUserNick, TestUserEmail, TestUserPassword, TestUserAvatar)
	mockInsertSessionReturnRows(mock, 1, TestUserID)

	h.Login(w, r)

	assert.NotEmpty(t, w.Result().Cookies())
	sess := w.Result().Cookies()[0]
	mockGetUserSessionReturnRows(mock, sess.Value, 1, sess.Expires, 1)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM session").WithArgs(sess.Value).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r = httptest.NewRequest("DELETE", "/api/v1/user/logout", http.NoBody)
	r.AddCookie(sess)
	w = httptest.NewRecorder()

	h.Logout(w, r)

	if w.Result().Cookies()[0].Expires.After(time.Now()) {
		t.Error("Session must expire in past")
	}
	mock.ExpectQuery("SELECT id, value, expires, profile_id " +
		"FROM session").WithArgs(sess).WillReturnError(sql.ErrNoRows)
	ok, err := h.SessionManager.IsAuthorized(&testUser)
	if err != nil {
		t.Error(err)
	}
	if ok == true {
		t.Error("User should be unauthorized")
	}
}
