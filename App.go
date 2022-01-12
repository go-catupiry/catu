package catu

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-catupiry/catu/configuration"
	"github.com/go-catupiry/catu/helpers"
	"github.com/go-catupiry/catu/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/event"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	gorm_logger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

type App struct {
	InitTime time.Time

	Events *event.Manager

	Configuration configuration.Configer

	DB *gorm.DB

	Plugins map[string]Pluginer

	router *echo.Echo

	routerGroups    map[string]*echo.Group
	apiRouterGroups map[string]*echo.Group

	RolesString string
	RolesList   map[string]Role

	templates *template.Template
}

func (r *App) RegisterPlugin(p Pluginer) {
	if p.GetName() == "" {
		panic("Plugin.RegisterPlugin Name should be returned from GetName method")
	}

	r.Plugins[p.GetName()] = p
}

func (r *App) GetRouter() *echo.Echo {
	return r.router
}

func (r *App) GetTemplates() *template.Template {
	return r.templates
}

func (r *App) Bootstrap() error {
	var err error

	logrus.Debug("Bootstrap running")
	// default roles and permissions, override it on your app
	json.Unmarshal([]byte(r.Configuration.Get("Roles")), &r.RolesList)

	for _, p := range r.Plugins {
		err = p.Init(r)
		if err != nil {
			return errors.Wrap(err, "App.Bootstrap | Error on run plugin init "+p.GetName())
		}
	}

	r.Events.MustTrigger("configuration", event.M{"app": r})
	r.Events.MustTrigger("bindMiddlewares", event.M{"app": r})
	r.Events.MustTrigger("bindRoutes", event.M{"app": r})
	r.Events.MustTrigger("setResponseFormats", event.M{"app": r})
	r.Events.MustTrigger("bootstrap", event.M{"app": r})

	return nil
}

func (r *App) StartHTTPServer() error {
	port := r.Configuration.Get("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Info("Server listening on port " + port)
	http.ListenAndServe(":"+port, r.GetRouter())

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

func (r *App) InitDatabase(name, path string) {
	dbURI := r.Configuration.Get("DB_URI")
	dbSlowThreshold := r.Configuration.GetInt64("DB_SLOW_THRESHOLD")
	logQuery := r.Configuration.Get("LOG_QUERY")

	dsn := dbURI + "?charset=utf8mb4&parseTime=True&loc=Local"

	dbLogger := gorm_logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gorm_logger.Config{
		SlowThreshold:             time.Duration(dbSlowThreshold) * time.Millisecond,
		LogLevel:                  gorm_logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})

	logg := dbLogger.LogMode(gorm_logger.Warn)

	if logQuery != "" {
		logg = dbLogger.LogMode(gorm_logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logg,
	})

	if err != nil {
		log.Panicln("Error on connect in database", err)
	}

	r.DB = db
}

func newApp() *App {
	var app App

	app.Events = event.NewManager("app")
	app.RolesString = configuration.Roles

	logger.Init()
	app.Configuration = configuration.NewCfg()
	app.routerGroups = make(map[string]*echo.Group)
	app.apiRouterGroups = make(map[string]*echo.Group)

	app.router = echo.New()
	app.Plugins = make(map[string]Pluginer)

	app.SetRouterGroup("main", "/")
	app.SetRouterGroup("public", "/public")

	apiRouterGroup := app.SetRouterGroup("api", "/api")
	apiRouterGroup.GET("", HealthCheck)

	app.router.Validator = &helpers.CustomValidator{Validator: validator.New()}

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
