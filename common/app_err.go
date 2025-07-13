package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"err_key"`
}

func NewAppError(root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}

}
func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

// Contructor
func NewFullErrorResponse(statuCode int, root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: statuCode,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}
func NewErrorResponse(root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewUnAuthorized(root error, msg, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Key:        key,
	}
}

func NewCustomError(root error, msg, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}
	return NewErrorResponse(errors.New(msg), msg, msg, key)
}

// Ham tien ich
func ErrEntityExisted(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s already existed", strings.ToLower(entity)),
		fmt.Sprintf("Err%sAlreadyExisted", entity))
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s not found", strings.ToLower(entity)),
		fmt.Sprintf("Err%sNotFound", entity))
}

func ErrCanNotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s Can Not Create", strings.ToLower(entity)),
		fmt.Sprintf("Err%sCanNotCreate", entity))
}

func ErrNoPermission(err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("You have no permission"),
		fmt.Sprintf("ErrNoPermission"))
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "Something went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "Invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "Something went wrong in the server", err.Error(), "ErrInternal")
}

func ErrCanNotListEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not List %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCanNotList%s", entity))
}
func ErrCanNotGetEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCanNoGet%s", entity))
}

func ErrCanNotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not Delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCanNotDelete%s", entity))
}

func ErrCanNotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not Update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCanNotUpdate%s", entity))
}
func ErrCanNotGeneratePassword(password string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not generate password %s", strings.ToLower(password)),
		fmt.Sprintf("ErrCanNotGeneratePassword%s", password))
}

var RecordNotFound = errors.New("Record not found") // for unit test
