package tests

import (
	"errors"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

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
				Email:    "o.gibadulin@yandex.ru",
			},
			err: errors.New("not enough input data"),
		},
		CreateUserTestCase{
			name: "Passwords that don't match",
			userInput: &handlers.UserInput{
				Nickname:         "Oleg",
				Email:            "o.gibadulin@yandex.ru",
				Password:         "hardpassword",
				RepeatedPassword: "otherpassword",
			},
			err: errors.New("passwords don't match"),
		},
		CreateUserTestCase{
			name: "Correct data",
			userInput: &handlers.UserInput{
				Nickname:         "Oleg",
				Email:            "o.gibadulin@yandex.ru",
				Password:         "hardpassword",
				RepeatedPassword: "hardpassword",
			},
			user: &user.User{
				ID:       0,
				Nickname: "Oleg",
				Email:    "o.gibadulin@yandex.ru",
				Password: "hardpassword",
			},
		},
		CreateUserTestCase{
			name: "Correct data with empty nickname",
			userInput: &handlers.UserInput{
				Email:            "o.gibadulin@yandex.ru",
				Password:         "hardpassword",
				RepeatedPassword: "hardpassword",
			},
			user: &user.User{
				ID:       0,
				Nickname: "o.gibadulin",
				Email:    "o.gibadulin@yandex.ru",
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
