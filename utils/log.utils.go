package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/deestarks/infiniti/utils/constants"
)

func LogMessage(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	// Using colorable log info
	logColors := constants.NewLoggerColors()
	logInfo := fmt.Sprintf("%s[%v Infiniti-API]:%s ", logColors.Blue, time.Now().Format("2006-01-02 15:04:05"), logColors.Reset)

	// Logging
	log.SetFlags(0)
	log.Println(logInfo+msg)
}