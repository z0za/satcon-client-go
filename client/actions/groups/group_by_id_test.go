package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/groups"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("GroupById", func() {
	var (
		orgID          string
		groupId        string
		fakeAuthClient authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		orgID = "someorg"
		groupId = "someid"
	})

	Describe("NewGroupByIdVariables", func() {
		It("Returns a correctly populated instance of GroupByIdVariables", func() {
			vars := NewGroupByIdVariables(orgID, groupId)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryGroupById))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.ID).To(Equal(groupId))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
				"uuid":  "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
				"orgId",
				"name",
				"created",
				"clusters{id,orgId,clusterId,name,metadata}",
			))
		})
	})

	Describe("GroupById", func() {
		var (
			c             GroupService
			h             *webfakes.FakeHTTPClient
			response      *http.Response
			groupResponse GroupByIdResponse
		)

		BeforeEach(func() {
			groupResponse = GroupByIdResponse{
				Data: &GroupByIdResponseData{
					&types.Group{
						UUID:  groupId,
						OrgID: orgID,
						Name:  "group1",
						Clusters: []types.Cluster{
							{
								ID:        "cid",
								OrgID:     "oid",
								ClusterID: "cid",
								Name:      "cluster1",
							},
						},
					},
				},
			}

			respBodyBytes, err := json.Marshal(groupResponse)
			Expect(err).NotTo(HaveOccurred())
			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			h = &webfakes.FakeHTTPClient{}
			Expect(h.DoCallCount()).To(Equal(0))
			h.DoReturns(response, nil)

			c, _ = NewClient("https://foo.bar", h, &fakeAuthClient)
		})

		It("Makes a valid http request", func() {
			_, err := c.GroupById(orgID, groupId)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the group", func() {
			groups, _ := c.GroupById(orgID, groupId)
			expected := groupResponse.Data.Group
			Expect(groups).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Kablooie!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.GroupById(orgID, groupId)
				Expect(err).To(MatchError("Kablooie!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(GroupsResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				groups, err := c.GroupById(orgID, groupId)
				Expect(err).NotTo(HaveOccurred())
				Expect(groups).To(BeNil())
			})
		})
	})
})
