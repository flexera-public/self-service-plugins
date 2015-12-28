package resources

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type (
	routesResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	routesRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	routesCreateParams struct {
		Name             string `json:"name"`
		Location         string `json:"location"`
		Prefix           string `json:"address_prefix"`
		NextHopType      string `json:"next_hop_type"`
		NextHopIpAddress string `json:"next_hop_ip_address,omitempty"`
		// the following params are required for building path
		Group          string `json:"group_name,omitempty"`
		RouteTableName string `json:"route_table_name,omitempty"`
	}
	// Route is base struct for Azure Route resource to store input create params,
	// request create params and response params gotten from cloud.
	Route struct {
		createParams   routesCreateParams
		requestParams  routesRequestParams
		responseParams routesResponseParams
	}
)

// SetupRoutes declares routes for Route resource
func SetupRoutes(e *echo.Group) {
	e.Get("/routes", listAllRoutes)

	//nested routes
	group := e.Group("/resource_groups/:group_name/route_tables/:route_table_name/routes")
	group.Get("", listRoutes)
	group.Get("/:id", listOneRoute)
	group.Post("", createRoute)
	group.Delete("/:id", deleteRoute)
}

func listAllRoutes(c *echo.Context) error {
	routes := make([]map[string]interface{}, 0)
	path := fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, routeTablePath, apiVersion)
	tables, err := GetResources(c, path)
	if err != nil {
		return err
	}
	for _, table := range tables {
		array := strings.Split(table["id"].(string), "/")
		groupName := array[len(array)-5]
		tableID := table["name"].(string)
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/routes?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, routeTablePath, tableID, apiVersion)
		resp, err := GetResources(c, path)
		if err != nil {
			return err
		}
		for _, route := range resp {
			route["href"] = fmt.Sprintf("/resource_groups/%s/route_tables/%s/routes/%s", groupName, tableID, route["name"])
			route["location"] = table["location"]
		}
		routes = append(routes, resp...)
	}
	return Render(c, 200, routes, "vnd.rightscale.route+json;type=collection")
}

// it doesn't return 'location' as listRoutes or listAllRoutes
func listRoutes(c *echo.Context) error {
	groupName := c.Param("group_name")
	tableID := c.Param("route_table_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/routes?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, routeTablePath, tableID, apiVersion)
	routes, err := GetResources(c, path)
	if err != nil {
		return err
	}
	//add href for each rule
	for _, route := range routes {
		route["href"] = fmt.Sprintf("/resource_groups/%s/route_tables/%s/routes/%s", groupName, tableID, route["name"])
	}
	return Render(c, 200, routes, "vnd.rightscale.routes+json;type=collection")
}

// it doesn't return 'location' as listRoutes or listAllRoutes
func listOneRoute(c *echo.Context) error {
	routes := Route{
		createParams: routesCreateParams{
			Name:           c.Param("id"),
			Group:          c.Param("group_name"),
			RouteTableName: c.Param("route_table_name"),
		},
	}
	return Get(c, &routes)
}

func createRoute(c *echo.Context) error {
	routes := new(Route)
	return Create(c, routes)
}

func deleteRoute(c *echo.Context) error {
	routes := Route{
		createParams: routesCreateParams{
			Name:           c.Param("id"),
			Group:          c.Param("group_name"),
			RouteTableName: c.Param("route_table_name"),
		},
	}
	return Delete(c, &routes)
}

// GetRequestParams prepares parameters for create route request to the cloud
func (r *Route) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&r.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	r.createParams.Group = c.Param("group_name")
	r.createParams.RouteTableName = c.Param("route_table_name")

	r.requestParams.Location = r.createParams.Location
	r.requestParams.Properties = map[string]interface{}{
		"addressPrefix":    r.createParams.Prefix,
		"nextHopType":      r.createParams.NextHopType,
		"nextHopIpAddress": r.createParams.NextHopIpAddress,
	}

	return r.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (r *Route) GetResponseParams() interface{} {
	return r.responseParams
}

// GetPath returns full path to the sigle route
func (r *Route) GetPath() string {
	rr := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/routes/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, r.createParams.Group, routeTablePath, r.createParams.RouteTableName, r.createParams.Name, apiVersion)
	log.Printf("Path: %s\n", rr)
	return rr
}

// GetCollectionPath returns full path to the collection of routes
func (r *Route) GetCollectionPath(groupName string) string { return "" }

// HandleResponse manage raw cloud response
func (r *Route) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &r.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := r.GetHref(r.responseParams.ID)
	if actionName == "create" {
		c.Response().Header().Add("Location", href)
	} else if actionName == "get" {
		r.responseParams.Href = href
	}
	return nil
}

// GetContentType returns route content type
func (r *Route) GetContentType() string {
	return "vnd.rightscale.route+json"
}

// GetHref returns route href
func (r *Route) GetHref(routesID string) string {
	array := strings.Split(routesID, "/")
	return fmt.Sprintf("resource_groups/%s/route_tables/%s/routes/%s", array[len(array)-7], array[len(array)-3], array[len(array)-1])
}
