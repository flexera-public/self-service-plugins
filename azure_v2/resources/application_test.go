package resources

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

const (
	authRsponse              = `{"access_token": "test_access_token"}`
	servicePrincipalResponse = `{"value":[{ "objectType":"ServicePrincipal","objectId":"b24988ac-6180-42a0-ab88-20f7382dd24c"}]}`
)

var _ = Describe("application", func() {

	var do *ghttp.Server
	var client *AzureClient
	var response *Response
	var err error

	BeforeEach(func() {
		do = ghttp.NewServer()
		config.AuthHost = do.URL()
		config.GraphURL = do.URL()
		config.BaseURL = do.URL()
		client = NewAzureClient()
	})

	AfterEach(func() {
		do.Close()
	})

	Describe("register", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/test_tenant/oauth2/token"),
					ghttp.RespondWith(http.StatusOK, authRsponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/test_tenant/servicePrincipals", "api-version=1.5&$filter=appId%20eq%20'test'"),
					ghttp.RespondWith(http.StatusOK, servicePrincipalResponse),
				),

				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", MatchRegexp(`/subscriptions/test_subscription/providers/microsoft.authorization/roleassignments/\w`)),
					ghttp.RespondWith(201, ""),
				),
			)
			//TODO: try to use analog of before_all
			response, err = client.Post("/application/register", "{\"tenant\":\"test_tenant\", \"client_id\": \"test\", \"client_secret\": \"test_secret\", \"subscription\": \"test_subscription\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(3))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns empty response in case of 201", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})
})
