package dto

type PaginationParams struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit" binding:"required"`
}

type SortingParams struct {
	SortBy    string `form:"sort_by" json:"sort_by"`
	SortOrder string `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc"`
}

type PermissionsByUserIDMap struct {
	PermissionsByUserID map[int][]string `json:"permissions_by_user_id" binding:"required,dive,keys,gt=0,endkeys,dive,oneof=edit_collection_info check_resource_validation edit_collection_users"`
}
