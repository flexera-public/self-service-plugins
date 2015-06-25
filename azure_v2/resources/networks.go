package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	networkPath = "providers/Microsoft.Network/virtualNetworks"
)

type (
	NetworkResponseParams struct {
		Id         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	NetworkRequestParams struct {
		Name       string                 `json:"name,omitempty"`
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	NetworkCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	Network struct {
		CreateParams   NetworkCreateParams
		RequestParams  NetworkRequestParams
		ResponseParams NetworkResponseParams
	}
)

func SetupNetworkRoutes(e *echo.Echo) {
	e.Get("/networks", listNetworks)
	e.Get("/networks/:id", listOneNetwork)
	e.Post("/networks", createNetwork)
	e.Delete("/networks/:id", deleteNetwork)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/networks")
	// group.Get("", listNetworks)
	// group.Post("", createNetwork)
	// group.Delete("/:id", deleteNetwork)
}

func listNetworks(c *echo.Context) error {
	return List(c, new(Network))
}

func listOneNetwork(c *echo.Context) error {
	params := c.Request.Form
	network := Network{
		CreateParams: NetworkCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Get(c, &network)
}

func createNetwork(c *echo.Context) error {
	network := new(Network)
	return Create(c, network)
}

func deleteNetwork(c *echo.Context) error {
	params := c.Request.Form
	network := Network{
		CreateParams: NetworkCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Delete(c, &network)
}

func (n *Network) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&n.CreateParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	var subnets []map[string]interface{}
	n.RequestParams.Name = n.CreateParams.Name
	n.RequestParams.Location = n.CreateParams.Location
	n.RequestParams.Properties = map[string]interface{}{
		"addressSpace": map[string]interface{}{
			"addressPrefixes": []string{"10.0.0.0/16"},
		},
		"subnets": append(subnets, map[string]interface{}{
			"name": n.CreateParams.Name,
			"properties": map[string]interface{}{
				"addressPrefix": "10.0.0.0/16",
			},
		}),
	}

	return n.RequestParams, nil
}

func (n *Network) GetResponseParams() interface{} {
	return n.ResponseParams
}

func (n *Network) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, n.CreateParams.Group, networkPath, n.CreateParams.Name, config.ApiVersion)
}

func (n *Network) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, networkPath, config.ApiVersion)
}

func (n *Network) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &n.ResponseParams)
	href := n.GetHref(n.CreateParams.Group, n.ResponseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		n.ResponseParams.Href = href
	}
}

func (n *Network) GetContentType() string {
	return "vnd.rightscale.network+json"
}

func (n *Network) GetHref(groupName string, networkName string) string {
	return fmt.Sprintf("/networks/%s?group_name=%s", networkName, groupName)
}
