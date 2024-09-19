package types

type PaginatedResult struct {
    Data interface{} `json:"data"`
    TotalCount int `json:"totalCount"`
}

