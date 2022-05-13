package format

import (
	"github.com/gin-gonic/gin"
	helper "github.com/mictsai60229/golang-logger/helper"
)

type Request struct {
	message     string
	system      SystemInterface
	requestData RequestDataInterface
	traceInfo   TraceInfoInterface
	logLabel    string
	tracing     TracingInterface
}

func (request *Request) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["message"] = request.getMessage()
	data["system"] = request.getSystem()
	data["request"] = request.getRequestData()
	data["trace_info"] = request.getTraceInfo()
	data["log_label"] = request.getLogLabel()
	data["datetime"] = request.getDateTime()
	//data["tracing"] = request.getTracing()

	return data
}

func (request *Request) getMessage() string {
	return request.message
}

func (request *Request) getSystem() map[string]interface{} {
	return request.system.Data()
}

func (request *Request) getRequestData() map[string]interface{} {
	return request.requestData.Data()
}

func (request *Request) getTraceInfo() map[string]interface{} {
	return request.traceInfo.Data()
}

func (request *Request) getLogLabel() string {
	return request.logLabel
}

func (request *Request) getDateTime() string {
	return helper.CurrentTimeWithMircoseconds()
}

func (request *Request) getTracing() map[string]interface{} {
	return request.tracing.Data()
}

type RequestRequestData struct {
	RequestData
}

func (requestData *RequestRequestData) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["method"] = requestData.getMethod()
	data["url"] = requestData.getUrl()
	data["host"] = requestData.getHost()
	data["route"] = requestData.getRoute()
	data["body"] = requestData.getBody()
	data["headers"] = requestData.getHeaders()
	data["ip"] = requestData.getIP()
	data["uuid"] = requestData.getUUID()

	return data
}

func NewRequest(context *gin.Context, serviceName string, appEnv string, requestUUID string, runTime int64) *Request {
	request := new(Request)

	request.message = "log request: " + context.Request.URL.Path
	request.system = &System{
		level:       "info",
		serviceName: serviceName,
		appEnv:      appEnv,
	}
	request.logLabel = "REQUEST"
	request.requestData = &RequestRequestData{
		RequestData{
			request: context.Request,
			uuid:    requestUUID,
		},
	}

	request.traceInfo = &TraceInfo{
		runTime: runTime,
	}

	return request
}
