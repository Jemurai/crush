package utils

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Timing provides a time slice and message
func Timing(start time.Time, message string) time.Time {
	current := time.Now()
	elapsed := current.Sub(start)
	log.Debugf(message, elapsed.Seconds())
	return current
}

// Fingerprint provides a unique string for a given input.
func Fingerprint(input string) string {

	return input // This works for now ...

	// TODO:  FIGURE out how to make a byte[] from a string
	// hasher := sha256.New()
	//	hasher.Write(input)
	//	hash := hex.EncodeToString(hasher.Sum(nil))
	// 	return hash
}
