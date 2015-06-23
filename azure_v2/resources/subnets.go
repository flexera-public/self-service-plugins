package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

type (
	SubnetResponseParams struct {
		Id         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"` // required by response
	}

	SubnetRequestParams struct {
		Name       string                 `json:"name,omitempty"`
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	SubnetCreateParams struct {
		Name          string `json:"name,omitempty"`
		Group         string `json:"group_name,omitempty"`
		NetworkId     string `json:"network_id,omitempty"`
		AddressPrefix string `json:"address_prefix,omitempty"`
	}
	Subnet struct {
		CreateParams   SubnetCreateParams
		RequestParams  SubnetRequestParams
		ResponseParams SubnetResponseParams
	}
)

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
	group_name := c.Param("group_name")
	if group_name != "" {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, c.Param("network_id"), config.ApiVersion)
		subnets, err := lib.GetResources(c, path)
		if err != nil {
			return err
		}
		//add href for each subnet
		for _, subnet := range subnets {
			subnet["href"] = fmt.Sprintf("/resource_groups/%s/networks/%s/subnets/%s", group_name, c.Param("network_id"), subnet["name"])
		}
		return c.JSON(200, subnets)
	} else {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
		resp, _ := lib.GetResources(c, path)
		//TODO: add error handling
		var subnets []*SubnetResponseParams
		for _, resource_group := range resp {
			groupName := resource_group["name"].(string)
			path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, networkPath, config.ApiVersion)
			networks, _ := lib.GetResources(c, path)
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
			subnets = make([]*SubnetResponseParams, 0)
		}
		return c.JSON(200, subnets)
	}

}

func createSubnet(c *echo.Context) error {
	subnet := new(Subnet)
	return lib.Create(c, subnet)
}

func deleteSubnet(c *echo.Context) error {
	params := c.Request.Form
	subnet := Subnet{
		CreateParams: SubnetCreateParams{
			Name:      c.Param("id"),
			Group:     params.Get("group_name"),
			NetworkId: params.Get("network_id"),
		},
	}
	return lib.Delete(c, &subnet)
}

func (s *Subnet) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&s.CreateParams)
	if err != nil {
		return nil, lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	s.RequestParams.Properties = map[string]interface{}{
		"addressPrefix": s.CreateParams.AddressPrefix,
		//"dhcpOptions":   postParams.Get("dhcp_options"),
	}

	return s.RequestParams, nil
}

func (s *Subnet) GetResponseParams() interface{} {
	return s.ResponseParams
}

func (s *Subnet) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, s.CreateParams.Group, networkPath, s.CreateParams.NetworkId, s.CreateParams.Name, config.ApiVersion)
}

//fake function to support AzureResource by Subnet
func (s *Subnet) GetCollectionPath(groupName string) string {
	return ""
}

func (s *Subnet) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &s.ResponseParams)
	href := s.GetHref(s.CreateParams.Group, s.ResponseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		s.ResponseParams.Href = href
	}
}

func (s *Subnet) GetContentType() string {
	return "vnd.rightscale.subnet+json"
}

func (s *Subnet) GetHref(groupName string, subnetName string) string {
	return fmt.Sprintf("/subnets/%s?group_name=%s&network=%s", subnetName, groupName, s.CreateParams.NetworkId)
}

//TODO: generify ListSubnets and getSubnets
func getSubnets(c *echo.Context, group_name string, network_name string) ([]*SubnetResponseParams, error) {
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, network_name, config.ApiVersion)
	log.Printf("Get Subents request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return nil, lib.GenericException(fmt.Sprintf("Error has occurred while getting subnet: %v", err))
	}
	defer resp.Body.Close()
	var m map[string][]*SubnetResponseParams
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	subnets := m["value"]

	for _, subnet := range subnets {
		subnet.Href = fmt.Sprintf("/resource_groups/%s/networks/%s/subnets/%s", group_name, network_name, subnet.Name)
	}

	return subnets, nil
}
