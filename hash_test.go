package main

import (
	"log"
	"testing"
)

func TestHashString(t *testing.T) {
	hash := HashString("Hello World!")
	log.Println("Hash:", hash)
	expectedHash := "ed076287532e86365e841e92bfc50d8c"
	if hash != expectedHash {
		log.Println("Expected:", expectedHash+".", "Found:", hash)
		t.FailNow()
	}
}

func TestDefaultHashFn(t *testing.T) {
	hash := DefaultHashFn(&PubSubMessage{NotificationType: "EXAMPLE_NOTIFICATION_TYPE"})
	log.Println("Hash:", hash)
}
