package catu

type Plugin struct {
	Name string
}

func (p *Plugin) Init(app *App) error {
	return nil
}

func (p *Plugin) BindRoutes(app interface{}) error {
	return nil
}
