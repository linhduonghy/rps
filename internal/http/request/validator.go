package request

import (
	"fmt"
	apperror "rps/internal/errors"
)

func (r *CreatePolicy) Validate() error {
	if r.ResourceType == "" {
		return apperror.InvalidInput("resource type is required")
	}

	if r.ResourceID == "" {
		return apperror.InvalidInput("resource id is required")
	}

	if r.Version == "" {
		return apperror.InvalidInput("version is required")
	}

	if len(r.Statements) == 0 {
		return apperror.InvalidInput("at least one statement is required")
	}

	for i, statement := range r.Statements {
		if err := statement.Validate(); err != nil {
			return apperror.InvalidInput(fmt.Sprintf("statement[%d]: %s", i, err.Error()))
		}
	}

	return nil

}

func (r *UpdatePolicy) Validate(resourceType string, resourceID string) error {
	if resourceType == "" {
		return apperror.InvalidInput("resource type is required")
	}

	if resourceID == "" {
		return apperror.InvalidInput("resource id is required")
	}

	if r.Version == "" {
		return apperror.InvalidInput("version is required")
	}

	if len(r.Statements) == 0 {
		return apperror.InvalidInput("at least one statement is required")
	}

	for i, statement := range r.Statements {
		if err := statement.Validate(); err != nil {
			return apperror.InvalidInput(fmt.Sprintf("statement[%d]: %s", i, err.Error()))
		}
	}

	return nil

}

func (s *Statement) Validate() error {
	if s.Effect == "" {
		return apperror.InvalidInput("effect is required")
	}

	if s.Effect != "ALLOW" && s.Effect != "DENY" {
		return apperror.InvalidInput(fmt.Sprintf("invalid effect: %s", s.Effect))
	}

	if len(s.Actions) == 0 {
		return apperror.InvalidInput("at least one action is required")
	}

	if len(s.Resources) == 0 {
		return apperror.InvalidInput("at least one resource is required")
	}

	if len(s.Principals) == 0 {
		return apperror.InvalidInput("at least one principal is required")
	}

	for i, principal := range s.Principals {
		if err := principal.Validate(); err != nil {
			return apperror.InvalidInput(fmt.Sprintf("principal[%d]: %s", i, err.Error()))
		}
	}

	return nil
}

func (p *Principal) Validate() error {
	if p.Type == "" {
		return apperror.InvalidInput("type is required")
	}

	if p.ID == "" {
		return apperror.InvalidInput("id is required")
	}

	return nil
}
