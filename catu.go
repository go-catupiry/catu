package catu

import "gorm.io/gorm"

var appInstance *App

func Init() *App {
	appInstance = newApp()
	return appInstance
}

func GetApp() *App {
	return appInstance
}

func GetDefaultDatabaseConnection() *gorm.DB {
	return appInstance.DB
}
