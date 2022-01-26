package catu

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
