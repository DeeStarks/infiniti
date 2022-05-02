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
	logInfo := fmt.Sprintf("%s[%v Infiniti]:%s ", logColors.Blue, time.Now().Format("2006-01-02 15:04:05"), logColors.Reset)

	// Logging
	log.SetFlags(0)
	log.Println(logInfo+msg)
}

// Logger to be used in the Jobs layer
func LogJobMessage(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	// Using colorable log info
	logColors := constants.NewLoggerColors()
	logInfo := fmt.Sprintf("%s[%v Infiniti Jobs]:%s ", logColors.Magenta, time.Now().Format("2006-01-02 15:04:05"), logColors.Reset)

	// Logging
	log.SetFlags(0)
	log.Println(logInfo+msg)
}

// 
type message struct {}

func Printer() *message {
	return &message{}
}

func (m *message) Message(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(msg)
}

func (m *message) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	colors := constants.NewLoggerColors()
	msg = fmt.Sprintf("%s%s%s", colors.Red, msg, colors.Reset)
	fmt.Println(msg)
}

func (m *message) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	colors := constants.NewLoggerColors()
	msg = fmt.Sprintf("%s%s%s", colors.Green, msg, colors.Reset)
	fmt.Println(msg)
}

func (m *message) Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	colors := constants.NewLoggerColors()
	msg = fmt.Sprintf("%s%s%s", colors.Yellow, msg, colors.Reset)
	fmt.Println(msg)
}
