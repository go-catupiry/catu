package catu

import (
	"github.com/go-catupiry/catu/configuration"
	"github.com/go-catupiry/catu/logger"
)

type App struct {
	Cfg configuration.Configer
}

func (a *App) Init() {
	a.Cfg.Init()
	logger.Init(a.Cfg)
}
