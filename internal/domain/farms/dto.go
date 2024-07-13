package farms

import "github.com/nmluci/da-farm-be/internal/core/httpres"

// FarmRequestQuery represent query parameter fetch from request
type FarmRequestQuery struct {
	ID      int64  `param:"farmID" example:"1"`
	Keyword string `query:"keyword" example:"Farm"`
	Limit   uint64 `query:"limit" example:"100"`
	Page    uint64 `query:"page" example:"2"`
}

// FarmPayload represent payload fetch from request
type FarmPayload struct {
	ID   int64  `param:"farmID" example:"1" json:"-"` // ignore any value assigned via JSON body
	Name string `json:"name" example:"Farm A"`
}

// FarmResponse represent domain response for Farm entity
type FarmResponse struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Farm A"`
}

// ListFarmResponse represent domain response for bulk Farm entities
type ListFarmResponse struct {
	Farms []*FarmResponse        `json:"farms"`
	Meta  httpres.ListPagination `json:"meta"`
}
