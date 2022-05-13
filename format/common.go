package format

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//system
type System struct {
	level       string
	serviceName string
	appEnv      string
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
	return system.serviceName
}

func (system *System) getPid() int {
	return os.Getpid()
}

func (system *System) getLevel() string {
	return system.level
}

func (system *System) getEnv() string {
	return system.appEnv
}

// request data
type RequestData struct {
	request *http.Request
	uuid    string
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
