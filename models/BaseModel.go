package models

import (
	"time"
)

type Base struct {
	ID        uint64    `gorm:"column:id;primary_key"  json:"id"`
	CreatedAt time.Time `gorm:"column:createdAt;type:datetime;not null" json:"createdAt" jsonapi:"attr,updatedAt,iso8601"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:datetime;not null" json:"updatedAt" jsonapi:"attr,updatedAt,iso8601"`

	LinkPermanent string `gorm:"-" json:"linkPermanent"`
}
