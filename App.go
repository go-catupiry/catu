package catu

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/go-catupiry/catu/acl"
	"github.com/go-catupiry/catu/configuration"
	"github.com/go-catupiry/catu/helpers"
	"github.com/go-catupiry/catu/http_client"
	"github.com/go-catupiry/catu/logger"
	"github.com/go-catupiry/catu/pagination"
	"github.com/go-catupiry/query_parser_to_db"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/event"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	gorm_logger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

type App interface {
	RegisterPlugin(p Pluginer)
	GetPlugins() map[string]Pluginer
	GetPlugin(name string) Pluginer
	SetPlugin(name string, plugin Pluginer) error

	GetRouter() *echo.Echo
	SetRouterGroup(name, path string) *echo.Group
	GetRouterGroup(name string) *echo.Group
	SetResource(name string, httpController HTTPController, routerGroup *echo.Group) error
	StartHTTPServer() error
	NewRequestContext(opts *RequestContextOpts) *RequestContext

	GetTemplates() *template.Template
	LoadTemplates() error
	SetTemplateFunction(name string, f interface{})

	InitDatabase(name, engine string, isDefault bool) error
	SetModel(name string, f interface{})
	GetModel(name string) interface{}

	Can(permission string, userRoles []string) bool

	GetEvents() *event.Manager

	GetConfiguration() configuration.ConfigurationInterface

	GetDB() *gorm.DB
	SetDB(db *gorm.DB) error
	Migrate() error

	Bootstrap() error
	Close() error
}

type AppStruct struct {
	InitTime time.Time

	Events *event.Manager

	Configuration configuration.ConfigurationInterface
	// Default database
	DB *gorm.DB
	// avaible databases
	DBs map[string]*gorm.DB

	Plugins map[string]Pluginer

	Models map[string]interface{}

	router    *echo.Echo
	Resources map[string]*HTTPResource

	routerGroups    map[string]*echo.Group
	apiRouterGroups map[string]*echo.Group

	RolesString string
	RolesList   map[string]acl.Role

	templates         *template.Template
	templateFunctions template.FuncMap
}

func (r *AppStruct) RegisterPlugin(p Pluginer) {
	if p.GetName() == "" {
		panic("Plugin.RegisterPlugin Name should be returned from GetName method")
	}

	r.Plugins[p.GetName()] = p
}

func (r *AppStruct) GetPlugin(name string) Pluginer {
	return r.Plugins[name]
}

func (r *AppStruct) SetPlugin(name string, plugin Pluginer) error {
	r.Plugins[name] = plugin
	return nil
}

func (r *AppStruct) GetPlugins() map[string]Pluginer {
	return r.Plugins
}

func (r *AppStruct) GetRouter() *echo.Echo {
	return r.router
}

func (app *AppStruct) NewRequestContext(opts *RequestContextOpts) *RequestContext {
	cfg := app.GetConfiguration()
	port := cfg.GetF("PORT", "8080")
	protocol := cfg.GetF("PROTOCOL", "http")
	domain := cfg.GetF("DOMAIN", "localhost")

	ctx := RequestContext{
		EchoContext: opts.EchoContext,
		Protocol:    protocol,
		Domain:      domain,
		AppOrigin:   cfg.GetF("APP_ORIGIN", protocol+"://"+domain+":"+port),
		// Title:               "",
		Layout: "site/layouts/default",
		ENV:    cfg.GetF("GO_ENV", "development"),
		Query:  query_parser_to_db.NewQuery(50),
		Pager:  pagination.NewPager(),
	}

	if opts.EchoContext == nil {
		opts.EchoContext = echo.New().NewContext(&http.Request{}, &helpers.FakeResponseWriter{})
	}

	// Is a context used on CLIs, not in HTTP request / echo then skip it
	if opts.EchoContext == nil || ctx.Request().URL == nil {
		return &ctx
	}

	ctx.Pager.CurrentUrl = ctx.Request().URL.Path
	ctx.Pager.Limit, _ = strconv.ParseInt(cfg.GetF("PAGER_LIMIT", "20"), 10, 64)

	ctx.MetaTags.Title = cfg.Get("SITE_NAME")
	ctx.MetaTags.Description = cfg.Get("SITE_DESCRIPTION")
	ctx.MetaTags.ImageURL = cfg.Get("SITE_IMAGE_URL")
	ctx.MetaTags.SiteName = cfg.Get("SITE_NAME")

	if opts.EchoContext.Request().Method != "GET" {
		return &ctx
	}

	limitMax, _ := strconv.ParseInt(app.GetConfiguration().GetF("PAGER_LIMIT_MAX", "50"), 10, 64)

	rawParams := opts.EchoContext.QueryParams()

	filteredParamArray := []string{}

	for key, param := range rawParams {
		// get limit with max value for security:
		if key == "limit" && len(param) == 1 {
			queryLimit, err := strconv.ParseInt(param[0], 10, 64)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"key":   key,
					"param": param,
				}).Error("NewRequestContext invalid query param limit")
				continue
			}
			if queryLimit > 0 && queryLimit < limitMax {
				ctx.Pager.Limit = queryLimit
			}
		}

		if key == "page" && len(param) == 1 {
			page, _ := strconv.ParseInt(param[0], 10, 64)
			ctx.Pager.Page = page
			continue
		}

		ctx.Query.AddQueryParamFromRaw(key, param)
	}

	if len(filteredParamArray) > 0 {
		strings.Join(filteredParamArray[:], ",")
	}

	return &ctx
}

