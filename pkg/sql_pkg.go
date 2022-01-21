package pkg

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