package gcp

//This package will contain all helper function to deal with
// the google pubsub service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

//CreateTopic Create a topic if it does not exist. Otherwise, return the current existin one.
func CreateTopic(topic string, client *pubsub.Client) (*pubsub.Topic, error) {
	ctx := context.Background()

	// Create a topic to subscribe to.
	t := client.Topic(topic)
	exists, err := t.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize topic '%s'. %v", t, err)
	}

	if exists {
		return t, nil
	}

	t, err = client.CreateTopic(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic '%s'. %v", t, err)
	}
	return t, nil
}

//ListSubs will return the available subscriptions
func ListSubs(client *pubsub.Client) ([]*pubsub.Subscription, error) {
	ctx := context.Background()
	// [START pubsub_list_subscriptions]
	var subs []*pubsub.Subscription
	it := client.Subscriptions(ctx)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list subscriptions. %v", err)
		}
		subs = append(subs, s)
	}
	// [END pubsub_list_subscriptions]
	return subs, nil
}

//CreateSub Create a subscription if it does't exist
func CreateSub(client *pubsub.Client, subName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	//first list to see if the subscription exists
	// Get Subscriptions
	var err error
	var subs []*pubsub.Subscription

	subs, err = ListSubs(client)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions. %v", err)
	}
	// List available subscriptions
	for _, sub := range subs {
		if strings.HasSuffix(sub.String(), subName) {
			// return existing subscriptions mathcing the name
			return sub, nil
		}
	}
	ctx := context.Background()

	var sub *pubsub.Subscription
	//TODO: For now, leave the AckDeadline to 20 seconds. Need to make it configurable
	sub, err = client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
		Topic:             topic,
		RetentionDuration: 48 * time.Hour,
		AckDeadline:       600 * time.Second,
	})

	if err != nil {
		return nil, fmt.Errorf("failed create subscription. %v", err)
	}
	// [END pubsub_create_pull_subscription]
	return sub, nil
}
