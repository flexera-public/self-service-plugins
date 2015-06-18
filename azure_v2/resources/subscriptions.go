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

const (
	subscriptionsPath = "subscriptions"
)

type Subscription struct {
	Id             string      `json:"id"`
	Name           string      `json:"displayName"`
	State          string      `json:"state"`
	SubscriptionId string      `json:"subscriptionId"`
	Policies       interface{} `json:"subscriptionPolicies"`
}

func SetupSubscriptionRoutes(e *echo.Echo) {
	//e.Get("/subscriptions", listSubscriptions)
	// get a current subscription
	e.Get("/subscription", GetSubscription)
}

func listSubscriptions(c *echo.Context) error {
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/%s?api-version=%s", config.BaseUrl, subscriptionsPath, config.ApiVersion)
	log.Printf("Get Subscriptions request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while getting subscriptions: %v", err))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var dat map[string][]*Subscription
	if err := json.Unmarshal(body, &dat); err != nil {
		return lib.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	return c.JSON(resp.StatusCode, dat["value"])
}

func GetSubscription(c *echo.Context) error {
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/%s/%s?api-version=%s", config.BaseUrl, subscriptionsPath, *config.SubscriptionIdCred, "2015-01-01")
	log.Printf("Get Subscription request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while getting subscription: %v", err))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var dat *Subscription
	if err := json.Unmarshal(body, &dat); err != nil {
		return lib.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	return c.JSON(resp.StatusCode, dat)
}
