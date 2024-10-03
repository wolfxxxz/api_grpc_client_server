package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Message  string
	Code     string
	HTTPCode int
}

func NewAppError() *AppError {
	return &AppError{}
}

var (
	EnvConfigLoadError = AppError{
		Message:  "Failed to parse env file",
		Code:     envParse,
		HTTPCode: http.StatusInternalServerError,
	}
	EnvConfigParseError = AppError{
		Message:  "Failed to parse env file",
		Code:     envParse,
		HTTPCode: http.StatusInternalServerError,
	}
	MongoInitFailedError = AppError{
		Message:  "Failed Init mongoDB",
		Code:     initMongo,
		HTTPCode: http.StatusInternalServerError,
	}
	MongoCreateUserFailedError = AppError{
		Message:  "Failed CreateUser",
		Code:     userRepo,
		HTTPCode: http.StatusInternalServerError,
	}
	MongoGetFailedError = AppError{
		Message:  "Failed decodeUsers",
		Code:     userRepo,
		HTTPCode: http.StatusInternalServerError,
	}
	MongoGetUserByIDErr = AppError{
		Message:  "Failed GetUserByID",
		Code:     userRepo,
		HTTPCode: http.StatusNotFound,
	}
	MongoGetUserByEmailErr = AppError{
		Message:  "Failed GetUserByEmail",
		Code:     userRepo,
		HTTPCode: http.StatusInternalServerError,
	}
	MongoUpdateUserByIDErr = AppError{
		Message:  "Failed UpdateUserByID",
		Code:     userRepo,
		HTTPCode: http.StatusInternalServerError,
	}
	MongoDropUserByIDErr = AppError{
		Message:  "Failed DropUserByID",
		Code:     userRepo,
		HTTPCode: http.StatusInternalServerError,
	}
	HashPasswordErr = AppError{
		Message:  "Hash Password Err",
		Code:     userInterfactor,
		HTTPCode: http.StatusInternalServerError,
	}
	MapCreateUserRequestToUserErr = AppError{
		Message:  "Create User Err",
		Code:     mapCreateUserRequestToUser,
		HTTPCode: http.StatusBadRequest,
	}
	ParsingStringToIntFailed = AppError{
		Message:  "Failed GetUsersByPageAndPerPage",
		Code:     userInterfactor,
		HTTPCode: http.StatusBadRequest,
	}
	CoverstionToUUIDFailed = AppError{
		Message:  "Failed GetUserByID",
		Code:     userInterfactor,
		HTTPCode: http.StatusBadRequest,
	}
	CoverstionToUUIDFailedUpdateUserByID = AppError{
		Message:  "Failed UpdateUserByID CoverstionToUUID",
		Code:     userInterfactor,
		HTTPCode: http.StatusBadRequest,
	}
	CoverstionToUUIDFailedDropUserByID = AppError{
		Message:  "Failed DropUserByID",
		Code:     userInterfactor,
		HTTPCode: http.StatusBadRequest,
	}
	ControllerCreateUserErr = AppError{
		Message:  "Failed CreateUser",
		Code:     controller,
		HTTPCode: http.StatusBadRequest,
	}
	ControllerGetUsersByPaginationError = AppError{
		Message:  "Failed GetUsersByPagination",
		Code:     controller,
		HTTPCode: http.StatusBadRequest,
	}
	ControllerGetUserByIDErr = AppError{
		Message:  "Failed GetUserByID",
		Code:     controller,
		HTTPCode: http.StatusForbidden,
	}
	ControllerUpdateUserByIDErr = AppError{
		Message:  "Failed UpdateUserByID",
		Code:     controller,
		HTTPCode: http.StatusForbidden,
	}
	ControllerUpdateUserByIdErr = AppError{
		Message:  "Failed UpdateUserByID",
		Code:     controller,
		HTTPCode: http.StatusBadRequest,
	}
	ControllerUpdateUserByIDRequestErr = AppError{
		Message:  "Failed UpdateUserByID",
		Code:     controller,
		HTTPCode: http.StatusConflict,
	}
	ControllerDropUserByIDErr = AppError{
		Message:  "Failed DropUserByID",
		Code:     controller,
		HTTPCode: http.StatusBadRequest,
	}
	ControllerGetUserByIdErr = AppError{
		Message:  "Failed GetUserById",
		Code:     controller,
		HTTPCode: http.StatusBadRequest,
	}
	RedisInitErr = AppError{
		Message:  "Failed InitRedisDB",
		Code:     dataStore,
		HTTPCode: http.StatusBadRequest,
	}
	RedisGetUserByIDErr = AppError{
		Message:  "Failed GetUserByID",
		Code:     cacheInterface,
		HTTPCode: http.StatusBadRequest,
	}
	UnmarshalError = AppError{
		Message:  "Failed GetUserByID",
		Code:     cacheInterface,
		HTTPCode: http.StatusBadRequest,
	}
	RedisGetUsersByPageAndPerPageErr = AppError{
		Message:  "Failed GetUserByPageAndPerPage",
		Code:     cacheInterface,
		HTTPCode: http.StatusBadRequest,
	}
	RedisGetUserByEmail = AppError{
		Message:  "Failed GetUserByEmail",
		Code:     cacheInterface,
		HTTPCode: http.StatusBadRequest,
	}
	RedisGetUserByID = AppError{
		Message:  "Failed GetUserByID",
		Code:     cacheInterface,
		HTTPCode: http.StatusBadRequest,
	}
	RedisKeyDoesntExist = AppError{
		Message:  "Failed RedisGet",
		Code:     cacheInterface,
		HTTPCode: http.StatusBadRequest,
	}
)

func (appError *AppError) HttpCode() int {
	return appError.HTTPCode
}

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}

func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message:  fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:     appError.Code,
		HTTPCode: appError.HTTPCode,
	}
}

func IsAppError(err1 error, err2 *AppError) bool {
	err, ok := err1.(*AppError)
	if !ok {
		return false
	}

	return err.Code == err2.Code
}
