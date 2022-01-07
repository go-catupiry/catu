package catu

import (
	"encoding/json"
	"html/template"
	"time"

	"github.com/go-catupiry/catu/configuration"
	"github.com/go-catupiry/catu/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type App struct {
	InitTime time.Time

	Configuration configuration.Configer

	Plugins map[string]Pluginer

	router *echo.Echo

	routerGroups    map[string]*echo.Group
	apiRouterGroups map[string]*echo.Group

	RolesString string
	RolesList   map[string]Role

	templates *template.Template
}

func (r *App) RegisterPlugin(name string, p Pluginer) {
	r.Plugins[name] = p
}

func (r *App) GetRouter() *echo.Echo {
	return r.router
}

func (r *App) GetTemplates() *template.Template {
	return r.templates
}

// -- Plugin lifecircle methods:

func (r *App) BeforeBindMiddlewares() {
	for i := range r.Plugins {
		r.Plugins[i].BeforeBindMiddlewares(r)
	}
}

func (r *App) BindMiddlewares() {
	for i := range r.Plugins {
		r.Plugins[i].BindMiddlewares(r)
	}
}

func (r *App) AfterBindMiddlewares() {
	for i := range r.Plugins {
		r.Plugins[i].AfterBindMiddlewares(r)
	}
}

func (r *App) BeforeBindRoutes() {
	for i := range r.Plugins {
		r.Plugins[i].BeforeBindRoutes(r)
	}
}

func (r *App) BindRoutes() {
	for i := range r.Plugins {
		r.Plugins[i].BindRoutes(r)
	}
}

func (r *App) AfterBindRoutes() {
	for i := range r.Plugins {
		r.Plugins[i].AfterBindRoutes(r)
	}
}

func (r *App) Bootstrap() error {
	json.Unmarshal([]byte(r.Configuration.Get("Roles")), &r.RolesList)

	return nil
}

func (r *App) StartHTTPServer() error {

	return nil
}

func (r *App) SetRouterGroup(name, path string) *echo.Group {
	if r.routerGroups[name] == nil {
		r.routerGroups[name] = r.router.Group(path)
	}
	return r.routerGroups[name]
}

func (r *App) SetAPIRouterGroup(name, path string) *echo.Group {
	if r.apiRouterGroups[name] == nil {
		r.apiRouterGroups[name] = r.routerGroups["api"].Group(path)
	}
	return r.apiRouterGroups[name]
}

func NewApp() *App {
	var app App

	app.RolesString = configuration.Roles

	app.Configuration = configuration.NewCfg()
	app.routerGroups = make(map[string]*echo.Group)
	app.apiRouterGroups = make(map[string]*echo.Group)

	app.router = echo.New()
	app.Plugins = make(map[string]Pluginer)

	app.SetRouterGroup("main", "/")
	app.SetRouterGroup("public", "/public")

	apiRouterGroup := app.SetRouterGroup("api", "/api")
	apiRouterGroup.GET("", HealthCheck)

	app.router.Validator = &utils.CustomValidator{Validator: validator.New()}

	app.router.Renderer = &TemplateRenderer{
		templates: app.GetTemplates(),
	}

	app.router.HTTPErrorHandler = CustomHTTPErrorHandler

	return &app
}

type Apper interface {
	GetApp() *App
	RegisterPlugin(p interface{}) error
	Init() error
	BindRoutes() error
}
