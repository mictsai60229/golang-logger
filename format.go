package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// request
func NewRequest(context *gin.Context, requestState *RequestState) *Request {
	request := new(Request)

	request.message = "log request: " + context.Request.URL.Path
	request.system = &System{
		level: "info",
	}
	request.logLabel = "REQUEST"
	request.requestData = &RequestRequestData{
		RequestData{
			request: context.Request,
			uuid:    requestState.GetRequestUUID(),
		},
	}

	request.traceInfo = &TraceInfo{
		runTime: requestState.DiffTime(),
	}

	return request
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
	return CurrentTimeWithMircoseconds()
}

func (request *Request) getTracing() map[string]interface{} {
	return request.tracing.Data()
}

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

func NewResponse(context *gin.Context, requestState *RequestState, writer *bodyLogWriter) *Response {
	response := new(Response)

	response.message = "log response: " + context.Request.URL.Path
	response.system = &System{
		level: "info",
	}
	response.logLabel = "RESPONSE"
	response.requestData = &ResponseRequestData{
		RequestData{
			request: context.Request,
			uuid:    requestState.GetRequestUUID(),
		},
	}

	response.responseData = &ResponseData{
		Writer: writer,
		time:   requestState.DiffTime(),
	}

	response.traceInfo = &TraceInfo{
		runTime: requestState.DiffTime(),
	}

	return response
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
	return CurrentTimeWithMircoseconds()
}

func (response *Response) getTracing() map[string]interface{} {
	return response.tracing.Data()
}

//system
type System struct {
	level string
}

func (system *System) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["service_name"] = system.getServiceName()
	data["pid"] = system.getPid()
	data["level"] = system.getLevel()
	data["env"] = system.getEnv()

	return data
}

func (system *System) getServiceName() string {
	return loggerConfig.serviceName
}

func (system *System) getPid() int {
	return os.Getpid()
}

func (system *System) getLevel() string {
	return system.level
}

func (system *System) getEnv() string {
	return loggerConfig.appEnv
}

// request data
type RequestData struct {
	request *http.Request
	uuid    string
}

type RequestRequestData struct {
	RequestData
}

type ResponseRequestData struct {
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

func (requestData *ResponseRequestData) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["method"] = requestData.getMethod()
	data["route"] = requestData.getRoute()
	data["ip"] = requestData.getIP()
	data["uuid"] = requestData.getUUID()

	return data
}

func (requestData *RequestData) getMethod() string {
	return requestData.request.Method
}

func (requestData *RequestData) getUrl() string {
	return requestData.getHost() + requestData.getRoute()
}

func (requestData *RequestData) getHost() string {
	return requestData.request.URL.Host
}

func (requestData *RequestData) getRoute() string {
	return requestData.request.URL.Path
}

func (requestData *RequestData) getBody() string {
	var buffer bytes.Buffer
	tee := io.TeeReader(requestData.request.Body, &buffer)
	body, _ := ioutil.ReadAll(tee)
	requestData.request.Body = ioutil.NopCloser(&buffer)
	return string(body)
}

func (requestData *RequestData) getHeaders() string {
	headers, err := json.Marshal(requestData.request.Header)
	if err != nil {
		panic("JSON decoding errors")
	}
	return string(headers)
}

func (requestData *RequestData) getIP() string {
	request := requestData.request

	IPAddress := request.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = request.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = request.RemoteAddr
	}
	return IPAddress
}
func (requestData *RequestData) getUUID() string {
	return requestData.uuid
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

type TraceInfo struct {
	runTime int64
}

func (traceInfo *TraceInfo) Data() (data map[string]interface{}) {
	data = make(map[string]interface{})

	data["run_time"] = traceInfo.getRunTime()

	return data
}

func (traceInfo *TraceInfo) getRunTime() int64 {
	return traceInfo.runTime
}

type ResponseData struct {
	Writer   *bodyLogWriter
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
	return responseData.Writer.body.String()
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
	err := json.Unmarshal(responseData.Writer.body.Bytes(), &body)
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

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (writer bodyLogWriter) Write(b []byte) (int, error) {
	writer.body.Write(b)
	return writer.ResponseWriter.Write(b)
}
