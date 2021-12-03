package catu

import "github.com/go-catupiry/catu/logger"

func NewApp() App {
	c := App{}
	return c
}

// Init global features like logrus configuration
func Init() error {
	logger.Init()

	return nil
}
