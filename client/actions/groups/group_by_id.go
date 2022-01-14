package groups

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryGroupById       = "group"
	GroupByIdVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"uuid":{{json .ID}}{{end}}`
)

type GroupByIdVariables struct {
	actions.GraphQLQuery
	OrgID string
	ID    string
}

func NewGroupByIdVariables(orgID string, id string) GroupByIdVariables {
	vars := GroupByIdVariables{
		OrgID: orgID,
		ID:    id,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryGroupById
	vars.Args = map[string]string{
		"orgId": "String!",
		"uuid":  "String!",
	}
	vars.Returns = []string{
		"uuid",
		"orgId",
		"name",
		"created",
		"clusters{id,orgId,clusterId,name,metadata}",
	}

	return vars
}

type GroupByIdResponse struct {
	Data *GroupByIdResponseData `json:"data,omitempty"`
}

type GroupByIdResponseData struct {
	Group *types.Group `json:"GroupById,omitempty"`
}

func (c *Client) GroupById(orgID string, id string) (*types.Group, error) {
	var response GroupByIdResponse

	vars := NewGroupByIdVariables(orgID, id)

	err := c.DoQuery(GroupByIdVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Group, err
	}

	return nil, err
}
