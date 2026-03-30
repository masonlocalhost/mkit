package dto

import (
	"mkit/example/ginapp/internal/model"
	"time"
)

var (
	ListTechnologiesOrderColumns = map[string]bool{
		model.Technology_CREATED_AT_COLUMN: true,
		model.Technology_UPDATED_AT_COLUMN: true,
	}
)

type Technology struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	Vendor       string    `json:"vendor,omitzero"`
	CPEType      string    `json:"cpe_type,omitzero"`
	ThumbnailUrl string    `json:"thumbnail_url,omitzero"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ListTechnologiesRequest struct {
	CPETypes      []string  `form:"cpe_types" binding:"omitempty,unique,dive,min=1"`
	Vendors       []string  `form:"vendors" binding:"omitempty,unique,dive,min=1"`
	Search        string    `form:"search"`
	CollectionIDs []string  `form:"collection_ids" binding:"omitempty,dive,uuid"`
	CreatedFrom   time.Time `form:"created_from"`
	CreatedTo     time.Time `form:"created_to"`

	PaginationParams
	SortingParams
}

type ListTechnologiesResponse []*Technology

type GetTechnologyRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type GetTechnologyResponse struct {
	Technology *Technology `json:"technology"`
}

type ListTechnologyVendorsRequest struct {
	CollectionIDs []string `form:"collection_ids" binding:"omitempty,dive,uuid"`
}

type ListTechnologyVendorsResponse []string

type ListTechnologyCPETypesRequest struct {
	CollectionIDs []string `form:"collection_ids" binding:"omitempty,dive,uuid"`
}

type ListTechnologyCPETypesResponse []string
