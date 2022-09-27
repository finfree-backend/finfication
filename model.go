package finfication

type PubSubMessage struct {
	NotificationType string             `json:"nt"`
	Hash             string             `json:"h"`
	Data             *PubSubMessageData `json:"d"`
}

type PubSubMessageData struct {
	Usernames          []string          `json:"u"`
	Parameters         map[string]string `json:"p"`
	OptionalParameters map[string]string `json:"op"`
}

type FinficationOption struct {
	NotificationType      string
	TopicName             string
	HashFunc              HashFn
	EnableMessageOrdering bool
	OrderingFn            OrderingFn
}

type FinficationError struct {
	Message string
}

func (f *FinficationError) Error() string {
	return f.Message
}
