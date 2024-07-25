package apperror

import (
	"database/sql"
	"errors"
	"net/http"
)

type Error struct {
	error          error //original error
	httpStatusCode int
}

func Wrap(err error) *Error {

	if err == nil {
		err = errors.New("nil error")
	}

	appError, ok := err.(*Error)
	if !ok {
		appError = &Error{
			error: err,
		}
	}

	return appError
}

func (e *Error) Error() string {
	return e.error.Error()
}

func (e *Error) HttpStatusCode() int {
	if e.httpStatusCode != 0 {
		return e.httpStatusCode
	}
	return http.StatusBadRequest
}

func (e *Error) SetHttpStatusCode(httpStatusCode int) *Error {
	e.httpStatusCode = httpStatusCode
	return e
}

func NewDatabaseError(err error) *Error {

	appError := Wrap(err)

	httpsStatusCode := http.StatusBadRequest

	if errors.Is(err, sql.ErrNoRows) {
		httpsStatusCode = http.StatusNotFound
		appError.httpStatusCode = httpsStatusCode
	}

	return appError
}
