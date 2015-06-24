package resources

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-errors/errors"
	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gm "github.com/rightscale/go_middleware"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	am "github.com/rightscale/self-service-plugins/azure_v2/middleware"

	"testing"
)

func TestAzureResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Azure Resources Suite")
}

// Port plugin listens on for testing...differs from default port in dev
const PluginPort = "8081"

// Run gdo once for all tests
// Can't shutdown http servers just yet https://github.com/golang/go/issues/4674
var _ = BeforeSuite(func() {
	*config.SubscriptionIdCred = "test"
	plugin := HttpServer()
	go plugin.Run(":" + PluginPort)
})

// basic gdo plugin HTTP client
type AzureClient struct {
	client *http.Client
	port   string
}

// Read HTTP response
type Response struct {
	Body    string
	Status  int
	Headers http.Header
}

// Instantiate new azure client
func NewAzureClient() *AzureClient {
	return &AzureClient{
		client: http.DefaultClient,
		port:   PluginPort,
	}
}

// Send GET request to gdo
func (c *AzureClient) Get(url string) (*Response, error) {
	return c.do("GET", url, "")
}

// Send POST request to gdo
func (c *AzureClient) Post(url, body string) (*Response, error) {
	return c.do("POST", url, body)
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
	req.AddCookie(&http.Cookie{Name: "AccessToken", Value: "fake"})
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
	}, nil
}

// Factory method for application
// Makes it possible to do integration testing.
// TODO: code duplication...the same method is placed in the main package
func HttpServer() *echo.Echo {
	// Setup middleware
	e := echo.New()
	e.Use(gm.RequestID)                 // Put that first so loggers can log request id
	e.Use(gm.HttpLogger(config.Logger)) // Log to syslog
	e.Use(am.AzureClientInitializer())
	e.Use(em.Recover())

	e.SetHTTPErrorHandler(AzureErrorHandler(e)) // override default error handler
	// Setup routes
	SetupSubscriptionRoutes(e)
	SetupInstanceRoutes(e)
	SetupGroupsRoutes(e)
	SetupStorageAccountsRoutes(e)
	SetupProviderRoutes(e)
	SetupNetworkRoutes(e)
	SetupSubnetsRoutes(e)
	SetupIpAddressesRoutes(e)
	SetupAuthRoutes(e)
	SetupNetworkInterfacesRoutes(e)
	SetupOperationRoutes(e)

	return e
}

type GenericError struct {
	echo.HTTPError
	StackTrace string `json:"StackTrace,omitempty"`
}

func AzureErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c *echo.Context) {
		ge := new(GenericError)
		ge.Code = http.StatusInternalServerError //default status code is 500
		ge.Message = http.StatusText(ge.Code)    // default message is 'Internal Server Error'
		switch error := err.(type) {
		case *errors.Error:
			if he, ok := error.Err.(*echo.HTTPError); ok {
				ge.Code = he.Code
				ge.Message = he.Message
			}
			if e.Debug() && ge.Code >= 500 {
				ge.StackTrace = error.ErrorStack()
			}
		case *echo.HTTPError:
			ge.Code = error.Code
			ge.Message = error.Message
		case error:
			if e.Debug() && ge.Code >= 500 {
				ge.Message = err.Error() //show original error message in case of debug mode https://github.com/labstack/echo/blob/1e117621e9006481bfc0fd8e6bafab48c1848639/echo.go#L161
			}
		}

		c.JSON(ge.Code, ge)
	}
}
