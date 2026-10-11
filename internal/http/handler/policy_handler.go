package handler

import (
	"encoding/json"
	"net/http"
	apperror "rps/internal/errors"
	"rps/internal/http/helper"
	"rps/internal/http/httpx"
	request2 "rps/internal/http/request"
	"rps/internal/http/response"
	"rps/internal/service"

	"github.com/go-chi/chi/v5"
)

type PolicyHandler struct {
	service service.PolicyService
}

func NewPolicyHandler(service service.PolicyService) *PolicyHandler {
	return &PolicyHandler{
		service: service,
	}
}

func (h *PolicyHandler) GetPolicy(writer http.ResponseWriter, request *http.Request) {

	resourceType := request.PathValue("resourceType")
	resourceID := request.PathValue("resourceId")
	policy, err := h.service.GetPolicy(request.Context(), resourceType, resourceID)

	if err != nil {
		httpx.WriteError(writer, err)
		return
	}

	if err = httpx.WriteSuccess(writer, response.NewPolicyDetail(policy)); err != nil {
		httpx.WriteError(writer, err)
	}
}

func (h *PolicyHandler) ListPolicies(writer http.ResponseWriter, request *http.Request) {

	page, _ := helper.QueryInt(request, "page", 1)
	size, _ := helper.QueryInt(request, "size", 10)

	req := request2.NewListPolicies(int32(page), int32(size), helper.QueryParam(request, "resource_type"), helper.QueryParam(request, "resource_id"))

	policies, total, err := h.service.ListPolicies(request.Context(), req)
	if err != nil {
		httpx.WriteError(writer, err)
		return
	}

	policiesRes := make([]response.PolicyInfo, len(policies))
	for i, policy := range policies {
		policiesRes[i] = response.PolicyInfo{
			ResourceType: policy.ResourceType,
			ResourceID:   policy.ResourceID,
			Version:      policy.Version,
			CreatedAt:    policy.CreatedAt,
			UpdatedAt:    policy.UpdatedAt,
		}
	}

	if err = httpx.WriteSuccess(writer, response.PageResponse[response.PolicyInfo]{
		Items: policiesRes,
		Total: total,
		Page:  req.GetPage(),
		Size:  req.GetPageSize(),
	}); err != nil {
		httpx.WriteError(writer, err)
	}
}

func (h *PolicyHandler) CreatePolicy(writer http.ResponseWriter, request *http.Request) {

	var req request2.CreatePolicy

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		httpx.WriteError(writer, apperror.InvalidInput("invalid request body"))
		return
	}

	if err := req.Validate(); err != nil {
		httpx.WriteError(writer, err)
		return
	}

	if err := h.service.CreatePolicy(request.Context(), &req); err != nil {
		httpx.WriteError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func (h *PolicyHandler) UpdatePolicy(writer http.ResponseWriter, request *http.Request) {

	resourceType := request.PathValue("resourceType")
	resourceID := request.PathValue("resourceId")

	var req request2.UpdatePolicy
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		httpx.WriteError(writer, apperror.InvalidInput("invalid request body"))
		return
	}

	if err := req.Validate(resourceType, resourceID); err != nil {
		httpx.WriteError(writer, err)
		return
	}

	if err := h.service.UpdatePolicy(request.Context(), resourceType, resourceID, &req); err != nil {
		httpx.WriteError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *PolicyHandler) DeletePolicy(writer http.ResponseWriter, request *http.Request) {
	resourceType := chi.URLParam(request, "resourceType")
	resourceID := chi.URLParam(request, "resourceId")

	if resourceType == "" || resourceID == "" {
		httpx.WriteError(writer, apperror.InvalidInput("resourceType and resourceId are required"))
		return
	}

	if err := h.service.DeletePolicy(request.Context(), resourceType, resourceID); err != nil {
		httpx.WriteError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
