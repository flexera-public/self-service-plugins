package lib

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/labstack/echo"
)

func GenericException(message string) error {
	return errors.New(&echo.HTTPError{
		Message: message,
		Code:    400,
	})
}

func RecordNotFound(resourceID string) error {
	message := fmt.Sprintf("Could not find resource with id: %s", resourceID)
	return errors.New(&echo.HTTPError{
		Message: message,
		Code:    404,
	})
}
