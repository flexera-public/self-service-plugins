package resources

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gm "github.com/rightscale/go_middleware"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
	am "github.com/rightscale/self-service-plugins/azure_v2/middleware"

	"testing"
)

func TestAzureResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Azure Resources Suite")
}

// Port plugin listens on for testing...differs from default port in dev
const (
	PluginPort     = "8081"
	subscriptionID = "test"
)

var AccessTokenTest = "fake"
var CredsTest = am.Credentials{}

// Run once for all tests
// Can't shutdown http servers just yet https://github.com/golang/go/issues/4674
var _ = BeforeSuite(func() {
	*config.SubscriptionIDCred = subscriptionID
	plugin := httpServer()
	go plugin.Run(":" + PluginPort)
})

// basic azure plugin HTTP client
type AzureClient struct {
	client *http.Client
	port   string
}

// Read HTTP response
type Response struct {
	Body    string
	Status  int
	Headers http.Header
	Cookies []*http.Cookie
}

// Instantiate new azure client
func NewAzureClient() *AzureClient {
	return &AzureClient{
		client: http.DefaultClient,
		port:   PluginPort,
	}
}

// Send GET request to cloud
func (c *AzureClient) Get(url string) (*Response, error) {
	return c.do("GET", url, "")
}

// Send POST request to cloud
func (c *AzureClient) Post(url, body string) (*Response, error) {
	return c.do("POST", url, body)
}

// Send DELETE request to cloud
func (c *AzureClient) Delete(url string) (*Response, error) {
	return c.do("DELETE", url, "")
}

// Helper generic send request method
func (c *AzureClient) do(verb, url, body string) (*Response, error) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(verb, "http://localhost:"+c.port+url, reader)
	if err != nil {
		return nil, err
	}
	if AccessTokenTest != "" {
		req.AddCookie(&http.Cookie{Name: "AccessToken", Value: AccessTokenTest})
	} else {
		req.AddCookie(&http.Cookie{Name: "TenantID", Value: CredsTest.TenantID})
		req.AddCookie(&http.Cookie{Name: "ClientID", Value: CredsTest.ClientID})
		req.AddCookie(&http.Cookie{Name: "ClientSecret", Value: CredsTest.ClientSecret})
		req.AddCookie(&http.Cookie{Name: "RefreshToken", Value: CredsTest.RefreshToken})
		req.AddCookie(&http.Cookie{Name: "SubscriptionID", Value: CredsTest.Subscription})
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		Body:    string(respBody),
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Cookies: resp.Cookies(),
	}, nil
}

// Factory method for application
// Makes it possible to do integration testing.
// TODO: code duplication...the same method is placed in the main package
func httpServer() *echo.Echo {
	// Setup middleware
	e := echo.New()
	e.Use(gm.RequestID)                 // Put that first so loggers can log request id
	e.Use(gm.HttpLogger(config.Logger)) // Log to syslog
	e.Use(am.AzureClientInitializer())
	e.Use(em.Recover())

	e.SetHTTPErrorHandler(eh.AzureErrorHandler(e)) // override default error handler
	// Setup routes
	SetupSubscriptionRoutes(e)
	SetupInstanceRoutes(e)
	SetupGroupsRoutes(e)
	SetupStorageAccountsRoutes(e)
	SetupProviderRoutes(e)
	SetupNetworkRoutes(e)
	SetupSubnetsRoutes(e)
	SetupIPAddressesRoutes(e)
	SetupAuthRoutes(e)
	SetupNetworkInterfacesRoutes(e)
	SetupOperationRoutes(e)
	SetupAvailabilitySetRoutes(e)

	return e
}
