package utils

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log *logrus.Logger

func init() {
	log = logrus.New()

	formatter := runtime.Formatter{ChildFormatter: &logrus.TextFormatter{
		FullTimestamp: true,
	}}
	formatter.File = true
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func Logger() *logrus.Logger {
	return log
}
