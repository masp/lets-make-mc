package log

import (
	log "github.com/sirupsen/logrus"
)

func SetupLogging() {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04:05"
	formatter.FullTimestamp = false

	log.SetFormatter(formatter)
}
