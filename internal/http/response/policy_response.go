package response

import (
	"rps/internal/domain"
	"time"
)

type PolicyInfo struct {
	ResourceType string    `json:"resourceType"`
	ResourceID   string    `json:"resourceId"`
	Version      string    `json:"version"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PolicyDetail struct {
	ResourceType string      `json:"resourceType"`
	ResourceID   string      `json:"resourceId"`
	Version      string      `json:"version"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	Statements   []Statement `json:"statements"`
}

type Statement struct {
	SID        string      `json:"sid,omitempty"`
	Effect     string      `json:"effect"`
	Actions    []string    `json:"actions"`
	Resources  []string    `json:"resources"`
	Principals []Principal `json:"principals"`
}

type Principal struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func NewPolicyDetail(policy *domain.Policy) *PolicyDetail {
	var res = &PolicyDetail{
		ResourceType: policy.ResourceType,
		ResourceID:   policy.ResourceID,
		Version:      policy.Version,
		CreatedAt:    policy.CreatedAt,
		UpdatedAt:    policy.UpdatedAt,
		Statements:   make([]Statement, 0, len(policy.Statements)),
	}
	for _, stmt := range policy.Statements {

		statement := Statement{
			SID:       stmt.SID,
			Effect:    string(stmt.Effect),
			Actions:   stmt.Actions,
			Resources: stmt.Resources,
		}

		for _, principal := range stmt.Principals {
			statement.Principals = append(
				statement.Principals,
				Principal{
					Type: principal.Type,
					ID:   principal.ID,
				},
			)
		}
		res.Statements = append(res.Statements, statement)
	}
	return res
}
