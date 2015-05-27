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

type ResourceGroup struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty,omitempty"`
	Location   string      `json:"location,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupGroupsRoutes(e *echo.Echo) {
	e.Get("/resource_groups", listResourceGroups)
}

func listResourceGroups(c *echo.Context) error {
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
	groups, err := lib.GetResources(c, path)
	if err != nil {
		return err
	}
	return c.JSON(200, groups)
}

func getResources(c *echo.Context, resource_group_id string) (int, []*ResourceGroup) {
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, resource_group_id, "2015-01-01")
	log.Printf("Get Resource Groups request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var dat map[string][]*ResourceGroup
	if err := json.Unmarshal(body, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}
	return resp.StatusCode, dat["value"]
}
