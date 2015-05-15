package resources

import (
	"encoding/json"
	"net/http"
	//"net/url"
	//"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	//"github.com/rightscale/gdo/middleware"
)

const (
	listInstancesResponse = `{"values":[]}`
)

var _ = Describe("instances", func() {

	var do *ghttp.Server
	var client *AzureClient

	BeforeEach(func() {
		do = ghttp.NewServer()
		config.BaseUrl = do.URL()
		client = NewAzureClient()
	})

	AfterEach(func() {
		do.Close()
	})

	Describe("listing", func() {
		BeforeEach(func() {
			subscriptionIdW := "test"
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/" + subscriptionIdW + "/resourceGroups/Group-1/" + virtualMachinesPath),
					ghttp.RespondWith(http.StatusOK, listInstancesResponse),
				),
			)
		})

		It("lists all instances inside one resource group", func() {
			resp, err := client.Get("/instances?group_name=Group-1")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			var actions map[string]interface{}
			err = json.Unmarshal([]byte(listInstancesResponse), &actions)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(actions)
			Expect(err).NotTo(HaveOccurred())
			Ω(resp.Body).Should(MatchJSON(expected))
		})
	})
})
