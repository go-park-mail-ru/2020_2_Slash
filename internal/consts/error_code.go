package consts

type ErrorCode uint16

const (
	CodeBadRequest ErrorCode = iota + 101
	CodeInternalError
	CodeUserUnauthorized
	CodeUserDoesNotExist
	CodeTooShortPassword
	CodeEmailDoesNotExist
	CodeWrongPassword
	CodeEmailAlreadyExists
	CodeWrongImgExtension
	CodeSessionDoesNotExist
	CodeSessionExpired
	CodeActorDoesNotExist
	CodeDirectorDoesNotExist
	CodeGenreNameAlreadyExists
	CodeGenreDoesNotExist
)
