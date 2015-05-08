package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/gdo/middleware"
	"github.com/rightscale/godo"
)

func SetupImagesRoutes(e *echo.Echo) {
	e.Get("", listImages)
	e.Get("/:id", showImage)
}

// Helper function that builds an image resource href from its id
func imageHref(id int) string {
	return fmt.Sprintf("/images/%d", id)
}

func listImages(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	list, er := paginateImages(client.Images.List)
	return Respond(c, list, er)
}

func showImage(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	image, _, er := client.Images.GetByID(id)
	return Respond(c, image, er)
}

func paginateImages(lister func(opt *godo.ListOptions) ([]godo.Image, *godo.Response, error)) ([]godo.Image, error) {
	list := []godo.Image{}
	opt := &godo.ListOptions{}
	for {
		images, resp, err := lister(opt)
		if err != nil {
			return nil, err
		}
		list = append(list, images...)
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
