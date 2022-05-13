package logger

import (
	"encoding/json"
	"time"
)

func CurrentTimeWithMircoseconds() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05.000000")
}

func JsonEncode(data map[string]interface{}) []byte {
	jsondata, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsondata
}
