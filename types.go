package dibus

import "context"

// EventName - semantic type for name of Event
type EventName string

// SubscriberProvider - type of provider of subscriber. Construct Subscriber
type SubscriberProvider func(ctx context.Context, bus Bus) SubscriberForBuild
