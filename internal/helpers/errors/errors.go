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
	CodeContentDoesNotExist: {
		Code:        CodeContentDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "content does not exist",
		UserMessage: "Такого контента не существует",
	},
	CodeMovieContentAlreadyExists: {
		Code:        CodeMovieContentAlreadyExists,
		HTTPCode:    http.StatusBadRequest,
		Message:     "movie with this content already exists",
		UserMessage: "Данный контент фильма уже существует",
	},
	CodeMovieDoesNotExist: {
		Code:        CodeMovieDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "movie does not exist",
		UserMessage: "Такого фильма не существует",
	},
	CodeRatingDoesNotExist: {
		Code:        CodeRatingDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "trying to update rating that not exist",
		UserMessage: "Что-то пошло не так",
	},
	CodeRatingAlreadyExist: {
		Code:        CodeRatingAlreadyExist,
		HTTPCode:    http.StatusConflict,
		Message:     "rating already exist",
		UserMessage: "Оценка уже поставлена",
	},
	CodeFavouriteAlreadyExist: {
		Code:        CodeFavouriteAlreadyExist,
		HTTPCode:    http.StatusConflict,
		Message:     "this favourite already exist",
		UserMessage: "Уже в избранном",
	},
	CodeFavouriteDoesNotExist: {
		Code:        CodeFavouriteDoesNotExist,
		HTTPCode:    http.StatusBadRequest,
		Message:     "this favourite does not exist",
		UserMessage: "Данный контент отсутствует в избранном",
	},
	CodeAccessDenied: {
		Code:        CodeAccessDenied,
		HTTPCode:    http.StatusUnauthorized,
		Message:     "user tries to execute request without sufficient rights",
		UserMessage: "Недостаточно прав",
	},
	CodeCSRFTokenWasNotPassed: {
		Code:        CodeCSRFTokenWasNotPassed,
		HTTPCode:    http.StatusBadRequest,
		Message:     "CSRF token was not passed",
		UserMessage: "Неверный формат запроса",
	},
	CodeWrongCSRFToken: {
		Code:        CodeWrongCSRFToken,
		HTTPCode:    http.StatusBadRequest,
		Message:     "wrong CSRF token",
		UserMessage: "Неверный формат запроса",
	},
	CodeErrorInNickname: {
		Code:        CodeErrorInNickname,
		HTTPCode:    http.StatusBadRequest,
		Message:     "error in nickname",
		UserMessage: "Логин должен содержать от 3 до 32 символов",
	},
	CodeErrorInEmail: {
		Code:        CodeErrorInEmail,
		HTTPCode:    http.StatusBadRequest,
		Message:     "error in email",
		UserMessage: "Email должен быть корректным и содержать до 64 символов",
	},
	CodeErrorInPassword: {
		Code:        CodeErrorInPassword,
		HTTPCode:    http.StatusBadRequest,
		Message:     "error in password",
		UserMessage: "Пароль должен содержать от 6 до 32 символов",
	},
	CodePasswordsDoesNotMatch: {
		Code:        CodePasswordsDoesNotMatch,
		HTTPCode:    http.StatusBadRequest,
		Message:     "passwords does not match",
		UserMessage: "Пароли не совпадают",
	},
}
