package service

import (
	"context"
	"rps/internal/domain"
	"rps/internal/http/request"
)

type PolicyService interface {
	CreatePolicy(ctx context.Context, request *request.CreatePolicy) error

	UpdatePolicy(ctx context.Context, resourceType string, resourceID string, r *request.UpdatePolicy) error

	GetPolicy(
		ctx context.Context,
		resourceType string,
		resourceID string,
	) (*domain.Policy, error)

	ListPolicies(ctx context.Context, req *request.ListPolicies) ([]*domain.Policy, int64, error)

	DeletePolicy(
		ctx context.Context,
		resourceType string,
		resourceID string,
	) error
}
