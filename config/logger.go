package config

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Level string
}

func (loggerConfig *Logger) GetLogLevel() (logrus.Level, error) {
	return logrus.ParseLevel(loggerConfig.Level)
}
