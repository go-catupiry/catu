package catu

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// BindMiddlewares - Bind middlewares in order
func BindMiddlewares(app *App, p *Plugin) {
	logrus.Debug("catu.BindMiddlewares " + p.GetName())

	goEnv := app.Configuration.Get("GO_ENV")

	router := app.GetRouter()

	router.Pre(initAppCtx())
	router.Pre(extensionMiddleware())
	router.Pre(contentNegotiationMiddleware())

	// router.Use(mw.Recover())
	router.Use(middleware.Gzip())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		MaxAge:           18000, // seccounds
	}))

	if goEnv == "dev" {
		router.Debug = true
	}

}

func initAppCtx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := NewRequestAppContext(c)
			c.Set("app", &ctx)

			logrus.WithFields(logrus.Fields{
				"URL":    c.Request().URL,
				"method": c.Request().Method,
				"header": c.Request().Header,
			}).Debug("initAppCtx init")

			return next(c)
		}
	}
}

func isPublicRoute(url string) bool {
	return strings.HasPrefix(url, "/health") || strings.HasPrefix(url, "/public")
}

// extensionMiddleware - handle url extensions
func extensionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			if isPublicRoute(req.URL.String()) {
				return next(c)
			}

			ctx := c.Get("app").(*AppContext)
			oldUrl := req.URL.Path

			logrus.WithFields(logrus.Fields{
				"url":    req.URL,
				"oldUrl": oldUrl,
			}).Debug("extensionMiddleware url before process")

			haveQueryParamJSONType := false

			query := c.QueryString()
			if query != "" {
				query = "?" + query
				responseFormat := c.QueryParam("responseType")
				haveQueryParamJSONType = responseFormat == "json"
			}

			if strings.HasSuffix(oldUrl, ".json") || haveQueryParamJSONType {
				newPath := strings.TrimSuffix(oldUrl, ".json")
				url, err := req.URL.Parse(newPath + query)
				if err != nil {
					return err
				}

				req.URL = url
				ctx.ResponseContentType = "application/json"
			}

			logrus.WithFields(logrus.Fields{
				"url": req.URL,
			}).Debug("extensionMiddleware url after process")

			return next(c)
		}
	}
}

func contentNegotiationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Get("app").(*AppContext)
			if ctx.ResponseContentType != "text/html" {
				// already filled
				return next(c)
			}

			contentTypeHeader := c.Request().Header.Get("Content-Type")
			acceptHeader := c.Request().Header.Get("Accept")

			logrus.WithFields(logrus.Fields{
				"acceptHeader":      acceptHeader,
				"contentTypeHeader": contentTypeHeader,
			}).Debug("contentNegotiationMiddleware headers")

			if contentTypeHeader != "" {
				switch contentTypeHeader {
				case "application/vnd.api+json":
					ctx.ResponseContentType = "application/vnd.api+json"
				}

				if ctx.ResponseContentType != "text/html" {
					logrus.WithFields(logrus.Fields{
						"contentTypeHeader":   contentTypeHeader,
						"ResponseContentType": ctx.ResponseContentType,
					}).Debug("contentNegotiationMiddleware found in contentTypeHeader header")

					return next(c)
				}
			}

			switch acceptHeader {
			case "application/json":
				ctx.ResponseContentType = "application/json"
			case "application/vnd.api+json":
				ctx.ResponseContentType = "application/vnd.api+json"
			}

			return next(c)
		}
	}
}
