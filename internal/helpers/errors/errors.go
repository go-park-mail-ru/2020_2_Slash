package errors

import (
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"net/http"
)

type Error struct {
	Code        ErrorCode `json:"code"`
	HTTPCode    int       `json:"-"`
	Message     string    `json:"message"`
	UserMessage string    `json:"user_message"`
}

var WrongErrorCode = &Error{
	HTTPCode:    http.StatusTeapot,
	Message:     "wrong error code",
	UserMessage: "Что-то пошло не так",
}

func New(code ErrorCode, err error) *Error {
	customErr, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	customErr.Message = err.Error()
	return customErr
}

func Get(code ErrorCode) *Error {
	err, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	return err
}

var Errors = map[ErrorCode]*Error{
	CodeBadRequest: {
		Code:        CodeBadRequest,
		HTTPCode:    http.StatusBadRequest,
		Message:     "wrong request data",
		UserMessage: "Неверный формат запроса",
	},
	CodeInternalError: {
		Code:        CodeInternalError,
		HTTPCode:    http.StatusInternalServerError,
		Message:     "something went wrong",
		UserMessage: "Что-то пошло не так",
	},
	CodeEmailAlreadyExists: {
		Code:        CodeEmailAlreadyExists,
		HTTPCode:    http.StatusBadRequest,
		Message:     "user with this email already exists",
		UserMessage: "Данный Email адрес уже существует",
	},
	CodeUserUnauthorized: {
		Code:        CodeUserUnauthorized,
		HTTPCode:    http.StatusUnauthorized,
		Message:     "user is unauthorized",
		UserMessage: "Вы не авторизированы",
	},
	CodeUserDoesNotExist: {
		Code:        CodeUserDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "user does not exist",
		UserMessage: "Такого пользователя не существует",
	},
	CodeWrongImgExtension: {
		Code:        CodeWrongImgExtension,
		HTTPCode:    http.StatusBadRequest,
		Message:     "file with this extension is prohibited",
		UserMessage: "Файлы с данным расширением запрещены",
	},
	CodeWrongPassword: {
		Code:        CodeWrongPassword,
		HTTPCode:    http.StatusBadRequest,
		Message:     "entered password don't match with saved password",
		UserMessage: "Неверный пароль",
	},
	CodeSessionDoesNotExist: {
		Code:        CodeSessionDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "user session doesn't exist in db",
		UserMessage: "Сессия невалидна",
	},
	CodeSessionExpired: {
		Code:        CodeSessionExpired,
		HTTPCode:    http.StatusUnauthorized,
		Message:     "session expired",
		UserMessage: "Сессия истекла",
	},
	CodeEmailDoesNotExist: {
		Code:        CodeEmailDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "email doesn't exist in db",
		UserMessage: "Пользователь с таким email не найден",
	},
	CodeGenreNameAlreadyExists: {
		Code:        CodeGenreNameAlreadyExists,
		HTTPCode:    http.StatusBadRequest,
		Message:     "genre with this name already exists",
		UserMessage: "Данный жанр уже существует",
	},
	CodeGenreDoesNotExist: {
		Code:        CodeGenreDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "genre does not exist",
		UserMessage: "Такого жанра не существует",
	},
	CodeActorDoesNotExist: {
		Code:        CodeActorDoesNotExist,
		HTTPCode:    http.StatusNotFound,
		Message:     "actor doesn't exist in db",
		UserMessage: "Такого актёра не существует",
	},
	CodeDirectorDoesNotExist: {
		Code:        CodeDirectorDoesNotExist,
		HTTPCode:    http.StatusNotFound,
		Message:     "director doesn't exist in db",
		UserMessage: "Такого режиссёра не существует",
	},
	CodeCountryNameAlreadyExists: {
		Code:        CodeCountryNameAlreadyExists,
		HTTPCode:    http.StatusBadRequest,
		Message:     "country with this name already exists",
		UserMessage: "Данная страна уже существует",
	},
	CodeCountryDoesNotExist: {
		Code:        CodeCountryDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "country does not exist",
		UserMessage: "Такой страны не существует",
	},
}
