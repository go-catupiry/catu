package helpers

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func RewriteURL(newPath string, c echo.Context) error {
	req := c.Request()

	rawURI := req.RequestURI
	if rawURI != "" && rawURI[0] != '/' {
		prefix := ""
		if req.URL.Scheme != "" {
			prefix = req.URL.Scheme + "://"
		}
		if req.URL.Host != "" {
			prefix += req.URL.Host // host or host:port
		}
		if prefix != "" {
			rawURI = strings.TrimPrefix(rawURI, prefix)
		}
	}

	query := c.QueryString()

	if query != "" {
		query = "?" + query
	}

	url, err := req.URL.Parse(newPath + query)
	if err != nil {
		return err
	}

	req.URL = url

	return nil
}
