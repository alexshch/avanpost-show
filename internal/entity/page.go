package entity

type PageInfo struct {
	StartIndex   int `json:"pageIndex"`
	ItemsPerPage int `json:"pageSize"`
	Total        int `json:"total"`
}

type PageParam struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}

func (p PageParam) GetOffset() int {
	return (p.PageIndex - 1) * p.PageSize
}

type PagedItemsList[T any] struct {
	StartIndex   int `json:"pageIndex"`
	ItemsPerPage int `json:"pageSize"`
	Total        int `json:"total"`
	Items        []T `json:"items"`
}

func NewPagedItemsList[T any](pageIndex, pageSize, total int, items []T) *PagedItemsList[T] {
	return &PagedItemsList[T]{
		StartIndex:   pageIndex,
		ItemsPerPage: pageSize,
		Total:        total,
		Items:        items,
	}
}
