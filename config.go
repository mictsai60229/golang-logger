package logger

import (
	"io"
)

type LoggerConfig struct {
	serviceName string
	appEnv      string
	logFile     io.Writer
}

var loggerConfig = new(LoggerConfig)

func SetServiceName(serviceName string) {
	loggerConfig.serviceName = serviceName
}

func SetAppEnv(appEnv string) {
	loggerConfig.appEnv = appEnv
}

func SetLogFile(filepath string) {
	file := getRotationLog(filepath)
	loggerConfig.logFile = file
}
