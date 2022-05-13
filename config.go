package logger

import (
	"fmt"
	"io"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
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
	file, err := getRotationLog(filepath)
	if err != nil {
		fmt.Printf("failed to create rotatelogs: %s", err)
	} else {
		loggerConfig.logFile = file
	}
}

func getRotationLog(logPath string) (*rotateLogs.RotateLogs, error) {
	logFile, err := rotateLogs.New(
		logPath+".%Y%m%d%H%M",
		rotateLogs.WithLinkName(logPath),
		rotateLogs.WithRotationCount(10),
		rotateLogs.WithRotationSize(100*1024*1024),
	)

	return logFile, err
}
