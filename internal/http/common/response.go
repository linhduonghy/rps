package common

type PageResponse[T any] struct {
	Items []T `json:"items"`

	Total int64 `json:"total"`

	Page int `json:"page"`

	PageSize int `json:"pageSize"`
}
