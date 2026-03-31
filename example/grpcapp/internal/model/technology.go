package model

import "time"

type Technology struct {
	ID           string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"column:name;type:varchar(100);not null;uniqueIndex:uk_technology_name_version" json:"name"`
	Version      string    `gorm:"column:version;type:varchar(100);not null;uniqueIndex:uk_technology_name_version" json:"version"`
	Vendor       string    `gorm:"column:vendor;type:varchar(100)" json:"vendor"`
	CPEType      string    `gorm:"column:cpe_type;type:varchar(250)" json:"cpe_type"`
	ThumbnailUrl string    `gorm:"column:thumbnail_url;type:varchar(250)" json:"thumbnail_url"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

const (
	TechnologyTableName = "technologies"

	Technology_ID_COLUMN         = "id"
	Technology_NAME_COLUMN       = "name"
	Technology_VERSION_COLUMN    = "version"
	Technology_CREATED_AT_COLUMN = "created_at"
	Technology_VENDOR_COLUMN     = "vendor"
	Technology_CPE_TYPE_COLUMN   = "cpe_type"
	Technology_UPDATED_AT_COLUMN = "updated_at"
)
