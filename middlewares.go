package catu

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// BindMiddlewares - Bind middlewares in order
func BindMiddlewares(app App, p *Plugin) {
	logrus.Debug("catu.BindMiddlewares " + p.GetName())

	goEnv := app.GetConfiguration().Get("GO_ENV")

	router := app.GetRouter()

	router.Pre(preRequestMiddleware())
	router.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	router.Pre(extensionMiddleware())
	router.Pre(contentNegotiationMiddleware())

	router.Use(middleware.Gzip())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: app.GetConfiguration().GetBoolF("CORS_ALLOW_CREDENTIALS", true),
		MaxAge:           app.GetConfiguration().GetIntF("CORS_MAX_AGE", 18000), // seccounds
	}))
	router.Use(initAppCtx())

	if goEnv == "dev" {
		router.Debug = true
	}
}

func preRequestMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// default required values for init the request
			c.Set("responseContentType", "text/html")

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
				c.Set("responseContentType", "application/json")
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
			responseTypeI := c.Get("responseContentType")
			responseType := "text/html"

			if responseTypeI != nil {
				responseType = responseTypeI.(string)
			}

			if responseType != "text/html" {
				// already set...
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
					responseType = "application/vnd.api+json"
				}

				if responseType != "text/html" {
					logrus.WithFields(logrus.Fields{
						"contentTypeHeader":   contentTypeHeader,
						"ResponseContentType": responseType,
					}).Debug("contentNegotiationMiddleware found in contentTypeHeader header")

					c.Set("responseContentType", responseType)
					return next(c)
				}
			}

			switch acceptHeader {
			case "application/json":
				responseType = "application/json"
			case "application/vnd.api+json":
				responseType = "application/vnd.api+json"
			}

			c.Set("responseContentType", responseType)
			return next(c)
		}
	}
}

// Middleare that update echo context to use custom methods
func initAppCtx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := NewRequestContext(&RequestContextOpts{EchoContext: c})

			logrus.WithFields(logrus.Fields{
				"URL":    ctx.Request().URL,
				"method": ctx.Request().Method,
				"header": ctx.Request().Header,
			}).Debug("init catu ctx init")

			return next(ctx)
		}
	}
}
