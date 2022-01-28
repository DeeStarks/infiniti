package lib

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

// Split map[string]interface{} into keys (string) and values ([]interface{})
// (e.g map[string]interface{"id": 1, "name": "John"} returns ("id, name" and []interface{}{1, "John"})
func SplitMap(data map[string]interface{}) (string, []interface{}) {
	var (
		colStr		string
		valArr		[]interface{}
	)

	for col, val := range data {
		colStr += col + ", "
		valArr = append(valArr, val)
	}
	colStr = colStr[:len(colStr)-2] // remove the last ", "
	return colStr, valArr
}