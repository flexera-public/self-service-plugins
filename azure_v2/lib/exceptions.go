package lib

import (
	"github.com/go-errors/errors"
	"github.com/labstack/echo"
)

func GenericException(message string) error {
	return errors.New(&echo.HTTPError{
		Message: message,
		Code:    400,
	})
}
