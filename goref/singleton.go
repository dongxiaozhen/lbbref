package goref

import (
	"encoding/json"
	"time"

	"github.com/donnie4w/go-logger/logger"
)

// singleton GoRef instance
var instance = NewGoRef()

// GetInstance -- Returns a scoped instance (matching the given scope path)
func GetInstance(path ...string) *GoRef {
	return instance.GetChild(path...)
}

// GetSnapshot -- Returns a Snapshot of the GoRef  (synchronously)
func GetSnapshot() Snapshot {
	return instance.GetSnapshot()
}

// Ref -- References an instance of 'key' (in singleton mode)
func Ref(key string) *Instance {
	return instance.Ref(key)
}

// Reset -- resets the internal state of the singleton GoRef instance
func Reset() {
	instance.Reset()
}

func init() {
	go hitPoint()
}

func hitPoint() {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		s := GetSnapshot().Data
		if len(s) <= 0 {
			continue
		}
		data, _ := json.MarshalIndent(s, "", "  ")
		logger.Warn(string(data))
	}
}
