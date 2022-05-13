package logger

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {

		writer := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: context.Writer,
		}
		context.Writer = writer

		requestState := NewRequestState()
		request := NewRequest(context, requestState)
		// log request
		LogRequest(request)

		context.Next()

		// log response
		response := NewResponse(context, requestState, writer)
		LogRequest(response)
	}
}

func LogRequest(request RequestInterface) {
	data := request.Data()
	jsondata := JsonEncode(data)
	loggerConfig.logFile.Write(jsondata)
	newline := []byte("\n")
	loggerConfig.logFile.Write(newline)
}

func LogResponse(response ResponseInterface) {
	data := response.Data()
	jsondata := JsonEncode(data)
	loggerConfig.logFile.Write(jsondata)
	newline := []byte("\n")
	loggerConfig.logFile.Write(newline)
}
