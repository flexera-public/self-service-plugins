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
	listEmptyResponse = `{"value":[]}`
	listASsResponse   = `{"value":[{"href":"/resource_groups/Group-3/availability_sets/khrvi1","id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Compute/availabilitySets/khrvi1","location":"eastus","name":"khrvi1","properties":{"platformFaultDomainCount":3,"platformUpdateDomainCount":5,"virtualMachines":[]},"type":"Microsoft.Compute/availabilitySets"}]}`
	listOneASResponse = `{"id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Compute/availabilitySets/khrvi1","name":"khrvi1","location":"eastus","properties":{"platformFaultDomainCount":3,"platformUpdateDomainCount":5,"virtualMachines":[]},"href":"/resource_groups/Group-3/availability_sets/khrvi1"}`
)

var _ = Describe("availability_sets", func() {

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
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+availabilitySetPath),
					ghttp.RespondWith(http.StatusOK, listASsResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-3/availability_sets")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.availability_set+json;type=collection"))
		})

		It("lists all availability sets inside one resource group", func() {
			availabilitySets := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listASsResponse), &availabilitySets)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(availabilitySets["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			subscriptionID := "test"
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups"),
					ghttp.RespondWith(http.StatusOK, `{"value": [{"name":"Group-3"}]}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+availabilitySetPath),
					ghttp.RespondWith(http.StatusOK, listASsResponse),
				),
			)
			response, err = client.Get("/availability_sets")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(2))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.availability_set+json;type=collection"))
		})

		It("lists all availability sets inside one resource group", func() {
			availabilitySets := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listASsResponse), &availabilitySets)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(availabilitySets["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+availabilitySetPath),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-3/availability_sets")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one availability set", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+availabilitySetPath+"/khrvi1"),
					ghttp.RespondWith(http.StatusOK, listOneASResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-3/availability_sets/khrvi1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.availability_set+json"))
		})

		It("retrieves an existing availability set", func() {
			var availabilitySet map[string]interface{}
			err := json.Unmarshal([]byte(listOneASResponse), &availabilitySet)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(availabilitySet)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+availabilitySetPath+"/khrvi"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-3/availability_sets/khrvi")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(404))
			Ω(response.Body).Should(Equal("{\"Code\":404,\"Message\":\"Could not find resource with id: khrvi\"}\n"))
		})
	})

	Describe("creating", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+availabilitySetPath+"/khrvi1"),
					ghttp.RespondWith(201, listOneASResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-3/availability_sets", "{\"name\": \"khrvi1\", \"location\": \"westus\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource availability set href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-3/availability_sets/khrvi1"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+availabilitySetPath+"/khrvi1"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-3/availability_sets/khrvi1")
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
