package utils

import (
	"fmt"
)

// Create placeholder string with the input "length" long
// (e.g. 1 returns $1, 2 returns $1, $2)
func CreatePlaceholder(length int) string {
	var placeholder string
	for i := 0; i < length; i++ {
		placeholder += "$" + fmt.Sprintf("%d", i+1) + ", "
	}
	return placeholder[:len(placeholder)-2]
}

type SplitMap struct {
	Key	 	string
	Value	interface{}
}
// Split map[string]interface{} into []SplitMap
// (e.g map[string]interface{"id": 1, "name": "John"} returns []SplitMap{{"id", 1}, {"name", "John"}})
func MapToStructSlice(data map[string]interface{}) []SplitMap {
	var splitMap []SplitMap
	for key, val := range data {
		splitMap = append(splitMap, SplitMap{key, val})
	}
	return splitMap
}

// Create "SET" conditions string with the input = "data: []SplitMap"
// Values are replaced with placeholders
// (e.g. []SplitMap{{"id", 1}, {"name", "John"}} returns "id = $1, name = $2")
func CreateSetConditions(data []SplitMap) string {
	var setConditions string
	for i, val := range data {
		setConditions += val.Key + " = $" + fmt.Sprintf("%d", i+1) + ", "
	}
	return setConditions[:len(setConditions)-2]
}