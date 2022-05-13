package helper

import "encoding/json"

func JsonEncode(data map[string]interface{}) []byte {
	jsondata, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsondata
}
