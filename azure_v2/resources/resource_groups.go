package resources

import (
	"log"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/middleware"
)

type ResourceGroup struct {
	Id string
	Name string
    Location string
    Properties interface {}
}

func listResourceGroups(c *echo.Context) *echo.HTTPError {
	//requestParams := c.Request.Form
	
	code, body := getResources(c, "")  
	//c.String(resp.StatusCode, string(body))
	return c.JSON(code, body)
}

func getResources(c *echo.Context, resource_group_id string) (int, string){
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, resource_group_id, "2015-01-01")
	log.Printf("Get Resource Groups request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}