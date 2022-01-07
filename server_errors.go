package catu

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		code = 404
	}

	switch code {
	case 401:
		forbiddenErrorHandler(err, c)
	case 404:
		notFoundErrorHandler(err, c)
	case 500:
		internalServerErrorHandler(err, c)
	default:
		log.Println("customHTTPErrorHandler Echo error handler", err)
		errorPage := fmt.Sprintf("site/%d.html", code)
		logrus.WithFields(logrus.Fields{
			"errorPage":  errorPage,
			"statusCode": code,
			"error":      fmt.Sprintf("%+v\n", err),
		}).Warn("customHTTPErrorHandler unknow error status code")

		if err := c.File(errorPage); err != nil {
			c.Logger().Error(err)
		}
		c.Logger().Error(err)
	}
}

func forbiddenErrorHandler(err error, c echo.Context) error {
	ctx := c.Get("app").(*AppContext)

	switch ctx.ResponseContentType {
	case "application/json":
		c.JSON(http.StatusUnauthorized, make(map[string]string))
		return nil
	case "application/vnd.api+json":
		c.JSON(http.StatusUnauthorized, make(map[string]string))
		return nil
	default:
		ctx.Title = "Acesso restrito"

		if err := c.Render(http.StatusNotFound, "site/401", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}

		return nil
	}

}

func notFoundErrorHandler(err error, c echo.Context) error {
	ctx := c.Get("app").(*AppContext)

	switch ctx.ResponseContentType {
	case "application/json":
		c.JSON(http.StatusNotFound, make(map[string]string))
		return nil
	case "application/vnd.api+json":
		c.JSON(http.StatusNotFound, make(map[string]string))
		return nil
	default:
		ctx.Title = "Não encontrado"

		if err := c.Render(http.StatusNotFound, "site/404", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}

		return nil
	}

}

func internalServerErrorHandler(err error, c echo.Context) error {
	ctx := c.Get("app").(*AppContext)

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	logrus.WithFields(logrus.Fields{
		"err":  fmt.Sprintf("%+v\n", err),
		"code": code,
	}).Warn("internalServerErrorHandler error")

	switch ctx.ResponseContentType {
	case "application/json":
		c.JSON(http.StatusInternalServerError, make(map[string]string))
		return nil
	case "application/vnd.api+json":
		c.JSON(http.StatusInternalServerError, make(map[string]string))
		return nil
	default:
		ctx.Title = "Internal server error"

		if err := c.Render(http.StatusInternalServerError, "site/500", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}

		return nil
	}

}
