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
	Href           string      `json:"href,omitempty"`
}

// SetupSubscriptionRoutes declares routes for Subscription resource
func SetupSubscriptionRoutes(e *echo.Group) {
	// get a current subscription
	e.Get("/subscription", getSubscription)
}

func getSubscription(c *echo.Context) error {
	return Get(c, new(Subscription))
}

// GetPath returns full path to the sigle subscription
func (s *Subscription) GetPath() string {
	return fmt.Sprintf("%s/%s/%s?api-version=%s", config.BaseURL, subscriptionsPath, *config.SubscriptionIDCred, "2015-01-01")
}

// HandleResponse manage raw cloud response
func (s *Subscription) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &s); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	s.Href = s.GetHref("")
	return nil
}

// GetContentType returns content type of subscription
func (s *Subscription) GetContentType() string {
	return "vnd.rightscale.subscription+json"
}

// GetHref returns subscription href
func (s *Subscription) GetHref(_ string) string {
	return "/subscription"
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (s *Subscription) GetResponseParams() interface{} { return s }

//GetCollectionPath is a fake function to support AzureResource by Subscription
func (s *Subscription) GetCollectionPath(groupName string) string { return "" }

//GetRequestParams is a fake function to support AzureResource by Subscription
func (s *Subscription) GetRequestParams(c *echo.Context) (interface{}, error) { return "", nil }
