package dibus

// SubscriberOptions options for create Subscriber.
// ApplicationBus read SubscriberOptions only once in Build process.
// Read-only Value Object.
type SubscriberOptions struct {
	// Order - the order in which subscriber called
	// Less - earlier
	Order int64

	// ImStoppedChannel - a channel in which a subscriber can notify the bus
	// that it has completed gracefully. Needed for graceful shutdown.
	ImStoppedChannel chan struct{}

	// SupportedEvents - a list of events that the subscriber is willing to process
	SupportedEvents []Event
}
