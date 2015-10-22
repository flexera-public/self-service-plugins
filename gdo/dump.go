package main

import "net/http"

// HTTP Request and response details
type RequestResponse struct {
	Verb, Uri  string      // request http verb and full uri with query string
	ReqHeader  http.Header // headers before std additions, such as user-agent
	ReqBody    string      // not []byte so that json.Marshal doesn't produce base64
	Status     int         // numerical response status
	RespHeader http.Header // full response headers
	RespBody   string      // not []byte so that json.Marshal doesn't produce base64
}
