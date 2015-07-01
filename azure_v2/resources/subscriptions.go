package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	subscriptionsPath = "subscriptions"
)

// Subscription is base struct for Azure Subscription resource
type Subscription struct {
	ID             string      `json:"id"`
	Name           string      `json:"displayName"`
	State          string      `json:"state"`
	SubscriptionID string      `json:"subscriptionId"`
	Policies       interface{} `json:"subscriptionPolicies"`
}

// SetupSubscriptionRoutes declares routes for Subscription resource
func SetupSubscriptionRoutes(e *echo.Echo) {
	// get a current subscription
	e.Get("/subscription", getSubscription)
}

// getSubscription return info about subscription provided in creds
func getSubscription(c *echo.Context) error {
	path := fmt.Sprintf("%s/%s/%s?api-version=%s", config.BaseURL, subscriptionsPath, *config.SubscriptionIDCred, "2015-01-01")
	body, err := GetResource(c, path)
	if err != nil {
		return nil
	}
	var dat *Subscription
	if err := json.Unmarshal(body, &dat); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}

	return c.JSON(200, dat)
}
