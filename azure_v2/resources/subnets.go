package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type (
	subnetResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"` // required by response
	}

	subnetRequestParams struct {
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	subnetCreateParams struct {
		Name                   string `json:"name,omitempty"`
		Group                  string `json:"group_name,omitempty"`
		NetworkID              string `json:"network_id,omitempty"`
		AddressPrefix          string `json:"address_prefix,omitempty"`
		NetworkSecurityGroupID string `json:"network_security_group_id,omitempty"`
	}
	// Subnet is base struct for Azure Subnet resource to store input create params,
	// request create params and response params gotten from cloud.
	Subnet struct {
		createParams   subnetCreateParams
		requestParams  subnetRequestParams
		responseParams subnetResponseParams
	}
)

// SetupSubnetsRoutes declares routes for Subnet resource
func SetupSubnetsRoutes(e *echo.Group) {
	e.Get("/subnets", listAllSubnets)
	// e.Post("/subnets", createSubnet)
	// e.Delete("/subnets/:id", deleteSubnet)

	//nested routes
	group := e.Group("/resource_groups/:group_name/networks/:network_id/subnets")
	group.Get("", listSubnets)
	group.Get("/:id", listOneSubnet)
	group.Post("", createSubnet)
	group.Delete("/:id", deleteSubnet)
}

func listSubnets(c *echo.Context) error {
	groupName := c.Param("group_name")
	networkID := c.Param("network_id")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkPath, networkID, microsoftNetworkApiVersion)
	subnets, err := GetResources(c, path)
	if err != nil {
		return err
	}
	//add href for each subnet
	for _, subnet := range subnets {
		subnet["href"] = fmt.Sprintf("/resource_groups/%s/networks/%s/subnets/%s", groupName, networkID, subnet["name"])
	}
	return Render(c, 200, subnets, "vnd.rightscale.subnet+json;type=collection")
}

// To get all subnets faster could be used Network resource since each network contains set of subnets
func listAllSubnets(c *echo.Context) error {
	subnets := make([]map[string]interface{}, 0)
	path := fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, networkPath, microsoftNetworkApiVersion)
	networks, err := GetResources(c, path)
	if err != nil {
		return err
	}
	for _, network := range networks {
		array := strings.Split(network["id"].(string), "/")
		groupName := array[len(array)-5]
		networkID := network["name"].(string)
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkPath, networkID, microsoftNetworkApiVersion)
		resp, err := GetResources(c, path)
		if err != nil {
			return err
		}
		for _, subnet := range resp {
			subnet["href"] = fmt.Sprintf("/resource_groups/%s/networks/%s/subnets/%s", groupName, networkID, subnet["name"])
		}
		subnets = append(subnets, resp...)
	}
	return Render(c, 200, subnets, "vnd.rightscale.subnet+json;type=collection")
}

func listOneSubnet(c *echo.Context) error {
	subnet := Subnet{
		createParams: subnetCreateParams{
			Name:      c.Param("id"),
			Group:     c.Param("group_name"),
			NetworkID: c.Param("network_id"),
		},
	}
	return Get(c, &subnet)
}

func createSubnet(c *echo.Context) error {
	subnet := new(Subnet)
	return Create(c, subnet)
}

func deleteSubnet(c *echo.Context) error {
	subnet := Subnet{
		createParams: subnetCreateParams{
			Name:      c.Param("id"),
			Group:     c.Param("group_name"),
			NetworkID: c.Param("network_id"),
		},
	}
	return Delete(c, &subnet)
}

// GetRequestParams prepares parameters for create  request to the cloud
func (s *Subnet) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&s.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	s.createParams.Group = c.Param("group_name")
	s.createParams.NetworkID = c.Param("network_id")

	s.requestParams.Properties = map[string]interface{}{
		"addressPrefix": s.createParams.AddressPrefix,
	}

	if s.createParams.NetworkSecurityGroupID != "" {
		s.requestParams.Properties["networkSecurityGroup"] = map[string]interface{}{
			"id": s.createParams.NetworkSecurityGroupID,
		}
	}

	return s.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (s *Subnet) GetResponseParams() interface{} {
	return s.responseParams
}

// GetPath returns full path to the sigle subnet
func (s *Subnet) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, s.createParams.Group, networkPath, s.createParams.NetworkID, s.createParams.Name, microsoftNetworkApiVersion)
}

// GetCollectionPath is a fake function to support AzureResource by Subnet
func (s *Subnet) GetCollectionPath(groupName string) string { return "" }

// HandleResponse manage raw cloud response
func (s *Subnet) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &s.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := s.GetHref(s.responseParams.ID)
	if actionName == "create" {
		c.Response().Header().Add("Location", href)
	} else if actionName == "get" {
		s.responseParams.Href = href
	}
	return nil
}

// GetContentType returns subnet content type
func (s *Subnet) GetContentType() string {
	return "vnd.rightscale.subnet+json"
}

// GetHref returns subnet href
func (s *Subnet) GetHref(subnetID string) string {
	array := strings.Split(subnetID, "/")
	return fmt.Sprintf("resource_groups/%s/networks/%s/subnets/%s", array[len(array)-7], array[len(array)-3], array[len(array)-1])
}
