package tests

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type SetAvatarTestCase struct {
	name    string
	resBody map[string]interface{}
	status  int
	image   string
}

func TestSetAvatarHandler(t *testing.T) {
	t.Parallel()

	method := "POST"
	target := url + "/user/avatar"

	storagePath := "./avatars" // TODO: config
	err := os.Mkdir(storagePath, 0777)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(storagePath)

	// Register user, create session and cookie
	UserHandler := handlers.NewUserHandler(db)
	curUser := &user.User{
		Nickname: "Oleg",
		Email:    "o@o.ru",
		Password: "hardpassword",
	}
	UserHandler.UserRepo.Register(curUser)
	sess, err := UserHandler.SessionManager.Create(curUser)
	if err != nil {
		t.Fatal(err)
	}
	cookie := handlers.CreateCookie(sess)

	cases := []SetAvatarTestCase{
		SetAvatarTestCase{
			name: "Empty request",
			resBody: map[string]interface{}{
				"error": "avatar field is expected",
			},
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		if tc.image != "" {
			var reader io.Reader
			var err error
			reader, err = os.Open(tc.image)
			if err != nil {
				t.Error(err)
			}
			var fw io.Writer
			if x, ok := reader.(io.Closer); ok {
				defer x.Close()
			}
			if _, ok := reader.(*os.File); ok {
				if fw, err = writer.CreateFormFile("avatar", "0.png"); err != nil {
					t.Error(err)
				}
			}
			if _, err := io.Copy(fw, reader); err != nil {
				t.Error(err)
			}
		}
		writer.Close()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, target, &b)
		r.AddCookie(cookie)
		r.Header.Set("Content-Type", writer.FormDataContentType())
		UserHandler.SetAvatar(w, r)

		// Check status
		assert := assert.New(t)
		assert.Equal(tc.status, w.Code, tc.name+": wrong status code")

		expResBody := new(bytes.Buffer)
		err := json.NewEncoder(expResBody).Encode(tc.resBody)
		if err != nil {
			t.Error(err)
		}

		// Check responce body
		if w.Code != http.StatusOK {
			res := w.Result()
			defer res.Body.Close()
			actResBody, _ := ioutil.ReadAll(res.Body)
			assert.Equal(expResBody.String(), string(actResBody)+"\n", tc.name+": exp and act resp bodies don't match")
		}
	}
}
