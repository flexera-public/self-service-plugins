package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rightscale/godo"
)

// Respond with 200 and JSON if *echo.HTTPError is nil, handle *echo.HTTPError as described in RespondStatus otherwise.
func Respond(c *echo.Context, respObj interface{}, respErr error) *echo.HTTPError {
	return RespondStatus(c, 200, respObj, respErr)
}

// Respond with 204 if *echo.HTTPError is nil, handle *echo.HTTPError as described in RespondStatus otherwise.
func RespondNoContent(c *echo.Context, respErr error) *echo.HTTPError {
	return RespondStatus(c, 204, nil, respErr)
}

// Helper function that considers the given *echo.HTTPError and response object and send the appropriate
// HTTP response. The logic is as follows:
// 1. If the *echo.HTTPError is not nil go to 2 else go to 3
// 2. Is the *echo.HTTPError a 404? if yes respond with 404 else respond with 500 and the *echo.HTTPError message
// 3. Is the response object nil? if so respond with given status and empty body else respond with
//    given status and JSON encoded response object.
func RespondStatus(c *echo.Context, status int, respObj interface{}, respErr error) *echo.HTTPError {
	if respErr != nil {
		if godoErr, ok := respErr.(*godo.ErrorResponse); ok {
			if godoErr.Response != nil && godoErr.Response.StatusCode == 404 {
				return c.String(404, http.StatusText(404))
			}
		}
		return c.String(500, respErr.Error())
	}
	if respObj == nil {
		return c.NoContent(status)
	}
	return c.JSON(status, respObj)
}
