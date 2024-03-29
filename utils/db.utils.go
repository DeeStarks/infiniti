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

func SliceToMap(cols []string, vals []interface{}) map[string]interface{} {
	var data = make(map[string]interface{}, len(cols))
	for i, col := range cols {
		data[col] = vals[i]
	}
	return data
}

// Create "SET" conditions string with the input "data"
// Values are replaced with placeholders
// (e.g. []SplitMap{{"id", 1}, {"name", "John"}} returns "id = $1, name = $2")
// NB: To be used together with the MapToStructSlice function
func CreateSetConditions(data []SplitMap) string {
	var setConditions string
	for i, val := range data {
		setConditions += val.Key + " = $" + fmt.Sprintf("%d", i+1) + ", "
	}
	return setConditions[:len(setConditions)-2]
}