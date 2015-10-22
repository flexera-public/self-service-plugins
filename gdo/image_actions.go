package main

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo"
	"github.com/rightscale/gdo/middleware"
	"github.com/rightscale/godo"
)

func SetupImageActionsRoutes(e *echo.Echo) {
	e.Get("/:actionId", getImageAction)
	e.Post("/transfer", transferImage)
	e.Post("/convert", convertImage)
}

func getImageAction(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	said := c.Param("actionId")
	if said == "" {
		return Error(fmt.Errorf("missing action id"))
	}
	aid, er := strconv.Atoi(said)
	if er != nil {
		return Error(fmt.Errorf("invalid action id '%s' - must be a number", said))
	}
	action, _, er := client.ImageActions.Get(id, aid)
	return Respond(c, action, er)
}

func transferImage(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		Region string `json:"region"`
	}{}
	if err = c.Bind(&req); err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	transferReq := godo.ActionRequest{"type": "transfer", "region": req.Region}
	action, _, er := client.ImageActions.Transfer(id, &transferReq)
	return Respond(c, action, er)
}

func convertImage(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	transferReq := godo.ActionRequest{"type": "convert"}
	action, _, er := client.ImageActions.Transfer(id, &transferReq)
	return Respond(c, action, er)
}
