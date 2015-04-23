package main

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo"
	"github.com/rightscale/go-digital-ocean/middleware"
	"github.com/rightscale/godo"
)

func SetupDropletActionsRoutes(e *echo.Echo) {
	e.Get("/:actionId", getDropletAction)
	e.Post("/disable_backups", disableDropletBackups)
	e.Post("/reboot", rebootDroplet)
	e.Post("/power", powerCycleDroplet)
	e.Post("/shutdown", shutdownDroplet)
	e.Post("/power_off", powerOffDroplet)
	e.Post("/power_on", powerOnDroplet)
	e.Post("/restore", restoreDroplet)
	e.Post("/password_reset", passwordResetDroplet)
	e.Post("/resize", resizeDroplet)
	e.Post("/rebuild", rebuildDroplet)
	e.Post("/rename", renameDroplet)
	e.Post("/change_kernel", changeDropletKernel)
	e.Post("/enable_ipv6", enableDropletIPv6)
	e.Post("/enable_private_networking", enableDropletPrivateNetworking)
	e.Post("/snapshot", snapshotDroplet)
	e.Post("/upgrade", upgradeDroplet)
}

func getDropletAction(c *echo.Context) error {
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
	action, _, err := client.DropletActions.Get(id, aid)
	if err != nil {
		return err
	}
	return Respond(c, action)
}

func disableDropletBackups(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.DisableBackups)
}

func rebootDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.Reboot)
}
func powerCycleDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.PowerCycle)
}

func shutdownDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.Shutdown)
}

func powerOffDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.PowerOff)
}

func powerOnDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.PowerOn)
}

func restoreDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		ImageID int `json:"imageID"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	return doResourceAction(c, func(id int) (*godo.Action, *godo.Response, error) {
		return client.DropletActions.Restore(id, req.ImageID)
	})
}

func passwordResetDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.PasswordReset)
}

func resizeDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		SizeSlug   string `json:"sizeSlug"`
		ResizeDisk bool   `json:"resizeDisk"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	return doResourceAction(c, func(id int) (*godo.Action, *godo.Response, error) {
		return client.DropletActions.Resize(id, req.SizeSlug, req.ResizeDisk)
	})
}

func rebuildDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		ImageID int `json:"imageID"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	return doResourceAction(c, func(id int) (*godo.Action, *godo.Response, error) {
		return client.DropletActions.RebuildByImageID(id, req.ImageID)
	})
}

func renameDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		Name string `json:"name"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	return doResourceAction(c, func(id int) (*godo.Action, *godo.Response, error) {
		return client.DropletActions.Rename(id, req.Name)
	})
}

func changeDropletKernel(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		KernelID int `json:"kernelID"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	return doResourceAction(c, func(id int) (*godo.Action, *godo.Response, error) {
		return client.DropletActions.ChangeKernel(id, req.KernelID)
	})
}

func enableDropletIPv6(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.EnableIPv6)
}

func enableDropletPrivateNetworking(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.EnablePrivateNetworking)
}

func snapshotDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := struct {
		Name string `json:"name"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	return doResourceAction(c, func(id int) (*godo.Action, *godo.Response, error) {
		return client.DropletActions.Snapshot(id, req.Name)
	})
}

func upgradeDroplet(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return doResourceAction(c, client.DropletActions.Upgrade)
}

// Helper function that calls given client function and builds response accordingly
func doResourceAction(c *echo.Context, actionFunc func(int) (*godo.Action, *godo.Response, error)) error {
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	action, _, err := actionFunc(id)
	if err != nil {
		return err
	}
	return Respond(c, action)
}
