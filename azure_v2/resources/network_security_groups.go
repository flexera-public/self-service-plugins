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
	networkSecurityGroupPath = "providers/Microsoft.Network/networkSecurityGroups"
)

type (
	networkSecurityGroupResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	networkSecurityGroupRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	networkSecurityGroupCreateParams struct {
		Name          string                   `json:"name,omitempty"`
		Location      string                   `json:"location,omitempty"`
		Group         string                   `json:"group_name,omitempty"`
		SecurityRules []map[string]interface{} `json:"security_rules,omitempty"`
	}
	// NetworkSecurityGroup is base struct for Azure Network Security Group resource to store input create params,
	// request create params and response params gotten from cloud.
	NetworkSecurityGroup struct {
		createParams   networkSecurityGroupCreateParams
		requestParams  networkSecurityGroupRequestParams
		responseParams networkSecurityGroupResponseParams
	}
)

// SetupNetworkSecurityGroupRoutes declares routes for NetworkSecurityGroup resource
func SetupNetworkSecurityGroupRoutes(e *echo.Echo) {
	e.Get("/network_security_groups", listNetworkSecurityGroup)

	//nested routes
	group := e.Group("/resource_groups/:group_name/network_security_groups")
	group.Get("", listNetworkSecurityGroup)
	group.Get("/:id", listOneNetworkSecurityGroup)
	group.Post("", createNetworkSecurityGroup)
	group.Delete("/:id", deleteNetworkSecurityGroup)
}

func listNetworkSecurityGroup(c *echo.Context) error {
	return List(c, new(NetworkSecurityGroup))
}

func listOneNetworkSecurityGroup(c *echo.Context) error {
	networkSecurityGroup := NetworkSecurityGroup{
		createParams: networkSecurityGroupCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Get(c, &networkSecurityGroup)
}

func createNetworkSecurityGroup(c *echo.Context) error {
	networkSecurityGroup := new(NetworkSecurityGroup)
	return Create(c, networkSecurityGroup)
}

func deleteNetworkSecurityGroup(c *echo.Context) error {
	networkSecurityGroup := NetworkSecurityGroup{
		createParams: networkSecurityGroupCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Delete(c, &networkSecurityGroup)
}

// GetRequestParams prepares parameters for create network security group request to the cloud
func (nsg *NetworkSecurityGroup) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&nsg.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	nsg.createParams.Group = c.Param("group_name")

	nsg.requestParams.Location = nsg.createParams.Location
	nsg.requestParams.Properties = map[string]interface{}{
		"securityRules": nsg.createParams.SecurityRules,
	}

	return nsg.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (nsg *NetworkSecurityGroup) GetResponseParams() interface{} {
	return nsg.responseParams
}

// GetPath returns full path to the sigle network security group
func (nsg *NetworkSecurityGroup) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, nsg.createParams.Group, networkSecurityGroupPath, nsg.createParams.Name, config.APIVersion)
}

// GetCollectionPath returns full path to the collection of network security groups
func (nsg *NetworkSecurityGroup) GetCollectionPath(groupName string) string {
	if groupName == "" {
		return fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, networkSecurityGroupPath, config.APIVersion)
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkSecurityGroupPath, config.APIVersion)
}

// HandleResponse manage raw cloud response
func (nsg *NetworkSecurityGroup) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &nsg.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := nsg.GetHref(nsg.responseParams.ID)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		nsg.responseParams.Href = href
	}
	return nil
}

// GetContentType returns network security group content type
func (nsg *NetworkSecurityGroup) GetContentType() string {
	return "vnd.rightscale.network_security_group+json"
}

// GetHref returns network security group href
func (nsg *NetworkSecurityGroup) GetHref(networkSecurityGroupID string) string {
	array := strings.Split(networkSecurityGroupID, "/")
	return fmt.Sprintf("/resource_groups/%s/network_security_groups/%s", array[len(array)-5], array[len(array)-1])
}
