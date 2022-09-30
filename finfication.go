package finfication

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
)

type ConsumerFn func(message *PubSubMessage)

type Client struct {
	pubsubClient *pubsub.Client
}

func New(credentialFilePath string) (*Client, error) {
	file, err := ioutil.ReadFile(credentialFilePath)
	if err != nil {
		return nil, err
	}

	return NewWithBytes(file)
}

func NewWithBytes(bytes []byte) (*Client, error) {
	credentials, err := google.CredentialsFromJSON(context.Background(), bytes, pubsub.ScopePubSub)
	if err != nil {
		return nil, err
	}

	var finfication Client
	finfication.pubsubClient, err = pubsub.NewClient(context.Background(), credentials.ProjectID, option.WithCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return &finfication, nil
}

func (f *Client) NewFinficationPublisher(opt *FinficationOption) (*NewFinficationPublisher, error) {
	if opt.NotificationType == "" {
		return nil, &FinficationError{Message: "'NotificationType' parameter can not be blank"}
	}

	if opt.TopicName == "" {
		return nil, &FinficationError{Message: "'TopicName' parameter can not be blank"}
	}

	fs := NewFinficationPublisher{
		notificationType:      opt.NotificationType,
		topic:                 f.pubsubClient.Topic(opt.TopicName),
		hashFn:                opt.HashFunc,
		EnableMessageOrdering: opt.EnableMessageOrdering,
		OrderingFn:            opt.OrderingFn,
	}

	if fs.hashFn == nil {
		fs.hashFn = DefaultHashFn
	}

	if fs.EnableMessageOrdering {
		// Enable message ordering in *pubsub.Topic
		fs.topic.EnableMessageOrdering = true
		if fs.OrderingFn == nil {
			// If message ordering is enabled and no ordering function set, default one will be set
			fs.OrderingFn = DefaultOrderingFn
		}
	}
	return &fs, nil
}

func (f *Client) NewFinficationConsumer(subscriptionName string, fn ConsumerFn) error {
	return f.pubsubClient.Subscription(subscriptionName).Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		defer message.Ack()
		var msg PubSubMessage
		if err := json.Unmarshal(message.Data, &msg); err != nil {
			log.Println("Error on unmarshalling Pub/Sub message. Err:", err, ". Data:", string(message.Data))
			return
		}
		fn(&msg)
	})
}

func (f *Client) Close() {
	f.pubsubClient.Close()
}

type NewFinficationPublisher struct {
	topic                 *pubsub.Topic
	notificationType      string
	hashFn                HashFn
	EnableMessageOrdering bool
	OrderingFn            OrderingFn
}

func (fs *NewFinficationPublisher) Publish(data *PubSubMessageData) error {
	pubSubMessage := PubSubMessage{
		NotificationType: fs.notificationType,
		Data:             data,
	}
	pubSubMessage.Hash = fs.hashFn(&pubSubMessage)

	byteData, err := json.Marshal(&pubSubMessage)
	if err != nil {
		return &FinficationError{Message: "Could not marshall the data you passed. Err:" + err.Error()}
	}
	msg := pubsub.Message{Data: byteData}
	if fs.EnableMessageOrdering {
		msg.OrderingKey = fs.OrderingFn(&pubSubMessage)
	}

	_, err = fs.topic.Publish(context.Background(), &msg).Get(context.Background())
	return err
}
