package logger

import (
	"os"

	"github.com/go-catupiry/catu/configuration"
	"github.com/sirupsen/logrus"
)

func Init() {
	GO_ENV := configuration.GetEnv("GO_ENV", "development")
	if GO_ENV != "development" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			DataKey: "data",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime: "timestamp",
				logrus.FieldKeyMsg:  "message",
			},
		})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	LOG_LV := configuration.GetEnv("LOG_LV", "")

	switch LOG_LV {
	case "verbose":
		// Only log the warning severity or above.
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
