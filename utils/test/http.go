package utils_test

// import (
// 	"net/http"
// 	"net/http/httptest"

// 	"github.com/go-catupiry/catu"
// 	"github.com/labstack/echo/v4"
// )

// func GetRequestContext(url string) (*echo.Echo, echo.Context, *httptest.ResponseRecorder, error) {
// 	e := echo.New()
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	ctx := catu.GetRequestCtx(c)
// 	c.Set("app", &ctx)

// 	return e, c, rec, nil
// }
