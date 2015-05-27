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
}

func SetupSubnetsRoutes(e *echo.Echo) {
	e.Get("/subnets", listSubnets)
	e.Post("/subnets", createSubnet)

	//nested routes
	group := e.Group("/resource_groups/:group_name/networks/:network_id/subnets")
	group.Get("", listSubnets)
	// group.Post("", createInstance)
	// group.Delete("/:id", deleteInstance)
}

func listSubnets(c *echo.Context) error {
	group_name := c.Param("group_name")
	if group_name != "" {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, c.Param("network_id"), config.ApiVersion)
		subnets, err := lib.GetResources(c, path)
		if err != nil {
			return err
		}
		return c.JSON(200, subnets)
	} else {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
		resp, _ := lib.GetResources(c, path)
		var subnets []*Subnet
		for _, resource_group := range resp {
			group := resource_group.(map[string]interface{})
			groupName := group["name"].(string)
			path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, networkPath, config.ApiVersion)
			networks, _ := lib.GetResources(c, path)

			for _, network := range networks {
				network := network.(map[string]interface{})
				_, resp := getSubnets(c, groupName, network["name"].(string))
				subnets = append(subnets, resp...)
			}
		}
		// [].to_json => null ... why?
		return c.JSON(200, subnets)
	}

}

func createSubnet(c *echo.Context) error {
	postParams := c.Request.Form
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, postParams.Get("group_name"), networkPath, postParams.Get("network_name"), postParams.Get("name"), config.ApiVersion)
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
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("PUT:", err)
	}

	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return lib.GenericException(fmt.Sprintf("Subnet creation failed: %s", string(b)))
	}

	var dat *Subnet
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}

	return c.JSON(response.StatusCode, dat)
}

func getSubnets(c *echo.Context, group_name string, network_name string) (int, []*Subnet) {
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/subnets?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, network_name, config.ApiVersion)
	log.Printf("Get Subents request: %s\n", path)
	resp, err := client.Get(path)

	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	var m map[string][]*Subnet
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	return resp.StatusCode, m["value"]
}
