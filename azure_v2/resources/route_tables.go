package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	routeTablePath = "providers/Microsoft.Network/routeTables"
	apiVersion     = "2015-06-15"
)

type (
	routeTableResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	routeTableRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	routeTableCreateParams struct {
		Name     string   `json:"name,omitempty"`
		Location string   `json:"location,omitempty"`
		Group    string   `json:"group_name,omitempty"`
		Routes   []string `json:"routes,omitempty"`
	}
	// RouteTable is base struct for Azure Route Table resource to store input create params,
	// request create params and response params gotten from cloud.
	RouteTable struct {
		createParams   routeTableCreateParams
		requestParams  routeTableRequestParams
		responseParams routeTableResponseParams
	}
)

// SetupRouteTablesRoutes declares routes for RouteTable resource
func SetupRouteTablesRoutes(e *echo.Group) {
	e.Get("/route_tables", listRouteTables)

	//nested routes
	group := e.Group("/resource_groups/:group_name/route_tables")
	group.Get("", listRouteTables)
	group.Get("/:id", listOneRouteTable)
	group.Post("", createRouteTable)
	group.Delete("/:id", deleteRouteTable)
}

func listRouteTables(c *echo.Context) error {
	return List(c, new(RouteTable))
}

func listOneRouteTable(c *echo.Context) error {
	routeTable := RouteTable{
		createParams: routeTableCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Get(c, &routeTable)
}

func createRouteTable(c *echo.Context) error {
	routeTable := new(RouteTable)
	return Create(c, routeTable)
}

func deleteRouteTable(c *echo.Context) error {
	routeTable := RouteTable{
		createParams: routeTableCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Delete(c, &routeTable)
}

// GetRequestParams prepares parameters for create route table request to the cloud
func (rt *RouteTable) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&rt.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	rt.createParams.Group = c.Param("group_name")

	rt.requestParams.Location = rt.createParams.Location
	rt.requestParams.Properties = map[string]interface{}{
		"routes": rt.createParams.Routes,
	}

	return rt.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (rt *RouteTable) GetResponseParams() interface{} {
	return rt.responseParams
}

// GetPath returns full path to the sigle route table
func (rt *RouteTable) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, rt.createParams.Group, routeTablePath, rt.createParams.Name, microsoftNetworkApiVersion)
}

// GetCollectionPath returns full path to the collection of route tables
func (rt *RouteTable) GetCollectionPath(groupName string) string {
	if groupName == "" {
		return fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, routeTablePath, microsoftNetworkApiVersion)
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, routeTablePath, microsoftNetworkApiVersion)
}

// HandleResponse manage raw cloud response
func (rt *RouteTable) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &rt.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := rt.GetHref(rt.responseParams.ID)
	if actionName == "create" {
		c.Response().Header().Add("Location", href)
	} else if actionName == "get" {
		rt.responseParams.Href = href
	}
	return nil
}

// GetContentType returns route table content type
func (rt *RouteTable) GetContentType() string {
	return "vnd.rightscale.route_table+json"
}

// GetHref returns route table href
func (rt *RouteTable) GetHref(routeTableID string) string {
	array := strings.Split(routeTableID, "/")
	return fmt.Sprintf("resource_groups/%s/route_tables/%s", array[len(array)-5], array[len(array)-1])
}
