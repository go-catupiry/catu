package catu

import (
	"os"
	"testing"

	"github.com/pkg/errors"
)

func TestMain(t *testing.T) {

}

func GetAppInstance() App {
	os.Setenv("DB_URI", "file::memory:?cache=shared")
	os.Setenv("DB_ENGINE", "sqlite")
	os.Setenv("TEMPLATE_FOLDER", "./_stubs/themes")

	app := Init(&AppOptions{})

	err := app.Bootstrap()
	if err != nil {
		panic(err)
	}

	err = app.GetDB().AutoMigrate()
	if err != nil {
		panic(errors.Wrap(err, "catu.GetAppInstance Error on run auto migration"))
	}

	return app
}