func (r *AppStruct) GetTemplates() *template.Template {
	return r.templates
}

func (r *AppStruct) GetEvents() *event.Manager {
	return r.Events
}

func (r *AppStruct) GetConfiguration() configuration.ConfigurationInterface {
	return r.Configuration
}

func (r *AppStruct) GetDB() *gorm.DB {
	return r.DB
}
func (r *AppStruct) SetDB(db *gorm.DB) error {
	r.DB = db
	return nil
}

func (r *AppStruct) Bootstrap() error {
	var err error

	logrus.Debug("catu.App.Bootstrap running")
	// default roles and permissions, override it on your app
	json.Unmarshal([]byte(r.RolesString), &r.RolesList)

	for _, p := range r.Plugins {
		err = p.Init(r)
		if err != nil {
			return errors.Wrap(err, "App.Bootstrap | Error on run plugin init "+p.GetName())
		}
	}

	r.Events.MustTrigger("configuration", event.M{"app": r})

	err = r.InitDatabase("default", configuration.GetEnv("DB_ENGINE", "sqlite"), true)
	if err != nil {
		return err
	}

	http_client.Init()

	r.Events.MustTrigger("bindMiddlewares", event.M{"app": r})
	r.Events.MustTrigger("bindRoutes", event.M{"app": r})
	r.Events.MustTrigger("setResponseFormats", event.M{"app": r})
	r.Events.MustTrigger("setTemplateFunctions", event.M{"app": r})

	logrus.WithFields(logrus.Fields{
		"count": len(r.templateFunctions),
	}).Debug("catu.App.Bootstrap template functions loaded")

	err = r.LoadTemplates()
	if err != nil {
		return errors.Wrap(err, "App.Bootstrap Error on LoadTemplates")
	}

	r.router.Renderer = &TemplateRenderer{
		templates: r.GetTemplates(),
	}

	r.Events.MustTrigger("bootstrap", event.M{"app": r})

	return nil
}

