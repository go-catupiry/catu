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

type HTTPErrorInterface interface {
	Error() string
	GetCode() int
	SetCode(code int) error
	GetMessage() interface{}
	SetMessage(message interface{}) error
}

// HTTPError implements HTTP Error interface, default error object
type HTTPError struct {
	Code     int         `json:"code"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}

// Error makes it compatible with `error` interface.
func (e *HTTPError) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", e.Code, e.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", e.Code, e.Message, e.Internal)
}

func (e *HTTPError) GetCode() int {
	return e.Code
}

func (e *HTTPError) SetCode(code int) error {
	e.Code = code
	return nil
}

func (e *HTTPError) GetMessage() interface{} {
	return e.Message
}

func (e *HTTPError) SetMessage(message interface{}) error {
	e.Message = message
	return nil
}

func (e *HTTPError) GetInternal() error {
	return e.Internal
}

func (e *HTTPError) SetInternal(internal error) error {
	e.Internal = internal
	return nil
}

type ValidationResponse struct {
	Errors []*ValidationFieldError `json:"errors"`
}

type ValidationFieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	logrus.WithFields(logrus.Fields{
		"err": fmt.Sprintf("%+v\n", err),
	}).Debug("catu.CustomHTTPErrorHandler running")

	var ctx *RequestContext

	switch v := c.(type) {
	case *RequestContext:
		ctx = v
	default:
		ctx = NewRequestContext(&RequestContextOpts{EchoContext: c})
	}

	code := 0
	if he, ok := err.(HTTPErrorInterface); ok {
		code = he.GetCode()
		if ctx.GetResponseContentType() == "application/json" {
			c.JSON(code, he)
			return
		}
	}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if ctx.GetResponseContentType() == "application/json" {
			c.JSON(code, he)
			return
		}
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		validationError(ve, err, ctx)
		return
	}

	if code == 0 && err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		code = 404
	}

	switch code {
	case 401:
		unAuthorizedErrorHandler(err, ctx)
	case 403:
		forbiddenErrorHandler(err, ctx)
	case 404:
		notFoundErrorHandler(err, ctx)
	case 500:
		internalServerErrorHandler(err, ctx)
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
	ctx := c.Get("ctx").(*RequestContext)

	logrus.WithFields(logrus.Fields{
		"err":               fmt.Sprintf("%+v\n", err),
		"code":              "403",
		"path":              c.Path(),
		"method":            c.Request().Method,
		"AuthenticatedUser": ctx.AuthenticatedUser,
		"roles":             ctx.GetAuthenticatedRoles(),
	}).Debug("catu.forbiddenErrorHandler running")

	switch ctx.GetResponseContentType() {
	case "text/html":
		ctx.Title = "Acesso restrito"

		if err := c.Render(http.StatusForbidden, "403", &TemplateCTX{
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

func unAuthorizedErrorHandler(err error, ctx *RequestContext) error {
	logrus.WithFields(logrus.Fields{
		"err":               fmt.Sprintf("%+v\n", err),
		"code":              "401",
		"path":              ctx.Path(),
		"method":            ctx.Request().Method,
		"AuthenticatedUser": ctx.AuthenticatedUser,
		"roles":             ctx.GetAuthenticatedRoles(),
	}).Info("catu.unAuthorizedErrorHandler running")

	switch ctx.GetResponseContentType() {
	case "text/html":
		ctx.Title = "Forbidden"

		if err := ctx.Render(http.StatusUnauthorized, "401", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			ctx.Logger().Error(err)
		}

		return nil
	default:
		ctx.JSON(http.StatusUnauthorized, err)
		return nil
	}

}

func notFoundErrorHandler(err error, ctx *RequestContext) error {
	logrus.WithFields(logrus.Fields{
		"err":  fmt.Sprintf("%+v\n", err),
		"code": "404",
	}).Debug("catu.notFoundErrorHandler running")

	switch ctx.GetResponseContentType() {
	case "text/html":
		ctx.Title = "Não encontrado"

		if err := ctx.Render(http.StatusNotFound, "404", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			ctx.Logger().Error(err)
		}
		return nil
	default:
		ctx.JSON(http.StatusNotFound, &HTTPError{Code: http.StatusNotFound, Message: "Not Found"})
		return nil
	}
}

func validationError(ve validator.ValidationErrors, err error, ctx *RequestContext) error {
	logrus.WithFields(logrus.Fields{
		"err":  fmt.Sprintf("%+v\n", err),
		"code": "400",
	}).Debug("catu.validationError running")

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

	switch ctx.GetResponseContentType() {
	case "text/html":
		ctx.Title = "Bad request"

		if err := ctx.Render(http.StatusInternalServerError, "400", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			ctx.Logger().Error(err)
		}

		return nil
	default:
		return ctx.JSON(http.StatusBadRequest, resp)
	}
}

func internalServerErrorHandler(err error, ctx *RequestContext) error {
	code := http.StatusInternalServerError
	if he, ok := err.(*HTTPError); ok {
		code = he.Code
	}

	logrus.WithFields(logrus.Fields{
		"err":               fmt.Sprintf("%+v\n", err),
		"code":              code,
		"path":              ctx.Path(),
		"method":            ctx.Request().Method,
		"AuthenticatedUser": ctx.AuthenticatedUser,
		"roles":             ctx.GetAuthenticatedRoles(),
	}).Warn("internalServerErrorHandler error")

	switch ctx.GetResponseContentType() {
	case "text/html":
		ctx.Title = "Internal server error"

		if err := ctx.Render(http.StatusInternalServerError, "500", &TemplateCTX{
			Ctx: ctx,
		}); err != nil {
			ctx.Logger().Error(err)
		}

		return nil
	default:
		if he, ok := err.(*HTTPError); ok {
			return ctx.JSON(he.Code, &HTTPError{Code: he.Code, Message: he.Message})
		}

		ctx.JSON(http.StatusInternalServerError, &HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		return nil
	}
}
