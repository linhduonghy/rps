package request

import "rps/internal/http/common"

type CreatePolicy struct {
	ResourceType string      `json:"resourceType"`
	ResourceID   string      `json:"resourceId"`
	Version      string      `json:"version"`
	Statements   []Statement `json:"statements"`
}

type UpdatePolicy struct {
	Version    string      `json:"version"`
	Statements []Statement `json:"statements"`
}

type Statement struct {
	SID        string      `json:"sid"`
	Effect     string      `json:"effect"`
	Actions    []string    `json:"actions"`
	Resources  []string    `json:"resources"`
	Principals []Principal `json:"principals"`
}

type Principal struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type ListPolicies struct {
	*common.ListRequest

	ResourceType string
	ResourceID   string
}

func NewListPolicies(page int32, size int32, resourceType string, resourceID string) *ListPolicies {
	return &ListPolicies{
		ListRequest:  common.NewListRequest(page, size, "", ""),
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}
}
