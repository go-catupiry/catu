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
	// router.Use(mw.Recover())
	router.Use(middleware.Gzip())

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		MaxAge:           18000, // seccounds
	}))

	if goEnv == "dev" {
		router.Debug = true
	}

	router.Pre(initAppCtx())
	router.Pre(extensionMiddleware())
	router.Pre(contentNegotiationMiddleware())
	// router.Pre(urlAliasMiddleware())

	// app.GetRouterGroup("main").Use(oauth2AuthenticationMiddleware())
	// app.GetRouterGroup("main").Use(sessionAuthenticationMiddleware())
	// app.GetRouterGroup("api").Use(oauth2AuthenticationMiddleware())
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

// // urlAliasMiddleware - Change url to handle aliased urls like /about to /content/1
// func urlAliasMiddleware() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			ctx := c.Get("app").(*AppContext)
// 			ctx.PathBeforeAlias = c.Request().URL.Path

// 			if configuration.CFGs.URL_ALIAS_ENABLE == "" {
// 				return next(c)
// 			}

// 			path, err := getPathFromReq(c.Request())
// 			if err != nil {
// 				return c.String(http.StatusInternalServerError, "Error on parse url")
// 			}

// 			if isPublicRoute(path) {
// 				// public folders dont have alias...
// 				return next(c)
// 			}

// 			logrus.WithFields(logrus.Fields{
// 				"url":           path,
// 				"c.path":        c.Path(),
// 				"c.QueryString": c.QueryString(),
// 			}).Debug("urlAliasMiddleware url after alias")

// 			// save path before alias for reuse ...
// 			ctx.PathBeforeAlias = path

// 			var record models.UrlAlias
// 			err = models.UrlAliasGetByURL(path, &record)
// 			if err != nil {
// 				log.Println("Error on get url alias", err)
// 			}

// 			if record.Target != "" && record.Alias != "" {
// 				if record.Target == path && ctx.ResponseContentType == "text/html" {
// 					// redirect to alias url keeping the query string
// 					queryString := c.QueryString()
// 					if queryString != "" {
// 						queryString = "?" + queryString
// 					}

// 					return c.Redirect(302, record.Alias+queryString)
// 				} else {
// 					// override and continue with target url
// 					RewriteURL(record.Target, c)
// 					ctx.Pager.CurrentUrl = path
// 				}
// 			}

// 			return next(c)
// 		}
// 	}
// }

// func getUrlFromReq(req *http.Request) (string, error) {
// 	rawURI := req.RequestURI
// 	if rawURI != "" && rawURI[0] != '/' {
// 		prefix := ""
// 		if req.URL.Scheme != "" {
// 			prefix = req.URL.Scheme + "://"
// 		}
// 		if req.URL.Host != "" {
// 			prefix += req.URL.Host // host or host:port
// 		}
// 		if prefix != "" {
// 			rawURI = strings.TrimPrefix(rawURI, prefix)
// 		}
// 	}

// 	return rawURI, nil
// }

// func getPathFromReq(req *http.Request) (string, error) {
// 	return req.URL.Path, nil
// }

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
			} else {
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
