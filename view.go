package catu

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type tplWrapper struct {
	Ctx *AppContext
}

type TemplateCTX struct {
	Ctx     interface{}
	Record  interface{}
	Records interface{}
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	switch v := data.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		fmt.Printf("Integer: %v", v)
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		fmt.Printf("Float64: %v", v)
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		fmt.Printf("String: %v", v)
	default:
		htmlContext := data.(*TemplateCTX)

		var contentBuffer bytes.Buffer
		err := t.templates.ExecuteTemplate(&contentBuffer, name, htmlContext)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"name":  name,
			}).Error("theme.Render error on execute template")
			return err
		}

		ctx := htmlContext.Ctx.(*AppContext)
		ctx.Content = template.HTML(contentBuffer.String())

		var layoutBuffer bytes.Buffer
		err = t.templates.ExecuteTemplate(&layoutBuffer, ctx.Layout, htmlContext)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"name":  name,
			}).Error("theme.Render error on execute layout template")
			return err
		}

		ctx.Content = template.HTML(layoutBuffer.String())

		return t.templates.ExecuteTemplate(w, "site/html", htmlContext)
	}

	return nil
}
