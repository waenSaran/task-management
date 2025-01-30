package models

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
	Data       []interface{}
}

func TransformPagination(pagination *Pagination) map[string]interface{} {
	return map[string]interface{}{
		"page":       pagination.Page,
		"pageSize":   pagination.PageSize,
		"total":      pagination.Total,
		"totalPages": pagination.TotalPages,
		"data":       pagination.Data,
	}
}
