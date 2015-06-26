package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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
		Name       string                 `json:"name,omitempty"`
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	subnetCreateParams struct {
		Name          string `json:"name,omitempty"`
		Group         string `json:"group_name,omitempty"`
		NetworkID     string `json:"network_id,omitempty"`
		AddressPrefix string `json:"address_prefix,omitempty"`
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
func SetupSubnetsRoutes(e *echo.Echo) {
	e.Get("/subnets", listSubnets)
	e.Post("/subnets", createSubnet)
	e.Delete("/subnets/:id", deleteSubnet)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/networks/:network_id/subnets")
	// group.Get("", listSubnets)
	// group.Post("", createSubnet)
	// group.Delete("/:id", deleteSubnet)
}

func listSubnets(c *echo.Context) error {
	groupName := c.Param("group_name")
	if groupName != "" {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkPath, c.Param("network_id"), config.APIVersion)
		subnets, err := GetResources(c, path)
		if err != nil {
			return err
		}
		//add href for each subnet
		for _, subnet := range subnets {
			subnet["href"] = fmt.Sprintf("/resource_groups/%s/networks/%s/subnets/%s", groupName, c.Param("network_id"), subnet["name"])
		}
		return c.JSON(200, subnets)
	}

	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, "2015-01-01")
	resp, _ := GetResources(c, path)
	//TODO: add error handling
	var subnets []*subnetResponseParams
	for _, resourceGroup := range resp {
		groupName := resourceGroup["name"].(string)
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkPath, config.APIVersion)
		networks, _ := GetResources(c, path)
		//TODO: add error handling
		for _, network := range networks {
			network := network //.(map[string]interface{})
			resp, _ := getSubnets(c, groupName, network["name"].(string))
			//TODO: add error handling
			subnets = append(subnets, resp...)
		}
	}

	// init empty array
	if len(subnets) == 0 {
		subnets = make([]*subnetResponseParams, 0)
	}

	return c.JSON(200, subnets)
}

func createSubnet(c *echo.Context) error {
	subnet := new(Subnet)
	return Create(c, subnet)
}

func deleteSubnet(c *echo.Context) error {
	params := c.Request.Form
	subnet := Subnet{
		createParams: subnetCreateParams{
			Name:      c.Param("id"),
			Group:     params.Get("group_name"),
			NetworkID: params.Get("network_id"),
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

	s.requestParams.Properties = map[string]interface{}{
		"addressPrefix": s.createParams.AddressPrefix,
		//"dhcpOptions":   postParams.Get("dhcp_options"),
	}

	return s.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (s *Subnet) GetResponseParams() interface{} {
	return s.responseParams
}

// GetPath returns full path to the sigle subnet
func (s *Subnet) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, s.createParams.Group, networkPath, s.createParams.NetworkID, s.createParams.Name, config.APIVersion)
}

// GetCollectionPath is a fake function to support AzureResource by Subnet
func (s *Subnet) GetCollectionPath(groupName string) string { return "" }

// HandleResponse manage raw cloud response
func (s *Subnet) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &s.responseParams)
	href := s.GetHref(s.createParams.Group, s.responseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		s.responseParams.Href = href
	}
}

// GetContentType returns subnet content type
func (s *Subnet) GetContentType() string {
	return "vnd.rightscale.subnet+json"
}

// GetHref returns subnet href
func (s *Subnet) GetHref(groupName string, subnetName string) string {
	return fmt.Sprintf("/subnets/%s?group_name=%s&network=%s", subnetName, groupName, s.createParams.NetworkID)
}

//TODO: generify ListSubnets and getSubnets
func getSubnets(c *echo.Context, groupName string, networkName string) ([]*subnetResponseParams, error) {
	client, _ := GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkPath, networkName, config.APIVersion)
	log.Printf("Get Subents request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while getting subnet: %v", err))
	}
	defer resp.Body.Close()
	var m map[string][]*subnetResponseParams
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	subnets := m["value"]

	for _, subnet := range subnets {
		subnet.Href = fmt.Sprintf("/resource_groups/%s/networks/%s/subnets/%s", groupName, networkName, subnet.Name)
	}

	return subnets, nil
}
