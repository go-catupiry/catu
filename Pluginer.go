package catu

type Pluginer interface {
	Init(app App) error
	GetName() string
}
