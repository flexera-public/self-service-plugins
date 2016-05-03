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
	virtualNetworkGatewayPath = "providers/Microsoft.Network/virtualNetworkGateways"
)

type (
	virtualNetworkGatewayResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	virtualNetworkGatewayRequestParams struct {
		Name       string                 `json:"name,omitempty"`
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	virtualNetworkGatewayCreateParams struct {
		Name        string `json:"name,omitempty"`
		Location    string `json:"location,omitempty"`
		Group       string `json:"group_name,omitempty"`
		GatewayType string `json:"gateway_type,omitempty"`
		IPAddressId string `json:"ip_address_id,omitempty"`
		SubnetId    string `json:"subnet_id,omitempty"`
	}
	// VirtualNetworkGateway is base struct for VirtualNetworkGateway resource to store input create params,
	// request create params and response params gotten from cloud.
	VirtualNetworkGateway struct {
		createParams   virtualNetworkGatewayCreateParams
		requestParams  virtualNetworkGatewayRequestParams
		responseParams virtualNetworkGatewayResponseParams
	}
)

// SetupNetworkRoutes declares routes for VirtualNetworkGateway resource
func SetupVirtualNetworkGatewayRoutes(e *echo.Group) {
	e.Get("/virtual_network_gateways", listVirtualNetworkGateways)

	//nested routes
	group := e.Group("/resource_groups/:group_name/virtual_network_gateways")
	group.Get("", listVirtualNetworkGateways)
	group.Get("/:id", listOneVirtualNetworkGateway)
	group.Post("", createVirtualNetworkGateway)
	group.Delete("/:id", deleteVirtualNetworkGateway)
}

// List virtualNetworkGateways within a resource group or subscription
func listVirtualNetworkGateways(c *echo.Context) error {
	return List(c, new(VirtualNetworkGateway))
}

func listOneVirtualNetworkGateway(c *echo.Context) error {
	virtualNetworkGateway := VirtualNetworkGateway{
		createParams: virtualNetworkGatewayCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Get(c, &virtualNetworkGateway)
}

func createVirtualNetworkGateway(c *echo.Context) error {
	virtualNetworkGateway := new(VirtualNetworkGateway)
	return Create(c, virtualNetworkGateway)
}

func deleteVirtualNetworkGateway(c *echo.Context) error {
	virtualNetworkGateway := VirtualNetworkGateway{
		createParams: virtualNetworkGatewayCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Delete(c, &virtualNetworkGateway)
}

// GetRequestParams prepares parameters for create virtualNetworkGateway request to the cloud
func (vng *VirtualNetworkGateway) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&vng.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	vng.createParams.Group = c.Param("group_name")

	//TODO: add validation
	vng.requestParams.Name = vng.createParams.Name
	vng.requestParams.Location = vng.createParams.Location

	vng.requestParams.Properties = map[string]interface{}{
		"gatewayType": vng.createParams.GatewayType,
		"ipConfigurations": map[string]interface{}{
			"id": vng.createParams.IPAddressId,
		},
		"bgpEnabled": false, //TODO: investigate what it should be
		"subnet": map[string]interface{}{
			"id": vng.createParams.SubnetId,
		},
	}

	return vng.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (vng *VirtualNetworkGateway) GetResponseParams() interface{} {
	return vng.responseParams
}

// GetPath returns full path to the sigle virtualNetworkGateway
func (vng *VirtualNetworkGateway) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, vng.createParams.Group, virtualNetworkGatewayPath, vng.createParams.Name, microsoftNetworkApiVersion)
}

// GetCollectionPath returns full path to the collection of virtualNetworkGateway
func (vng *VirtualNetworkGateway) GetCollectionPath(groupName string) string {
	if groupName == "" {
		return fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, virtualNetworkGatewayPath, microsoftNetworkApiVersion)
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, virtualNetworkGatewayPath, microsoftNetworkApiVersion)
}

// HandleResponse manage raw cloud response
func (vng *VirtualNetworkGateway) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &vng.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := vng.GetHref(vng.responseParams.ID)
	if actionName == "create" {
		c.Response().Header().Add("Location", href)
	} else if actionName == "get" {
		vng.responseParams.Href = href
	}
	return nil
}

// GetContentType returns virtualNetworkGateway content type
func (vng *VirtualNetworkGateway) GetContentType() string {
	return "vnd.rightscale.virtual_network_gateway+json"
}

// GetHref returns virtualNetworkGateway href
func (vng *VirtualNetworkGateway) GetHref(virtualNetworkGatewayID string) string {
	array := strings.Split(virtualNetworkGatewayID, "/")
	return fmt.Sprintf("resource_groups/%s/virtualNetworkGateways/%s", array[len(array)-5], array[len(array)-1])
}
