package models

import "time"

type Route struct {
	ID   uint `gorm:"primaryKey;column:id;type:int(11);not null" json:"-"`
	Path string
	Type string

	Handler *Handler

	CreatedAt time.Time `gorm:"column:createdAt;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:datetime;not null" json:"updatedAt"`
}

// TableName - Set db table name for Route table
func (r *Route) TableName() string {
	return "routes"
}

type Handler struct {
	ID         uint `gorm:"primaryKey;column:id;type:int(11);not null" json:"-"`
	Components string

	CreatedAt time.Time `gorm:"column:createdAt;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:datetime;not null" json:"updatedAt"`
}

// TableName - Set db table name for Route table
func (r *Handler) TableName() string {
	return "route_handlers"
}
