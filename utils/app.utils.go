package utils

import (
	"encoding/json"
	"net/mail"
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

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}