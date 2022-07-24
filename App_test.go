package catu

var testAppInstance App

func GetTestAppInstance() App {
	if testAppInstance != nil {
		return testAppInstance
	}

	app := Init(&AppOptions{})

	err := app.Bootstrap()
	if err != nil {
		panic(err)
	}

	testAppInstance = app

	return app
}
