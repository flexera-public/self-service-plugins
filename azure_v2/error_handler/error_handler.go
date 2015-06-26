package errorHandler

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/labstack/echo"
	"net/http"
)

type genericError struct {
	echo.HTTPError
	StackTrace string `json:"StackTrace,omitempty"`
}

// AzureErrorHandler is a custom Echo.HTTPErrorHandler
func AzureErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c *echo.Context) {
		ge := new(genericError)
		ge.Code = http.StatusInternalServerError //default status code is 500
		ge.Message = http.StatusText(ge.Code)    // default message is 'Internal Server Error'
		if e.Debug() && ge.Code == 500 {
			ge.Message = err.Error() //show original error message in case of debug mode https://github.com/labstack/echo/blob/1e117621e9006481bfc0fd8e6bafab48c1848639/echo.go#L161
		}
		switch errorType := err.(type) {
		case *errors.Error:
			if he, ok := errorType.Err.(*echo.HTTPError); ok {
				ge.Code = he.Code
				ge.Message = he.Message
			}
			if e.Debug() && ge.Code == 500 {
				ge.StackTrace = errorType.ErrorStack()
			}
		case *echo.HTTPError:
			ge.Code = errorType.Code
			ge.Message = errorType.Message
		}

		c.JSON(ge.Code, ge)
	}
}

// GenericException represents error with status code 400
func GenericException(message string) error {
	return errors.New(&echo.HTTPError{
		Message: message,
		Code:    400,
	})
}

// RecordNotFound represents error with status code 404
func RecordNotFound(resourceID string) error {
	message := fmt.Sprintf("Could not find resource with id: %s", resourceID)
	return errors.New(&echo.HTTPError{
		Message: message,
		Code:    404,
	})
}
