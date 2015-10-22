package main

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo"
	"github.com/rightscale/gdo/middleware"
	"github.com/rightscale/godo"
)

// Setup routes for droplet actions
func SetupDropletsRoutes(e *echo.Echo) {
	e.Get("", listDroplets)
	e.Get("/:id", showDroplet)
	e.Post("", createDroplet)
	e.Delete("/:id", deleteDroplet)
	e.Get("/:id/kernels", listDropletKernels)
	e.Get("/:id/snapshots", listDropletSnapshots)
	e.Get("/:id/backups", listDropletBackups)
	e.Get("/:id/actions", listDropletActions)
	e.Get("/:id/neighbors", listDropletNeighbors)
}

// Helper function that builds a droplet resource href from its id
func dropletHref(id int) string {
	return fmt.Sprintf("/droplets/%d", id)
}

func listDroplets(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	list, er := paginateDroplets(client.Droplets.List)
	return Respond(c, list, er)
}

func showDroplet(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	root, _, er := client.Droplets.Get(id)
	var droplet *Droplet
	if err == nil {
		droplet = DropletFromApi(root.Droplet)
	}
	return Respond(c, droplet, er)
}

func createDroplet(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := godo.DropletCreateRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}
	root, _, er := client.Droplets.Create(&req)
	if er == nil {
		c.Response.Header().Set("Location", dropletHref(root.Droplet.ID))
	}
	return RespondNoContent(c, er)
}

func deleteDroplet(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	_, er := client.Droplets.Delete(id)
	return RespondNoContent(c, er)
}

func listDropletKernels(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	list := []godo.Kernel{}
	opt := &godo.ListOptions{}
	for {
		kernels, resp, err := client.Droplets.Kernels(id, opt)
		if err != nil {
			return Error(err)
		}
		list = append(list, kernels...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return Error(err)
		}
		opt.Page = page + 1
	}
	return Respond(c, list, nil)
}

func listDropletSnapshots(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return listDropletImages(c, client.Droplets.Snapshots)
}

func listDropletBackups(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	return listDropletImages(c, client.Droplets.Backups)
}

func listDropletImages(c *echo.Context, lister func(int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)) *echo.HTTPError {
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	list, er := paginateImages(func(opt *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
		return lister(id, opt)
	})
	return Respond(c, list, er)
}

func listDropletActions(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	list, er := paginateActions(func(opt *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
		return client.Droplets.Actions(id, opt)
	})
	return Respond(c, list, er)
}

func listDropletNeighbors(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	droplets, _, er := client.Droplets.Neighbors(id)
	var list []*Droplet
	if er == nil {
		for _, d := range droplets {
			list = append(list, DropletFromApi(&d))
		}
	}
	return Respond(c, list, er)
}

// Helper function that retrieves the droplet id (number) from the request parameters
func getIDParam(c *echo.Context) (int, *echo.HTTPError) {
	sid := c.Param("id")
	if sid == "" {
		return 0, Error(fmt.Errorf("missing droplet id"))
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		return 0, Error(fmt.Errorf("invalid droplet id '%s' - must be a number", sid))
	}
	return id, nil
}

// Paginate over droplet listing
func paginateDroplets(lister func(opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)) ([]godo.Droplet, error) {
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := lister(opt)
		if err != nil {
			return nil, err
		}
		list = append(list, droplets...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opt.Page = page + 1
	}
	return list, nil
}
