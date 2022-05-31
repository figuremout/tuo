package json

import (
	"encoding/json"
)

// Dumps,Loads is valid for go objects like struct, string, slice
// encode golang object into json string
func Dumps(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("")
	}
	return b
}

// parse json bytes into golang object
func Loads(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func StructToMap(v interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	b, _ := json.Marshal(v)
	json.Unmarshal(b, &m)
	return m
}
