package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
	"github.com/rightscale/godo"
)

// Name of cookie created by SS that contains the credentials needed to send API requests to DO
const CredCookieName = "ServiceCred"

// Digital Ocean API endpoint, exposed so tests can change it
var DOBaseURL *url.URL

// Middleware that creates DO client using credentials in cookie
func DOClientInitializer(dump bool) echo.Middleware {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token, err := c.Request.Cookie(CredCookieName)
			if err != nil {
				return fmt.Errorf("cookie '%s' is missing", CredCookieName)
			}
			t := &oauth.Transport{Token: &oauth.Token{AccessToken: token.Value}}
			client := godo.NewClient(t.Client())
			if DOBaseURL != nil {
				client.BaseURL = DOBaseURL
			}
			if dump {
				client.OnRequestCompleted(dumpRequestResponse)
			}
			c.Set("doC", client)
			return h(c)
		}
	}
}

// Retrieve client initialized by middleware, send error response if not found
// This function should be used by controller actions that need to use the client
func GetDOClient(c *echo.Context) (*godo.Client, error) {
	client, _ := c.Get("doC").(*godo.Client)
	if client == nil {
		return nil, fmt.Errorf("failed to retrieve Digital Ocean client, check middleware")
	}
	return client, nil
}

// HTTP Request and response details to be dumped
type RequestResponse struct {
	Verb, Uri  string      // request http verb and full uri with query string
	ReqHeader  http.Header // headers before std additions, such as user-agent
	ReqBody    string      // not []byte so that json.Marshal doesn't produce base64
	Status     int         // numerical response status
	RespHeader http.Header // full response headers
	RespBody   string      // not []byte so that json.Marshal doesn't produce base64
}

// Dump the HTTP request and response
func dumpRequestResponse(req *http.Request, resp *http.Response) {
	reqBody, err := dumpReqBody(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load request body for dump: %s\n", err)
	}
	respBody, err := dumpRespBody(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load response body for dump: %s\n", err)
	}
	dumped := RequestResponse{
		Verb:       req.Method,
		Uri:        req.URL.String(),
		ReqHeader:  req.Header,
		ReqBody:    string(reqBody),
		Status:     resp.StatusCode,
		RespHeader: resp.Header,
		RespBody:   string(respBody),
	}
	b, err := json.MarshalIndent(dumped, "", "    ")
	if err == nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(b))
	} else {
		fmt.Fprintf(os.Stderr, "Failed to dump request content - %s\n", err)
	}
}

// Dump request body, strongly inspired from httputil.DumpRequest
func dumpReqBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}
	var save io.ReadCloser
	var err error
	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	var dest io.Writer = &b
	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}
	_, err = io.Copy(dest, req.Body)
	if chunked {
		dest.(io.Closer).Close()
		io.WriteString(&b, "\r\n")
	}
	req.Body = save
	return b.Bytes(), err
}

// Dump response body, strongly inspired from httputil.DumpResponse
func dumpRespBody(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, nil
	}
	var b bytes.Buffer
	savecl := resp.ContentLength
	var save io.ReadCloser
	var err error
	save, resp.Body, err = drainBody(resp.Body)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&b, resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = save
	resp.ContentLength = savecl
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// debugging only.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
