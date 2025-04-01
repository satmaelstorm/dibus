package dibus

// SubscriberOptions options for create Subscriber.
// ApplicationBus read SubscriberOptions only once in Build process.
// Read-only Value Object.
type SubscriberOptions struct {
	// InitOrder - the order in which callbacks are called
	// Less - earlier
	InitOrder int64

	// AfterBusBuildCallback - callbacks are called after Build
	AfterBusBuildCallback func()

	// ImStoppedChannel - a channel in which a subscriber can notify the bus
	// that it has completed gracefully. Needed for graceful shutdown.
	ImStoppedChannel chan struct{}

	// SupportedEvents - a list of events that the subscriber is willing to process
	SupportedEvents []Event
}
