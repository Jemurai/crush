package utils

import (
	"crypto/sha256"
	"encoding/hex"
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
	var hash string
	if input != "" {
		hasher := sha256.New()
		hasher.Write([]byte(input))
		hash = hex.EncodeToString(hasher.Sum(nil))
	} else {
		hash = input
	}
	return hash
}

// Truncate ensures the sizes of strings are under control
func Truncate(input string, size int) string {
	result := input
	if len(input) > size {
		if size > 3 {
			size -= 3
		}
		result = input[0:size] + "..."
	}
	return result
}
