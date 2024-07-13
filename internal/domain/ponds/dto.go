package ponds

import "github.com/nmluci/da-farm-be/internal/core/httpres"

// PondRequestQuery represent query parameters fetch from request
type PondRequestQuery struct {
	ID      int64  `param:"pondID" example:"1"`
	FarmID  int64  `param:"farmID" example:"1"`
	Keyword string `query:"keyword" example:"Pond A"`
	Limit   uint64 `query:"limit" example:"100"`
	Page    uint64 `query:"page" example:"2"`
}

// PondPayload represent payload fetch from request body
type PondPayload struct {
	ID     int64  `param:"pondID" json:"-" example:"1"`
	FarmID int64  `param:"farmID" json:"-" example:"1"`
	Name   string `json:"name" example:"Pond 1"`
}

// PondResponse represent domain response for Pond entity
type PondResponse struct {
	ID       int64  `json:"id" example:"1"`
	FarmID   int64  `json:"farm_id" example:"1"`
	FarmName string `json:"farm_name" example:"Farm A"`
	Name     string `json:"pond_name" example:"Pond A"`
}

// ListPondResponse represent domain response for bulk Pond entities
type ListPondResponse struct {
	Ponds []*PondResponse        `json:"ponds"`
	Meta  httpres.ListPagination `json:"meta"`
}
