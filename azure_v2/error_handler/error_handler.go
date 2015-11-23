package errorHandler

import (
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/labstack/echo"
)

type genericError struct {
	Code       int
	Message    string
	StackTrace string `json:"StackTrace,omitempty"`
}

func (e *genericError) Error() string {
	return e.Message
}

// AzureErrorHandler is a custom Echo.HTTPErrorHandler
func AzureErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c *echo.Context) {
		ge := new(genericError)
		ge.Code = http.StatusInternalServerError // default status code is 500
		ge.Message = http.StatusText(ge.Code)    // default message is 'Internal Server Error'
		if e.Debug() {
			ge.Message = err.Error() //show original error message in case of debug mode https://github.com/labstack/echo/blob/1e117621e9006481bfc0fd8e6bafab48c1848639/echo.go#L161
		}
		switch errorType := err.(type) {
		case *errors.Error:
			if he, ok := errorType.Err.(*genericError); ok {
				ge = he
			}
			if e.Debug() && ge.Code == 500 {
				ge.StackTrace = errorType.ErrorStack()
			}
		case *echo.HTTPError:
			ge.Code = errorType.Code()
			ge.Message = errorType.Error()
		}

		c.JSON(ge.Code, ge)
	}
}

// GenericException represents error with status code 400
func GenericException(message string) error {
	return errors.New(&genericError{
		Code:    400,
		Message: message,
	})
}

// RecordNotFound represents error with status code 404
func RecordNotFound(resourceID string) error {
	message := fmt.Sprintf("Could not find resource with id: %s", resourceID)
	return errors.New(&genericError{
		Code:    404,
		Message: message,
	})
}

// InvalidParamException returns generic error massage for invalid value of paramName
func InvalidParamException(paramName string) error {
	message := fmt.Sprintf("You have specified an invalid '%s' parameter.", paramName)
	return errors.New(&genericError{
		Code:    400,
		Message: message,
	})
}
