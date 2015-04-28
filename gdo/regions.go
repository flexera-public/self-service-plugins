package main

import (
	"github.com/labstack/echo"
	"github.com/rightscale/gdo/middleware"
	"github.com/rightscale/godo"
)

func SetupRegionsRoutes(e *echo.Echo) {
	e.Get("", listRegions)
}

func listRegions(c *echo.Context) error {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	list, err := paginateRegions(client.Regions.List)
	return Respond(c, list, err)
}

func paginateRegions(lister func(opt *godo.ListOptions) ([]godo.Region, *godo.Response, error)) ([]godo.Region, error) {
	list := []godo.Region{}
	opt := &godo.ListOptions{}
	for {
		regions, resp, err := lister(opt)
		if err != nil {
			return nil, err
		}
		list = append(list, regions...)
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
