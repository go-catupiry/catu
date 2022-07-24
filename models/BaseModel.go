package models

import (
	"time"
)

type Base struct {
	ID        uint64    `gorm:"column:id;primary_key"  json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt" jsonapi:"attr,updatedAt,iso8601"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt" jsonapi:"attr,updatedAt,iso8601"`

	LinkPermanent string `gorm:"-" json:"linkPermanent"`
}
