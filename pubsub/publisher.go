package pubsub

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/pubsub"
)

var client *pubsub.Client

func InitPubSub() error {
	var err error
	ctx := context.Background()
	client, err = pubsub.NewClient(ctx, os.Getenv("PUBSUB_PROJECT_ID"))
	return err
}

func PublishOrderEvent(data map[string]interface{}) error {
	topic := client.Topic(os.Getenv("ORDERS_TOPIC"))
	bytes, _ := json.Marshal(data)
	result := topic.Publish(context.Background(), &pubsub.Message{
		Data: bytes,
	})
	_, err := result.Get(context.Background())
	return err
}
