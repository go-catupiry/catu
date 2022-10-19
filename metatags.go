package catu

import "time"

type HTMLMetaTags struct {
	Title       string
	Description string
	Canonical   string
	SiteName    string
	Type        string
	ImageURL    string
	ImageHeight string
	ImageWidth  string
	Author      string
	Keywords    string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	PublishedAt *time.Time
	TwitterSite string
}
