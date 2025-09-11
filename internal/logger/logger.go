package logger

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(mode string) *logrus.Logger {
	log := logrus.New()
	if mode == "release" {
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
		log.SetLevel(logrus.DebugLevel)
	}
	return log
}
