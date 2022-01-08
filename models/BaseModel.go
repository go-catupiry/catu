package models

import (
	"time"
)

type Base struct {
	ID        uint64    `gorm:"column:id;primary_key"  json:"id"`
	CreatedAt time.Time `gorm:"column:createdAt;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:datetime;not null" json:"updatedAt"`

	LinkPermanent string `gorm:"-" json:"linkPermanent"`
}
