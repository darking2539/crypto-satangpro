package models

type Pagination struct {
	Page         int64 `json:"page"`
	PerPage      int64 `json:"perPage"`
	TotalResults int64 `json:"totalResults"`
	TotalPages   int64 `json:"totalPages"`
}
