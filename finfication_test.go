package finfication

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
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
		NotificationType:      "EXAMPLE_NOTIFICATION-v2",
		TopicName:             "finfication",
		HashFunc:              nil,
		EnableMessageOrdering: false,
		OrderingFn:            nil,
	})

	if err != nil {
		log.Println("Error on initializing notification sender. Err:", err)
		t.FailNow()
	}

	data := []*PubSubMessageData{
		{
			Username:           "ErtugrulAcar",
			Parameters:         map[string]string{"name": "Ertugrul", "last_name": "Acar"},
			OptionalParameters: map[string]string{"example_key": "example_value"},
		},
		{
			Username:           "cemunuvar",
			Parameters:         map[string]string{"name": "Cem", "last_name": "Unuvar"},
			OptionalParameters: map[string]string{"example_key": "example_value"},
		},
		{
			Username:           "ahmetakil",
			Parameters:         map[string]string{"name": "Ahmet", "last_name": "Akıl"},
			OptionalParameters: map[string]string{"example_key": "example_value"},
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

	client.pubsubClient.Subscription("finfication-sub").Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		message.Ack()
		go func(messagex *pubsub.Message) {
			var msg PubSubMessage
			if err := json.Unmarshal(messagex.Data, &msg); err != nil {
				log.Println("Error on unmarshalling data. Err:", err)
				return
			}

			log.Println("Msg received:", msg)
		}(message)
	})

}
