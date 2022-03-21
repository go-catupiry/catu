package catu

import "github.com/labstack/echo/v4"

type HTTPController interface {
	Query(c echo.Context) error
	Create(c echo.Context) error
	Count(c echo.Context) error
	FindOne(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}
