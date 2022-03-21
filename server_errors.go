package catu

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ValidationResponse struct {
	Errors []*ValidationFieldError `json:"errors"`
}

type ValidationFieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type HTTPError struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

// func NewHTTPError(code, message string) *HTTPError {
// 	return HTTPError{ Code: code, Message: message }
// }

// type NotFoundResponseError

func CustomHTTPErrorHandler(err error, c echo.Context) {
	logrus.WithFields(logrus.Fields{
		"err": fmt.Sprintf("%+v\n", err),
	}).Debug("catu.CustomHTTPErrorHandler running")

	ctx := c.Get("app").(*RequestContext)

	code := 0
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if ctx.ResponseContentType == "application/json" {
			c.JSON(he.Code, &HTTPError{Code: he.Code, Message: he.Message})
			return
		}
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		validationError(ve, err, c)
		return
	}

	if code == 0 && err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		code = 404
	}

	switch code {
	case 401:
		unAuthorizedErrorHandler(err, c)
	case 403:
		forbiddenErrorHandler(err, c)
	case 404:
		notFoundErrorHandler(err, c)
	case 500:
		internalServerErrorHandler(err, c)
	default:
		logrus.WithFields(logrus.Fields{
			"error":             fmt.Sprintf("%+v\n", err),
			"statusCode":        code,
			"path":              c.Path(),
			"method":            c.Request().Method,
			"AuthenticatedUser": ctx.AuthenticatedUser,
			"roles":             ctx.GetAuthenticatedRoles(),
		}).Warn("customHTTPErrorHandler unknown error status code")
		c.JSON(http.StatusInternalServerError, &HTTPError{Code: 500, Message: "Unknown Error"})
	}
}

func forbiddenErrorHandler(err error, c echo.Context) error {
	ctx := c.Get("app").(*RequestContext)

	logrus.WithFields(logrus.Fields{
		"err":               fmt.Sprintf("%+v\n", err),
		"code":              "403",
		"path":              c.Path(),
		"method":            c.Request().Method,
		"AuthenticatedUser": ctx.AuthenticatedUser,
		"roles":             ctx.GetAuthenticatedRoles(),
	}).Debug("catu.forbiddenErrorHandler running")

	switch ctx.ResponseContentType {
	case "text/html":
		ctx.Title = "Acesso restrito"

		if err := c.Render(http.StatusForbidden, "site/403", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}

		return nil
	default:
		c.JSON(http.StatusForbidden, err)
		return nil
	}
}

func unAuthorizedErrorHandler(err error, c echo.Context) error {
	ctx := c.Get("app").(*RequestContext)

	logrus.WithFields(logrus.Fields{
		"err":               fmt.Sprintf("%+v\n", err),
		"code":              "401",
		"path":              c.Path(),
		"method":            c.Request().Method,
		"AuthenticatedUser": ctx.AuthenticatedUser,
		"roles":             ctx.GetAuthenticatedRoles(),
	}).Info("catu.unAuthorizedErrorHandler running")

	switch ctx.ResponseContentType {
	case "text/html":
		ctx.Title = "Forbidden"

		if err := c.Render(http.StatusUnauthorized, "site/401", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}

		return nil
	default:
		c.JSON(http.StatusUnauthorized, err)
		return nil
	}

}

func notFoundErrorHandler(err error, c echo.Context) error {
	logrus.WithFields(logrus.Fields{
		"err":  fmt.Sprintf("%+v\n", err),
		"code": "404",
	}).Debug("catu.notFoundErrorHandler running")

	ctx := c.Get("app").(*RequestContext)

	switch ctx.ResponseContentType {
	case "text/html":
		ctx.Title = "Não encontrado"

		if err := c.Render(http.StatusNotFound, "site/404", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}
		return nil
	default:
		c.JSON(http.StatusNotFound, &HTTPError{Code: http.StatusNotFound, Message: "Not Found"})
		return nil
	}
}

func validationError(ve validator.ValidationErrors, err error, c echo.Context) error {
	logrus.WithFields(logrus.Fields{
		"err":  fmt.Sprintf("%+v\n", err),
		"code": "400",
	}).Debug("catu.validationError running")

	ctx := c.Get("app").(*RequestContext)

	resp := ValidationResponse{}

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el ValidationFieldError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			el.Message = err.Error()
			resp.Errors = append(resp.Errors, &el)
		}
	}

	switch ctx.ResponseContentType {
	case "text/html":
		ctx.Title = "Bad request"

		if err := c.Render(http.StatusInternalServerError, "site/400", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}

		return nil
	default:
		return c.JSON(http.StatusBadRequest, resp)
	}
}

func internalServerErrorHandler(err error, c echo.Context) error {
	ctx := c.Get("app").(*RequestContext)

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	logrus.WithFields(logrus.Fields{
		"err":               fmt.Sprintf("%+v\n", err),
		"code":              code,
		"path":              c.Path(),
		"method":            c.Request().Method,
		"AuthenticatedUser": ctx.AuthenticatedUser,
		"roles":             ctx.GetAuthenticatedRoles(),
	}).Warn("internalServerErrorHandler error")

	switch ctx.ResponseContentType {
	case "text/html":
		ctx.Title = "Internal server error"

		if err := c.Render(http.StatusInternalServerError, "site/500", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			c.Logger().Error(err)
		}

		return nil
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			return c.JSON(he.Code, &HTTPError{Code: he.Code, Message: he.Message})
		}

		c.JSON(http.StatusInternalServerError, &HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		return nil
	}
}
