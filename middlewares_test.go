package catu

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestContentNegotiationMiddleware(t *testing.T) {
	assert := assert.New(t)

	GetTestAppInstance()

	url := "/symbol"
	e := echo.New()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Errorf("TestContentNegotiationMiddleware error: %v", err)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := NewRequestRequestContext(c)
	c.Set("app", &ctx)

	t.Run("Should start with text/html", func(t *testing.T) {
		assert.Equal("text/html", ctx.ResponseContentType)
	})

	assert.Equal("text/html", ctx.ResponseContentType)

	f := func(c echo.Context) error {
		return nil
	}

	t.Run("Should set application/json based in Accept header", func(t *testing.T) {
		c.Request().Header.Set("Accept", "application/json")
		middleware := contentNegotiationMiddleware()
		middleware(f)(c)
		assert.Equal("application/json", ctx.ResponseContentType)
	})

	t.Run("Should set application/vnd.api+json based in Accept header", func(t *testing.T) {
		// reset data:
		ctx.ResponseContentType = "text/html"
		c.Request().Header.Del("Content-Type")
		// mocked data:
		c.Request().Header.Set("Accept", "application/vnd.api+json")
		// run it:
		middleware := contentNegotiationMiddleware()
		middleware(f)(c)
		assert.Equal("application/vnd.api+json", ctx.ResponseContentType)
	})

	t.Run("Should set application/vnd.api+json based in Content type header", func(t *testing.T) {
		// reset data:
		ctx.ResponseContentType = "text/html"
		c.Request().Header.Del("Accept")
		// mock:
		c.Request().Header.Set("Content-Type", "application/vnd.api+json")
		middleware := contentNegotiationMiddleware()
		// run it:
		middleware(f)(c)
		assert.Equal("application/vnd.api+json", ctx.ResponseContentType)
	})
}
