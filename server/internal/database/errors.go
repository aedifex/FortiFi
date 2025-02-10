package database

import (
	"errors"
	"fmt"
	"net/http"
)

type DatabaseError struct {
	Err        error
	HttpStatus int
}

var (
	USER_EXISTS_ERROR = &DatabaseError {
		Err:        errors.New("user exists"),
		HttpStatus: http.StatusConflict,
	}

	PREPARE_ERROR = func(err error) *DatabaseError {
		return &DatabaseError{
			Err:        fmt.Errorf("error preparing statement: %s", err),
			HttpStatus: http.StatusInternalServerError,
		}
	}

	QUERY_ERROR = func(err error) *DatabaseError {
		return &DatabaseError{
			Err:        fmt.Errorf("error querying database: %s", err),
			HttpStatus: http.StatusInternalServerError,
		}
	}

	EXEC_ERROR = func(err error) *DatabaseError {
		return &DatabaseError{
			Err:        fmt.Errorf("error executing statement: %s", err),
			HttpStatus: http.StatusInternalServerError,
		}
	}

	ROWS_AFFECTED_ERROR = func(err error) *DatabaseError {
		return &DatabaseError{
			Err:        fmt.Errorf("error fetching rows affected: %s", err),
			HttpStatus: http.StatusInternalServerError,
		}
	}

	DNE_ERROR = &DatabaseError{
		Err:        errors.New("resource does not exist"),
		HttpStatus: http.StatusNotFound,
	}

	INVALID_INPUT_ERROR = &DatabaseError{
		Err:        errors.New("invalid input"),
		HttpStatus: http.StatusBadRequest,
	}

	SCAN_ERROR = func(err error) *DatabaseError {
		return &DatabaseError{
			Err:        fmt.Errorf("error scanning table: %s", err),
			HttpStatus: http.StatusInternalServerError,
		}
	}

	UNAUTHORIZED_ERROR = &DatabaseError{
		Err:        errors.New("unauthorized"),
		HttpStatus: http.StatusUnauthorized,
	}

	HASH_ERROR = func(err error) *DatabaseError {
		return &DatabaseError{
			Err:        err,
			HttpStatus: http.StatusInternalServerError,
		}
	}
)
