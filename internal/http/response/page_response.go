package response

type PageResponse[T any] struct {
	Items []T `json:"items"`

	Total int64 `json:"total"`

	Page int32 `json:"page"`

	Size int32 `json:"size"`
}
