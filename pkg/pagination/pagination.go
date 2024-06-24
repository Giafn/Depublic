package pkg

import "math"

type Pagination struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
	TotalCount int         `json:"totalCount"`
}

func Paginate(data interface{}, count int, page int, pageSize int) Pagination {
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	return Pagination{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: count,
	}
}
