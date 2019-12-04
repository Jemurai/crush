package utils

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func Timing(start time.Time, message string) time.Time {
	current := time.Now()
	elapsed := current.Sub(start)
	log.Debugf(message, elapsed.Seconds())
	return current
}