func (r *AppStruct) StartHTTPServer() error {
	port := r.Configuration.Get("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Info("Server listening on port " + port)
	return http.ListenAndServe(":"+port, r.GetRouter())
}

func (r *AppStruct) SetRouterGroup(name, path string) *echo.Group {
	if r.routerGroups[name] == nil {
		r.routerGroups[name] = r.router.Group(path)
	}
	return r.routerGroups[name]
}

func (r *AppStruct) GetRouterGroup(name string) *echo.Group {
	return r.routerGroups[name]
}

// Set Resource CRUD.
// Now we only supports HTTP Resources / Ex Rest
func (r *AppStruct) SetResource(name string, httpController HTTPController, routerGroup *echo.Group) error {
	routerGroup.GET("", httpController.Query)
	routerGroup.GET("/count", httpController.Count)
	routerGroup.POST("", httpController.Create)
	routerGroup.GET("/:id", httpController.FindOne)
	routerGroup.POST("/:id", httpController.Update)
	routerGroup.PATCH("/:id", httpController.Update)
	routerGroup.PUT("/:id", httpController.Update)
	routerGroup.DELETE("/:id", httpController.Delete)

	r.Resources[name] = &HTTPResource{
		Name:       name,
		Controller: &httpController,
	}

	return nil
}

func (r *AppStruct) InitDatabase(name, engine string, isDefault bool) error {
	var err error
	var db *gorm.DB

	dbURI := r.Configuration.GetF("DB_URI", "test.sqlite?charset=utf8mb4")
	dbSlowThreshold := r.Configuration.GetInt64F("DB_SLOW_THRESHOLD", 400)
	logQuery := r.Configuration.GetF("LOG_QUERY", "")

	logrus.WithFields(logrus.Fields{
		"dbURI":           dbURI,
		"dbSlowThreshold": dbSlowThreshold,
		"logQuery":        logQuery,
	}).Debug("catu.App.InitDatabase starting db with configs")

	if dbURI == "" {
		return errors.New("catu.App.InitDatabase DB_URI environment variable is required")
	}

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

	switch engine {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logg,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbURI), &gorm.Config{
			Logger: logg,
		})

	default:
		return errors.New("catu.App.InitDatabase invalid database engine. Options available: mysql or sqlite")
	}

	if err != nil {
		return errors.Wrap(err, "catu.App.InitDatabase error on database connection")
	}

	if isDefault {
		r.DB = db
	}

	return nil
}

func (r *AppStruct) SetModel(name string, f interface{}) {
	r.Models[name] = f
}

func (r *AppStruct) GetModel(name string) interface{} {
	return r.Models[name]
}

func (r *AppStruct) SetTemplateFunction(name string, f interface{}) {
	r.templateFunctions[name] = f
}

func (r *AppStruct) Can(permission string, userRoles []string) bool {
	// first check if user is administrator
	for i := range userRoles {
		if userRoles[i] == "administrator" {
			return true
		}
	}

	for j := range userRoles {
		R := r.RolesList[userRoles[j]]
		if R.Can(permission) {
			return true
		}
	}

	return false
}

func (r *AppStruct) LoadTemplates() error {
	rootDir := r.Configuration.GetF("TEMPLATE_FOLDER", "./templates")
	disableTemplating := r.Configuration.GetBool("TEMPLATE_DISABLE")

	if disableTemplating {
		return nil
	}

	tpls, err := findAndParseTemplates(rootDir, r.templateFunctions)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			// "error":   errHealthCheckHandlerr,
			"rootDir": rootDir,
		}).Error("catu.App.LoadTemplates Error on parse templates")
		r.templates = tpls
		return err
	}

	r.templates = tpls

	logrus.WithFields(logrus.Fields{
		"count": len(r.templates.Templates()),
	}).Debug("catu.App.ParseTemplates templates loaded")

	return nil
}

// Run migrations
func (r *AppStruct) Migrate() error {
	err, _ := r.Events.Fire("migrate", event.M{"app": r})
	if err != nil {
		return errors.Wrap(err, "App.Migrate migrate error")
	}

	return nil
}

// Method for close and end all app operations, use that before close the app execution
func (r *AppStruct) Close() error {
	err, _ := r.Events.Fire("close", event.M{"app": r})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Debug("catu.App.Close error")
	}

	return nil
}

func newApp() App {
	var app AppStruct

	app.Events = event.NewManager("app")
	app.RolesString, _ = acl.LoadRoles()

	logger.Init()
	app.Configuration = configuration.NewCfg()
	app.routerGroups = make(map[string]*echo.Group)

	app.Resources = make(map[string]*HTTPResource)
	app.router = echo.New()

	app.router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &RequestContext{
				EchoContext: c,
			}
			return next(cc)
		}
	})

	app.router.Binder = &CustomBinder{}
	app.router.HTTPErrorHandler = CustomHTTPErrorHandler
	app.router.Validator = &helpers.CustomValidator{Validator: validator.New()}

	app.router.GET("/health", HealthCheckHandler)
	app.Plugins = make(map[string]Pluginer)

	app.templates = &template.Template{}

	app.SetRouterGroup("main", "/")
	app.SetRouterGroup("public", "/public")

	apiRouterGroup := app.SetRouterGroup("api", "/api")
	apiRouterGroup.GET("", HealthCheckHandler)

	app.templateFunctions = sprig.FuncMap()

	return &app
}
