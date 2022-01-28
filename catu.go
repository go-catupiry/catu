package catu

import (
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

func initDotEnvConfigSupport() {
	godotenv.Load()
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
