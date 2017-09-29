package goref

import (
	"encoding/json"
	"sync"
	"test/consul2/lbbconsul"
	"time"

	"github.com/dongxiaozhen/lbbutil"
	"github.com/donnie4w/go-logger/logger"
)

var ConsulKey = make(map[string]bool)
var mutex sync.Mutex

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
		for k := range ConsulKey {
			if d, ok := s[k]; ok {
				count := lbbutil.Int64ToString(d.Count)
				lbbconsul.GConsulClient.PutKV(k, []byte(count))
			}
		}
	}
}

func SetConsulKey(key string) {
	mutex.Lock()
	ConsulKey[key] = true
	mutex.Unlock()
}

func DelConsulKey(key string) {
	mutex.Lock()
	delete(ConsulKey, key)
	mutex.Unlock()
}
