package dibus

import "context"

type EventName string

type SubscriberProvider func(ctx context.Context, bus Bus) Subscriber
