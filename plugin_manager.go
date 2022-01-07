package catu

import "github.com/labstack/echo/v4"

type PluginManager interface {
	GetRouter(name string) *echo.Group
	SetRouter(name string) *echo.Group
	GetAPIRouter(name string) *echo.Group
}
