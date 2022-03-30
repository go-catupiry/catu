package catu

import (
	"os"

	"github.com/go-catupiry/catu/configuration"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var appInstance *App

func init() {
	initDotEnvConfigSupport()
}

func Init() *App {
	appInstance = newApp()

	InitSanitizer()

	appInstance.RegisterPlugin(&Plugin{Name: "catu"})
	return appInstance
}

// Init doc env config with default development.env configuration file
// The configuration file pattern is: [environment].env
func initDotEnvConfigSupport() {
	env, _ := os.LookupEnv("GO_ENV")

	if env == "" {
		env = "development"
	}

	if _, err := os.Stat(env + ".env"); err == nil {
		godotenv.Load(env + ".env")
	}
}

func GetApp() *App {
	return appInstance
}

func GetConfiguration() configuration.Configer {
	return appInstance.Configuration
}

func GetDefaultDatabaseConnection() *gorm.DB {
	return appInstance.DB
}
