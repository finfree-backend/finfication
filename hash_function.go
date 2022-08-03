package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

type HashFn func(message *PubSubMessage) string

var DefaultHashFn HashFn = func(message *PubSubMessage) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return HashString(message.Hash + timestamp)
}

func HashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
