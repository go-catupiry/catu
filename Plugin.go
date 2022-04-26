package catu

import (
	"github.com/gookit/event"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	Name string
}

func (p *Plugin) Init(a App) error {
	logrus.WithFields(logrus.Fields{
		"PluginName": p.Name,
	}).Debug("catu.Plugin.Init Running init")

	a.GetEvents().On("bindMiddlewares", event.ListenerFunc(func(e event.Event) error {
		return p.BindMiddlewares(a)
	}), event.High)

	a.GetEvents().On("setTemplateFunctions", event.ListenerFunc(func(e event.Event) error {
		return p.setTemplateFunctions(a)
	}), event.Normal)

	return nil
}

func (p *Plugin) GetName() string {
	return p.Name
}

func (p *Plugin) BindMiddlewares(a App) error {
	BindMiddlewares(a, p)
	return nil
}

func (p *Plugin) setTemplateFunctions(app App) error {
	app.SetTemplateFunction("paginate", paginate)
	app.SetTemplateFunction("contentDates", contentDates)
	app.SetTemplateFunction("truncate", truncate)
	app.SetTemplateFunction("formatCurrency", formatCurrency)
	app.SetTemplateFunction("formatDecimalWithDots", formatDecimalWithDots)
	app.SetTemplateFunction("html", noEscapeHTML)

	return nil
}

func (p *Plugin) SetTemplateFuncMap(app App) error {
	return nil
}
