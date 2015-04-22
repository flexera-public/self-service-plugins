package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/godo"
)

func SetupKeysRoutes(e *echo.Echo) {
	e.Get("", listKeys)
	e.Get("/:id", showKey)
	e.Post("", createKey)
	e.Put("/:id", updateKey)
	e.Delete("/:id", deleteKey)
}

func listKeys(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	list, err := paginateKeys(client.Keys.List)
	if err != nil {
		return err
	}
	Respond(c, list)
	return nil
}

func showKey(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	key, _, err := client.Keys.GetByID(id)
	if err != nil {
		return err
	}
	Respond(c, key)
	return nil
}

func createKey(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	req := godo.KeyCreateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	key, _, err := client.Keys.Create(&req)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Location", keyHref(key.ID))
	return RespondNoContent(c)
}

func updateKey(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	req := godo.KeyUpdateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	_, _, err = client.Keys.UpdateByID(id, &req)
	if err != nil {
		return err
	}
	return RespondNoContent(c)
}

func deleteKey(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	_, err = client.Keys.DeleteByID(id)
	if err != nil {
		return err
	}
	return RespondNoContent(c)
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
