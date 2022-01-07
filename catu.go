package catu

import "github.com/go-catupiry/catu/logger"

// Init global features like logrus configuration
func Init() error {
	logger.Init()

	return nil
}
