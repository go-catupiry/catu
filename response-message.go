package catu

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ResponseMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func SetResponseMessage(c echo.Context, key string, message *ResponseMessage) error {
	messages, err := GetResponseMessages(c)
	if err != nil {
		return err
	}

	messages[key] = message
	c.Set("requestMessages", messages)

	return nil
}

func GetResponseMessages(c echo.Context) (map[string]*ResponseMessage, error) {
	iMessages := c.Get("requestMessages")

	switch ms := iMessages.(type) {
	case map[string]*ResponseMessage:
		return ms, nil
	}

	return map[string]*ResponseMessage{}, nil
}

func ResponseMessagesRender(c echo.Context, tpl string) template.HTML {
	html := ""

	messages, err := GetResponseMessages(c)
	if err != nil {
		return template.HTML(html)
	}

	if tpl == "" {
		tpl = "blocks/response/messages"
	}

	app := c.Get("app").(App)
	if app.HasTemplate(tpl) {
		var htmlBuffer bytes.Buffer
		err := app.RenderTemplateWithTheme(&htmlBuffer, GetTheme(c), tpl, messages)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":    fmt.Sprintf("%+v\n", err),
				"template": tpl,
			}).Error("ResponseMessageRender error on render template")

			return template.HTML(html)
		}

		html = htmlBuffer.String()
	}

	for _, m := range messages {
		html += `<div>` + m.Type + ": " + m.Message + `</div>'`
	}

	return template.HTML(html)
}
