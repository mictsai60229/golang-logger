package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mictsai60229/golang-logger/format"
)

// request
func NewRequest(context *gin.Context, requestState *RequestState) *format.Request {
	return format.NewRequest(context, loggerConfig.serviceName, loggerConfig.appEnv, requestState.GetRequestUUID(), requestState.DiffTime())
}

func NewResponse(context *gin.Context, requestState *RequestState, writer *bodyLogWriter) *format.Response {
	return format.NewResponse(context, loggerConfig.serviceName, loggerConfig.appEnv, requestState.GetRequestUUID(), requestState.DiffTime(), writer)
}

// request state
type RequestState struct {
	time time.Time
	uuid string
}

func NewRequestState() *RequestState {
	requestState := new(RequestState)

	uuidBtye, _ := uuid.NewRandom()
	requestState.uuid = uuidBtye.String()
	requestState.time = time.Now()

	return requestState
}

func (requestState *RequestState) GetRequestUUID() string {
	return requestState.uuid
}

func (requestState *RequestState) DiffTime() int64 {
	now := time.Now()
	duration := now.Sub(requestState.time)
	return duration.Milliseconds()
}
