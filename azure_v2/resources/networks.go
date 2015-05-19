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
	"github.com/rightscale/self-service-plugins/azure_v2/middleware"
)

const (
	networkPath = "providers/Microsoft.Network/virtualNetworks"
)

type Network struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Location   string      `json:"location"`
	Tags       interface{} `json:"tags,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupNetworkRoutes(e *echo.Echo) {
	e.Get("/networks", listNetworks)
	e.Post("/networks", createNetwork)
}

func listNetworks(c *echo.Context) *echo.HTTPError {
	requestParams := c.Request.Form
	if requestParams.Get("group_name") != "" {
		code, resp := getNetworks(c, requestParams.Get("group_name"))
		return c.JSON(code, resp)
	} else {
		code, resp := getResources(c, "")
		var networks []*Network
		for _, resource_group := range resp {
			_, resp := getNetworks(c, resource_group.Name)
			networks = append(networks, resp...)
		}
		// [].to_json => null ... why?
		return c.JSON(code, networks)
	}

}

func createNetwork(c *echo.Context) *echo.HTTPError {
	postParams := c.Request.Form
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, postParams.Get("group_name"), networkPath, postParams.Get("name"), config.ApiVersion)
	log.Printf("Create Network request with params: %s\n", postParams)
	log.Printf("Create Network path: %s\n", path)
	var subnets []map[string]interface{}
	data := Network{
		Name:     postParams.Get("name"),
		Location: postParams.Get("location"),
		Properties: map[string]interface{}{
			"addressSpace": map[string]interface{}{
				"addressPrefixes": []string{"10.0.0.0/16"},
			},
			"subnets": append(subnets, map[string]interface{}{
				"name": postParams.Get("name"),
				"properties": map[string]interface{}{
					"addressPrefix": "10.0.0.0/16",
				},
			}),
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
	var dat *Network
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}

	return c.JSON(response.StatusCode, dat)
}

func getNetworks(c *echo.Context, group_name string) (int, []*Network) {
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, config.ApiVersion)
	log.Printf("Get Networks request: %s\n", path)
	resp, err := client.Get(path)

	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	var m map[string][]*Network
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	return resp.StatusCode, m["value"]
}
