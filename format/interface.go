package format

import "github.com/gin-gonic/gin"

type LoggerInterface interface {
	Data() (data map[string]interface{})
}

type RequestInterface interface {
	LoggerInterface
}

type ResponseInterface interface {
	LoggerInterface
}

type RequestStateInterface interface {
	GetRequestUUID() string
	DiffTime() int64
}

type SystemInterface interface {
	LoggerInterface
	getServiceName() string
	getPid() int
	getLevel() string
	getEnv() string
}

type RequestDataInterface interface {
	LoggerInterface
	getMethod() string
	getUrl() string
	getHost() string
	getRoute() string
	getBody() string
	getHeaders() string
	getIP() string
	getUUID() string
}
type ResponseDataInterface interface {
	LoggerInterface
	getHttpStatus() int
	getBody() string
	getHeaders() string
	getTime() int64
	getMetaStatus() string
	getMetaDescription() string
}

type TraceInfoInterface interface {
	LoggerInterface
	getRunTime() int64
}

type TracingInterface interface {
	LoggerInterface
}

type ResponseBodyLogger interface {
	gin.ResponseWriter
	GetResponseBody() string
}
