package utils

import (
    "encoding/json"
)

// The struct - "obj", fields must have json tags
func StructToMap(obj interface{}) map[string]interface{} {
	b, _ := json.Marshal(obj)
	var result map[string]interface{}
	json.Unmarshal(b, &result)
	return result
}

func StructSliceToMap(obj interface{}) []map[string]interface{} {
	b, _ := json.Marshal(obj)
	var result []map[string]interface{}
	json.Unmarshal(b, &result)
	return result
}