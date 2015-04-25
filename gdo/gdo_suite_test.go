package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rightscale/gdo/middleware"

	"testing"
)

func TestGdo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gdo Suite")
}

// Port gdo listens on for testing
const PluginPort = "8080"

// Run gdo once for all tests
// Can't shutdown http servers just yet https://github.com/golang/go/issues/4674
var _ = BeforeSuite(func() {
	plugin := HttpServer()
	go plugin.Run(":" + PluginPort)
})

// basic gdo plugin HTTP client
type GDOClient struct {
	client *http.Client
	port   string
}

// Read HTTP response
type Response struct {
	Body    string
	Headers http.Header
}

// Instantiate new gdo client
func NewGDOClient() *GDOClient {
	return &GDOClient{
		client: http.DefaultClient,
		port:   PluginPort,
	}
}

// Send GET request to gdo
func (c *GDOClient) Get(url string) (*Response, error) {
	return c.do("GET", url, "")
}

// Send POST request to gdo
func (c *GDOClient) Post(url, body string) (*Response, error) {
	return c.do("POST", url, body)
}

// Helper generic send request method
func (c *GDOClient) do(verb, url, body string) (*Response, error) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(verb, "http://localhost:"+c.port+url, reader)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: middleware.CredCookieName, Value: "fake"})
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
		Headers: resp.Header,
	}, nil
}
