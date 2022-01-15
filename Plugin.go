package catu

import (
	"github.com/gookit/event"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	Name string
}

func (p *Plugin) Init(app *App) error {
	logrus.WithFields(logrus.Fields{
		"PluginName": p.Name,
	}).Debug("tupi.Plugin.Init Running init")

	app.Events.On("bindMiddlewares", event.ListenerFunc(func(e event.Event) error {
		return p.BindMiddlewares(app)
	}), event.Normal)

	app.Events.On("setTemplateFunctions", event.ListenerFunc(func(e event.Event) error {
		return p.setTemplateFunctions(app)
	}), event.Normal)

	return nil
}

func (p *Plugin) GetName() string {
	return p.Name
}

func (p *Plugin) BindMiddlewares(app *App) error {
	BindMiddlewares(app, p)
	return nil
}

func (p *Plugin) setTemplateFunctions(app *App) error {
	app.SetTemplateFunction("paginate", paginate)
	app.SetTemplateFunction("contentDates", contentDates)
	app.SetTemplateFunction("shareMenu", shareMenu)
	app.SetTemplateFunction("truncate", truncate)

	return nil
}

func (p *Plugin) SetTemplateFuncMap(app *App) error {
	return nil
}
