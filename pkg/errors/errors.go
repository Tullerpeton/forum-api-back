package errors

import (
	"fmt"
)

type Error struct {
	Message string `json:"error"`
}

func (err Error) Error() string {
	return fmt.Sprintf("error: happened %s", err.Message)
}

func CreateError(err error) error {
	if _, ok := err.(Error); ok {
		return err
	}

	return Error{Message: err.Error()}
}

var (
	ErrBadRequest error = Error{
		Message: "bad body of request",
	}
	ErrBadArguments error = Error{
		Message: "bad arguments of request",
	}
)