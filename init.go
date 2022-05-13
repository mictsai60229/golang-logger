package logger

import (
	"os"
)

func init() {
	loggerConfig.logFile = os.Stdout
}
