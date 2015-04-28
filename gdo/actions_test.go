package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/gdo/middleware"
)

const (
	doGetActions = `{"actions":[{"id":48912908,"status":"completed","type":"create","started_at":"2015-04-24T22:55:43Z","completed_at":"2015-04-24T22:56:39Z","resource_id":5004436,"resource_type":"droplet","region":{"name":"New York 3","slug":"nyc3","sizes":["512mb","1gb","2gb","4gb","8gb","16gb","32gb","48gb","64gb"],"features":["virtio","private_networking","backups","ipv6","metadata"],"available":true},"region_slug":"nyc3"}],"links":{},"meta":{"total":1}}`
	doGetAction  = `{"action":{"id":48913731,"status":"completed","type":"power_off","started_at":"2015-04-24T23:11:41Z","completed_at":"2015-04-24T23:12:01Z","resource_id":5004436,"resource_type":"droplet","region":{"name":"New York 3","slug":"nyc3","sizes":["512mb","1gb","2gb","4gb","8gb","16gb","32gb","48gb","64gb"],"features":["virtio","private_networking","backups","ipv6","metadata"],"available":true},"region_slug":"nyc3"}}`
	doGetEmpty   = `{"actions":[],"links":{},"meta":{"total":0}}`
	doNotFound   = `{"id":"not_found","message":"The resource you were accessing could not be found."}`
)

var _ = Describe("actions", func() {

	var do *ghttp.Server
	var client *GDOClient

	BeforeEach(func() {
		do = ghttp.NewServer()
		u, err := url.Parse(do.URL())
		Expect(err).NotTo(HaveOccurred())
		middleware.DOBaseURL = u
		client = NewGDOClient()
	})

	AfterEach(func() {
		do.Close()
	})

	Describe("listing", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/actions"),
					ghttp.RespondWith(http.StatusOK, doGetActions),
				),
			)
		})

		It("lists all actions", func() {
			resp, err := client.Get("/actions")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			var actions map[string]interface{}
			err = json.Unmarshal([]byte(doGetActions), &actions)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(actions["actions"])
			Expect(err).NotTo(HaveOccurred())
			Ω(resp.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/actions"),
					ghttp.RespondWith(http.StatusOK, doGetEmpty),
				),
			)
		})

		It("returns empty array", func() {
			resp, err := client.Get("/actions")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(resp.Status).Should(Equal(200))
			Ω(resp.Body).Should(Equal("[]\n"))
		})
	})

	Describe("retrieving", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/actions/36804636"),
					ghttp.RespondWith(http.StatusOK, doGetAction),
				),
			)
		})

		It("retrieves an existing action", func() {
			resp, err := client.Get("/actions/36804636")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(resp.Status).Should(Equal(200))
			var action map[string]interface{}
			err = json.Unmarshal([]byte(doGetAction), &action)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(action["action"])
			Expect(err).NotTo(HaveOccurred())
			Ω(resp.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/actions/1"),
					ghttp.RespondWith(http.StatusNotFound, doNotFound),
				),
			)
		})

		It("returns 404", func() {
			resp, err := client.Get("/actions/1")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(resp.Status).Should(Equal(404))
			Ω(resp.Body).Should(Equal(http.StatusText(http.StatusNotFound)))
		})
	})
})
