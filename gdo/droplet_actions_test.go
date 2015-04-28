package main

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/gdo/middleware"
)

const (
	doGetDropletActions      = `{"actions":[{"id":48912908,"status":"completed","type":"create","started_at":"2015-04-24T22:55:43Z","completed_at":"2015-04-24T22:56:39Z","resource_id":5004436,"resource_type":"droplet","region":{"name":"New York 3","slug":"nyc3","sizes":["512mb","1gb","2gb","4gb","8gb","16gb","32gb","48gb","64gb"],"features":["virtio","private_networking","backups","ipv6","metadata"],"available":true},"region_slug":"nyc3"}],"links":{},"meta":{"total":1}}`
	doGetDropletActionsEmpty = `{"actions":[],"links":{},"meta":{"total":0}}`
	doGetDropletAction       = `{"action":{"id":48913731,"status":"completed","type":"power_off","started_at":"2015-04-24T23:11:41Z","completed_at":"2015-04-24T23:12:01Z","resource_id":5004436,"resource_type":"droplet","region":{"name":"New York 3","slug":"nyc3","sizes":["512mb","1gb","2gb","4gb","8gb","16gb","32gb","48gb","64gb"],"features":["virtio","private_networking","backups","ipv6","metadata"],"available":true},"region_slug":"nyc3"}}`
)

var _ = Describe("droplet_actions", func() {

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

})
