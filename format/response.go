package format

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	helper "github.com/mictsai60229/golang-logger/helper"
)

// response
type Response struct {
	message      string
	system       SystemInterface
	requestData  RequestDataInterface
	responseData ResponseDataInterface
	traceInfo    TraceInfoInterface
	logLabel     string
	tracing      TracingInterface
}

func (response *Response) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["message"] = response.getMessage()
	data["system"] = response.getSystem()
	data["request"] = response.getRequestData()
	data["response"] = response.getResponseData()
	data["trace_info"] = response.getTraceInfo()
	data["log_label"] = response.getLogLabel()
	data["datetime"] = response.getDateTime()
	//data["tracing"] = request.getTracing()

	return data
}

func (response *Response) getMessage() string {
	return response.message
}

func (response *Response) getSystem() map[string]interface{} {
	return response.system.Data()
}

func (response *Response) getRequestData() map[string]interface{} {
	return response.requestData.Data()
}

func (response *Response) getResponseData() map[string]interface{} {
	return response.responseData.Data()
}

func (response *Response) getTraceInfo() map[string]interface{} {
	return response.traceInfo.Data()
}

func (response *Response) getLogLabel() string {
	return response.logLabel
}

func (response *Response) getDateTime() string {
	return helper.CurrentTimeWithMircoseconds()
}

func (response *Response) getTracing() map[string]interface{} {
	return response.tracing.Data()
}

type ResponseRequestData struct {
	RequestData
}

func (requestData *ResponseRequestData) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["method"] = requestData.getMethod()
	data["route"] = requestData.getRoute()
	data["ip"] = requestData.getIP()
	data["uuid"] = requestData.getUUID()

	return data
}

type ResponseData struct {
	Writer   ResponseBodyLogger
	time     int64
	metaData map[string]interface{}
}

func (responseData *ResponseData) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["http_status"] = responseData.getHttpStatus()
	data["body"] = responseData.getBody()
	data["headers"] = responseData.getHeaders()
	data["time"] = responseData.getTime()
	data["meta_status"] = responseData.getMetaStatus()
	data["meta_description"] = responseData.getMetaDescription()

	return data
}

func (responseData *ResponseData) getHttpStatus() int {
	return responseData.Writer.Status()
}
func (responseData *ResponseData) getBody() string {
	return responseData.Writer.GetResponseBody()
}
func (responseData *ResponseData) getHeaders() string {
	headers, err := json.Marshal(responseData.Writer.Header())
	if err != nil {
		panic("JSON decoding errors")
	}
	return string(headers)
}
func (responseData *ResponseData) getTime() int64 {
	return responseData.time
}

func (responseData *ResponseData) getMetaData() map[string]interface{} {
	if responseData.metaData != nil {
		return responseData.metaData
	}

	body := make(map[string]interface{})
	body_bytes := []byte(responseData.Writer.GetResponseBody())

	err := json.Unmarshal(body_bytes, &body)
	if err != nil {
		panic("JSON encoding errors")
	}

	if value, ok := body["metadata"]; ok {
		responseData.metaData = value.(map[string]interface{})
	}

	return responseData.metaData
}

func (responseData *ResponseData) getMetaStatus() string {

	metaData := responseData.getMetaData()
	if value, ok := metaData["status"]; ok {
		return value.(string)
	}
	return ""
}

func (responseData *ResponseData) getMetaDescription() string {
	metaData := responseData.getMetaData()
	if value, ok := metaData["desc"]; ok {
		return value.(string)
	}
	return ""
}

func NewResponse(context *gin.Context, serviceName string, appEnv string, requestUUID string, runTime int64, writer ResponseBodyLogger) *Response {
	response := new(Response)

	response.message = "log response: " + context.Request.URL.Path
	response.system = &System{
		level:       "info",
		serviceName: serviceName,
		appEnv:      appEnv,
	}
	response.logLabel = "RESPONSE"
	response.requestData = &ResponseRequestData{
		RequestData{
			request: context.Request,
			uuid:    requestUUID,
		},
	}

	response.responseData = &ResponseData{
		Writer: writer,
		time:   runTime,
	}

	response.traceInfo = &TraceInfo{
		runTime: runTime,
	}

	return response
}
