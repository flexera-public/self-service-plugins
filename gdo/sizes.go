package main

import (
	"github.com/labstack/echo"
	"github.com/rightscale/gdo/middleware"
	"github.com/rightscale/godo"
)

func SetupSizesRoutes(e *echo.Echo) {
	e.Get("", listSizes)
}

func listSizes(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	list, er := paginateSizes(client.Sizes.List)
	return Respond(c, list, er)
}

func paginateSizes(lister func(opt *godo.ListOptions) ([]godo.Size, *godo.Response, error)) ([]godo.Size, error) {
	list := []godo.Size{}
	opt := &godo.ListOptions{}
	for {
		sizes, resp, err := lister(opt)
		if err != nil {
			return nil, err
		}
		list = append(list, sizes...)
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
