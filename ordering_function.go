package finfication

import (
	"strconv"
	"time"
)

type OrderingFn func(message *PubSubMessage) string

// Default ordering function returns the current time as nano second
var DefaultOrderingFn OrderingFn = func(message *PubSubMessage) string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
