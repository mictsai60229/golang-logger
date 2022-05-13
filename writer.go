package logger

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (writer bodyLogWriter) Write(b []byte) (int, error) {
	writer.body.Write(b)
	return writer.ResponseWriter.Write(b)
}

func (writer bodyLogWriter) GetResponseBody() string {
	return writer.body.String()
}

func getRotationLog(logPath string) io.Writer {
	logFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500,
		MaxBackups: 10,
	}

	return logFile
}
