package main

import (
	"log"
)

func logMessage(msg string, level string) {
	log.Printf("[%s] %s", level, msg)
}

func logWarn(msg string) {
	logMessage(msg, "WARN")
}

func logWarnErr(err error, msg string) {
	if err != nil {
		if msg != "" {
			msg += ", "
		}
		msg += err.Error()
	}
	logWarn(msg)
}

func logInfo(msg string) {
	logMessage(msg, "INFO")
}
