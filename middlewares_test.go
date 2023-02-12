package catu

// Tests example:
// func TestContentNegotiationMiddleware(t *testing.T) {
// 	assert := assert.New(t)

// 	GetTestAppInstance()

// 	url := "/symbol"
// 	e := echo.New()
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		t.Errorf("TestContentNegotiationMiddleware error: %v", err)
// 	}
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	ctx := NewRequestContext(&RequestContextOpts{EchoContext: c})
// 	c.Set("ctx", &ctx)

// 	t.Run("Should start with text/html", func(t *testing.T) {
// 		assert.Equal("text/html", ctx.GetResponseContentType())
// 	})

// 	assert.Equal("text/html", ctx.GetResponseContentType())

// 	f := func(c echo.Context) error {
// 		return nil
// 	}

// 	// TODO! add tests here
// }
