package catu

import "github.com/gookit/event"

type Plugin struct {
	Name string
}

func (p *Plugin) Init(app *App) error {
	app.Events.On("bindMiddlewares", event.ListenerFunc(func(e event.Event) error {
		return p.BindMiddlewares(app)
	}), event.Normal)

	return nil
}

func (p *Plugin) GetName() string {
	return p.Name
}

func (p *Plugin) BindMiddlewares(app *App) error {
	return nil
}

func (p *Plugin) SetTemplateFuncMap(app *App) error {
	return nil
}
