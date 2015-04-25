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

func getImageAction(c *echo.Context) error {
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
		return fmt.Errorf("missing action id")
	}
	aid, err := strconv.Atoi(said)
	if err != nil {
		return fmt.Errorf("invalid action id '%s' - must be a number", said)
	}
	action, _, err := client.ImageActions.Get(id, aid)
	if err != nil {
		return err
	}
	return Respond(c, action)
}

func transferImage(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		Region string `json:"region"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	transferReq := godo.ActionRequest{"type": "transfer", "region": req.Region}
	action, _, err := client.ImageActions.Transfer(id, &transferReq)
	if err != nil {
		return err
	}
	return Respond(c, action)
}

func convertImage(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	transferReq := godo.ActionRequest{"type": "convert"}
	action, _, err := client.ImageActions.Transfer(id, &transferReq)
	if err != nil {
		return err
	}
	return Respond(c, action)
}
