package consts

import "errors"

const (
	CodeBadRequest         = 101
	CodeInternalError      = 102
	CodeUserUnauthorized   = 103
	CodeUserDoesNotExist   = 104
	CodeTooShortPassword   = 105
	CodeEmailDoesNotExist  = 106
	CodeWrongPassword      = 107
	CodeEmailAlreadyExists = 108
	CodeWrongImgExtension  = 109
)

var Errors = map[int]error{
	CodeBadRequest:         errors.New("wrong request data"),
	CodeInternalError:      errors.New("internal server error"),
	CodeUserUnauthorized:   errors.New("user is unauthorized"),
	CodeUserDoesNotExist:   errors.New("user does not exist"),
	CodeTooShortPassword:   errors.New("password too small"),
	CodeEmailDoesNotExist:  errors.New("this email does not exist"),
	CodeWrongPassword:      errors.New("wrong password"),
	CodeEmailAlreadyExists: errors.New("user with this email already exists"),
	CodeWrongImgExtension:  errors.New("file with this extension is prohibited"),
}
