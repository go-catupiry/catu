package catu

import "github.com/labstack/echo/v4"

type CustomBinder struct{}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	// You may use default binder
	db := &echo.DefaultBinder{}
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return
	}

	return
}
