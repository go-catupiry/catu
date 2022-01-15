package pagination

import (
	"encoding/json"
)

type Pager struct {
	CurrentUrl string

	Page  int64
	Limit int64
	Count int64

	MaxLinks int64

	FirstPath   string
	FirstNumber string

	LastPath   string
	LastNumber string

	HasPrevius bool

	PreviusPath   string
	PreviusNumber string

	HasMoreBefore bool

	HasMoreAfter bool
	HasNext      bool

	NextPath   string
	NextNumber string

	Links []Link
}

type Link struct {
	Path     string
	Number   string
	IsActive bool
}

func (r *Pager) ToJSON() []byte {
	jsonString, _ := json.MarshalIndent(r, "", "  ")
	return jsonString
}

func NewPager() *Pager {
	var p Pager

	p.MaxLinks = 2
	p.Page = 1

	return &p
}
