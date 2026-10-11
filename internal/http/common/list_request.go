package common

type ListRequest struct {
	page     int32
	pageSize int32
	sortBy   string
	order    string
}

func NewListRequest(page int32, pageSize int32, sortBy string, order string) *ListRequest {
	return &ListRequest{page: page, pageSize: pageSize, sortBy: sortBy, order: order}
}

func (r *ListRequest) GetPage() int32 {
	if r.page <= 0 {
		r.page = 1
	}
	return r.page
}

func (r *ListRequest) GetPageSize() int32 {
	if r.pageSize <= 0 {
		r.pageSize = 10
	}
	return r.pageSize
}

func (r *ListRequest) GetSortBy() string {
	return r.sortBy
}

func (r *ListRequest) GetOrder() string {
	return r.order
}

func (r *ListRequest) GetOffset() int32 {
	return (r.GetPage() - 1) * r.GetPageSize()
}

func (r *ListRequest) GetLimit() int32 {
	return r.GetPageSize()
}
