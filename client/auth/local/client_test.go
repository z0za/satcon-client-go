package local_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/IBM/satcon-client-go/client/auth/local"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("returns a LocalRazeeClient", func() {
		local, err := local.NewClient("http://foo.bar", "user", "password")
		Expect(err).NotTo(HaveOccurred())
		Expect(local).NotTo(BeNil())
	})

	Describe("Local razee errors", func() {
		It("Should error when url is empty", func() {
			local, err := local.NewClient("", "user", "password")
			Expect(err).To(HaveOccurred())
			Expect(local).To(BeNil())
		})

		It("Should error when login is empty", func() {
			local, err := local.NewClient("http://foo.bar", "", "password")
			Expect(err).To(HaveOccurred())
			Expect(local).To(BeNil())
		})

		It("Should error when password is empty", func() {
			local, err := local.NewClient("http://foo.bar", "user", "")
			Expect(err).To(HaveOccurred())
			Expect(local).To(BeNil())
		})
	})

	Describe("Local razee testing", func() {
		var token string
		var h *webfakes.FakeHTTPClient
		var response *http.Response
		var signInResponse local.SignInResponse

		BeforeEach(func() {
			var err error
			hmacSampleSecret := []byte("secret")
			tokenWithClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().Add(4 * time.Hour).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret
			token, err = tokenWithClaim.SignedString(hmacSampleSecret)

			h = &webfakes.FakeHTTPClient{}
			response = &http.Response{
				Header: http.Header{},
			}

			signInResponse = local.SignInResponse{
				Data: &local.SignInResponseData{
					Details: &local.SignInResponseDataDetails{
						Token: types.Token(token),
					},
				},
			}

			respBodyBytes, err := json.Marshal(signInResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			h.DoReturns(response, nil)
		})

		It("executes token retrieval", func() {
			localClient, err := local.NewClientWithHttpClient(h, "http://foo.bar", "user", "password")
			Expect(err).NotTo(HaveOccurred())
			Expect(localClient).NotTo(BeNil())
			request := http.Request{
				Header: http.Header{},
			}
			request.Header.Add("content-type", "application/json")
			// Call authenticate to check if the bearer token gets injected
			err = localClient.Authenticate(&request)
			Expect(err).NotTo(HaveOccurred())
			Expect(request.Header.Get(local.AuthorizationHeaderKey)).NotTo(BeEmpty())
			Expect(request.Header.Get(local.AuthorizationHeaderKey)).To(Equal("Bearer " + token))

			// Call authenticate to check if the bearer token gets injected again
			err = localClient.Authenticate(&request)
			Expect(err).NotTo(HaveOccurred())
			Expect(request.Header.Get(local.AuthorizationHeaderKey)).NotTo(BeEmpty())
			Expect(request.Header.Get(local.AuthorizationHeaderKey)).To(Equal("Bearer " + token))

			// Check that there was only one invocation (the second authenticate should come from the cache)
			Expect(len(h.Invocations())).To(Equal(1))
		})
	})
})
