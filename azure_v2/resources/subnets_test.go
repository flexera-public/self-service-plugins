package resources

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

const (
	listSubnetsResponse   = `{"value":[{"etag":"W/\"bfa01c55-48e2-41fc-8f0b-2189de29e495\"","href":"/resource_groups/Group-3/networks/khrvi-3/subnets/sub1","id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Network/virtualNetworks/khrvi-3/subnets/sub1","name":"sub1","properties":{"addressPrefix":"10.0.0.0/16","provisioningState":"Succeeded"}}]}`
	listOneSubnetResponse = `{"id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Network/virtualNetworks/khrvi-3/subnets/sub1","name":"sub1","etag":"W/\"bfa01c55-48e2-41fc-8f0b-2189de29e495\"","properties":{"addressPrefix":"10.0.0.0/16","provisioningState":"Succeeded"}, "href": "/resource_groups/Group-3/networks/khrvi-3/subnets/sub1"}`
)

var _ = Describe("subnets", func() {

	var do *ghttp.Server
	var client *AzureClient
	var response *Response
	var err error

	BeforeEach(func() {
		do = ghttp.NewServer()
		config.BaseURL = do.URL()
		client = NewAzureClient()
	})

	AfterEach(func() {
		do.Close()
	})

	Describe("listing", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi-3/subnets"),
					ghttp.RespondWith(http.StatusOK, listSubnetsResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-3/networks/khrvi-3/subnets")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.subnet+json;type=collection"))
		})

		It("lists all subnets inside network", func() {
			subnets := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listSubnetsResponse), &subnets)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(subnets["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/"+networkPath),
					ghttp.RespondWith(http.StatusOK, listNetworksResponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi-3/subnets"),
					ghttp.RespondWith(http.StatusOK, listSubnetsResponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/net2/subnets"),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
			response, err = client.Get("/subnets")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(3))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.subnet+json;type=collection"))
		})

		It("lists all subnets inside network", func() {
			subnets := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listSubnetsResponse), &subnets)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(subnets["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi-3/subnets"),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-3/networks/khrvi-3/subnets")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one subnet", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi-3/subnets/sub1"),
					ghttp.RespondWith(http.StatusOK, listOneSubnetResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-3/networks/khrvi-3/subnets/sub1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.subnet+json"))
		})

		It("retrieves an existing subnet", func() {
			subnet := new(subnetResponseParams)
			err := json.Unmarshal([]byte(listOneSubnetResponse), &subnet)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(subnet)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi-3/subnets/sub2"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-3/networks/khrvi-3/subnets/sub2")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(404))
			Ω(response.Body).Should(Equal("{\"Code\":404,\"Message\":\"Could not find resource with id: sub2\"}\n"))
		})
	})

	Describe("creating", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi-3/subnets/sub1"),
					ghttp.VerifyJSONRepresenting(subnetRequestParams{
						Properties: map[string]interface{}{
							"addressPrefix": "10.0.0.0/16",
							"networkSecurityGroup": map[string]interface{}{
								"id": "/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-3/providers/Microsoft.Network/networkSecurityGroups/khrvi1",
							},
						},
					}),
					ghttp.RespondWith(201, listOneSubnetResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-3/networks/khrvi-3/subnets", "{\"name\": \"sub1\", \"address_prefix\": \"10.0.0.0/16\", \"network_security_group_id\": \"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-3/providers/Microsoft.Network/networkSecurityGroups/khrvi1\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource subnet href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-3/networks/khrvi-3/subnets/sub1"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi-3/subnets/sub1"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-3/networks/khrvi-3/subnets/sub1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 204 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(204))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})
})
