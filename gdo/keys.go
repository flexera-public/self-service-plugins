package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/gdo/middleware"
	"github.com/rightscale/godo"
)

func SetupKeysRoutes(e *echo.Echo) {
	e.Get("", listKeys)
	e.Get("/:id", showKey)
	e.Post("", createKey)
	e.Put("/:id", updateKey)
	e.Delete("/:id", deleteKey)
}

func listKeys(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	list, er := paginateKeys(client.Keys.List)
	return Respond(c, list, er)
}

func showKey(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	key, _, er := client.Keys.GetByID(id)
	if er != nil {
		return Error(er)
	}
	return Respond(c, key, er)
}

func createKey(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := godo.KeyCreateRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}
	key, _, er := client.Keys.Create(&req)
	if er == nil {
		c.Response.Header().Set("Location", keyHref(key.ID))
	}
	return RespondNoContent(c, er)
}

func updateKey(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	req := godo.KeyUpdateRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	_, _, er := client.Keys.UpdateByID(id, &req)
	return RespondNoContent(c, er)
}

func deleteKey(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	_, er := client.Keys.DeleteByID(id)
	return RespondNoContent(c, er)
}

func keyHref(keyID int) string {
	return fmt.Sprintf("/v2/account/keys/%d", keyID)
}

func paginateKeys(lister func(opt *godo.ListOptions) ([]godo.Key, *godo.Response, error)) ([]godo.Key, error) {
	list := []godo.Key{}
	opt := &godo.ListOptions{}
	for {
		keys, resp, err := lister(opt)
		if err != nil {
			return nil, err
		}
		list = append(list, keys...)
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
