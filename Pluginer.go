package catu

type Pluginer interface {
	Init(app *App) error

	BeforeBindMiddlewares(app *App) error
	BindMiddlewares(app *App) error
	AfterBindMiddlewares(app *App) error

	BeforeBindRoutes(app *App) error
	BindRoutes(app *App) error
	AfterBindRoutes(app *App) error

	SetTemplateFuncMap(app *App) error

	OnBootstrap(app *App) error
}
