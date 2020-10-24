package helpers

import (
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
)

type Error struct {
	Code        int    `json:"code"`
	UserMessage string `json:"user_message"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Err *Error `json:"error"`
}

func NewErrorResponce(code int, err error) *ErrorResponse {
	return &ErrorResponse{
		Err: GetCustomError(code, err),
	}
}

var CustomErrors = map[int]*Error{
	CodeBadRequest: {
		Code:        CodeBadRequest,
		UserMessage: "Неверный формат запроса",
	},
	CodeInternalError: {
		Code:        CodeInternalError,
		UserMessage: "Что-то пошло не так",
	},
	CodeUserUnauthorized: {
		Code:        CodeUserUnauthorized,
		UserMessage: "Вы не авторизированы",
	},
	CodeUserDoesNotExist: {
		Code:        CodeUserDoesNotExist,
		UserMessage: "Такого пользователя не существует",
	},
	CodeTooShortPassword: {
		Code:        CodeTooShortPassword,
		UserMessage: "Пароль должен быть длинне 5 символов",
	},
	CodeEmailDoesNotExist: {
		Code:        CodeEmailDoesNotExist,
		UserMessage: "Такого Email не существует",
	},
	CodeWrongPassword: {
		Code:        CodeWrongPassword,
		UserMessage: "Неверный пароль",
	},
	CodeEmailAlreadyExists: {
		Code:        CodeEmailAlreadyExists,
		UserMessage: "Данный Email адрес уже существует",
	},
	CodeWrongImgExtension: {
		Code:        CodeWrongImgExtension,
		UserMessage: "Файлы с таким расширением запрещены",
	},
}

func GetCustomError(code int, err error) *Error {
	customError := CustomErrors[code]
	customError.Message = err.Error()
	return customError
}
