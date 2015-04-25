package main

import "github.com/labstack/echo"

func Respond(c *echo.Context, resp interface{}) error {
	return c.JSON(200, resp)
}

func RespondError(c *echo.Context, resp error) error {
	return c.String(500, resp.Error())
}

func RespondNoContent(c *echo.Context) error {
	return c.NoContent(204)
}
