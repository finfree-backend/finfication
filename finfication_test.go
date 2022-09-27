package finfication

import (
	"log"
	"os"
	"testing"
)

func TestFinfication(t *testing.T) {
	client, err := New(os.Getenv("FINFICATION_CREDENTIALS_FILE"))
	if err != nil {
		log.Println("Error on initializing client. Err:", err)
		t.FailNow()
	}
	defer client.Close()

	exampleNotificationSender, err := client.NewFinficationPublisher(&FinficationOption{
		NotificationType: "INC_PRICE_ALERT_NOTIFICATION",
		TopicName:        os.Getenv("FINFICATION_TOPIC"),
		//HashFunc: func(message *PubSubMessage) string {
		//	return HashString("PARSN:1.87:42.46")
		//},
		HashFunc:              DefaultHashFn,
		EnableMessageOrdering: false,
		OrderingFn:            nil,
	})

	if err != nil {
		log.Println("Error on initializing notification sender. Err:", err)
		t.FailNow()
	}

	data := []*PubSubMessageData{
		{
			Usernames:  []string{"test44"},
			Parameters: map[string]string{"stock": "PARSN", "percentage": "1.87", "price": "42.46"},
		},
	}
	err = exampleNotificationSender.Publish(data)
	if err != nil {
		log.Println("Error on publishing the notification message to Pub/Sub. Err:", err)
		t.FailNow()
	}

	log.Println("Message sent successfully!")

}

func TestSubscribe(t *testing.T) {
	client, err := New(os.Getenv("FINFICATION_CREDENTIALS_FILE"))
	if err != nil {
		log.Println("Error on initializing client. Err:", err)
		t.FailNow()
	}
	defer client.Close()

	if err := client.NewFinficationConsumer(os.Getenv("FINFICATION_SUBSCRIPTION_NAME"), func(message *PubSubMessage) {
		log.Println("Message received:", message)
	}); err != nil {
		log.Println("Err:", err)
	}

}
