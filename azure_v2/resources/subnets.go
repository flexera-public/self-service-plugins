package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

type Subnet struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
	Href       string      `json:"href,omitempty"` // required by response
}

func SetupSubnetsRoutes(e *echo.Echo) {
	e.Get("/subnets", listSubnets)

	//nested routes
	group := e.Group("/resource_groups/:group_name/networks/:network_id/subnets")
	group.Get("", listSubnets)
	group.Post("", createSubnet)
	group.Delete("/:id", deleteSubnet)
}

func listSubnets(c *echo.Context) error {
	group_name := c.Param("group_name")
	if group_name != "" {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, c.Param("network_id"), config.ApiVersion)
		subnets, err := lib.GetResources(c, path, "/subnets/%s?group_name="+group_name)
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
		resp, _ := lib.GetResources(c, path, "/resource_group/%s")
		//TODO: add error handling
		var subnets []*Subnet
		for _, resource_group := range resp {
			groupName := resource_group["name"].(string)
			path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, networkPath, config.ApiVersion)
			networks, _ := lib.GetResources(c, path, "/networks/%s?group_name="+groupName)
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
			subnets = make([]*Subnet, 0)
		}
		return c.JSON(200, subnets)
	}

}

func deleteSubnet(c *echo.Context) error {
	group_name := c.Param("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, c.Param("network_id"), c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

func createSubnet(c *echo.Context) error {
	postParams := c.Request.Form
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, c.Param("group_name"), networkPath, c.Param("network_id"), postParams.Get("name"), config.ApiVersion)
	log.Printf("Create Subnet request with params: %s\n", postParams)
	log.Printf("Create Subnet path: %s\n", path)
	data := Subnet{
		Properties: map[string]interface{}{
			"addressPrefix": postParams.Get("address_prefix"),
			//"dhcpOptions":   postParams.Get("dhcp_options"),
		},
	}

	by, err := json.Marshal(data)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	request, err := http.NewRequest("PUT", path, reader)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while creating subnet: %v", err))
	}
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while creating subnet: %v", err))
	}

	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return lib.GenericException(fmt.Sprintf("Subnet creation failed: %s", string(b)))
	}

	var dat *Subnet
	if err := json.Unmarshal(b, &dat); err != nil {
		return lib.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	return c.JSON(response.StatusCode, dat)
}

func getSubnets(c *echo.Context, group_name string, network_name string) ([]*Subnet, error) {
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, network_name, config.ApiVersion)
	log.Printf("Get Subents request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return nil, lib.GenericException(fmt.Sprintf("Error has occurred while getting subnet: %v", err))
	}
	defer resp.Body.Close()
	var m map[string][]*Subnet
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	subnets := m["value"]

	for _, subnet := range subnets {
		subnet.Href = fmt.Sprintf("/resource_groups/%s/networks/%s/subnets/%s", group_name, network_name, subnet.Name)
	}

	return subnets, nil
}
