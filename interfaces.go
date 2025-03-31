package dibus

type Event interface {
	Name() EventName
}

type Query interface {
	Event
	SetExecuted()
}

type Command interface {
	Event
	IsStopPropagation() bool
}

type SubscriberForBuild interface {
	Subscriber

	// InitOrder return order of call AfterBusBuild
	// less - earlier
	InitOrder() int64

	// AfterBusBuild call when bus is built
	// If you need init some dependencies, you can Query them from other Subscriber's
	AfterBusBuild()

	// SupportedEvents - list of supported events
	SupportedEvents() []Event
}

type Subscriber interface {
	// ProcessQuery - execute query, set query result and return this (or another, if needed) Query
	ProcessQuery(query Query) Query

	// ProcessCommand - execute command
	ProcessCommand(command Command)

	// IamStopChan - returns the channel from which the bus will wait for a signal during a graceful shutdown
	IamStopChan() <-chan struct{}
}

type Bus interface {
	ExecQuery(query Query) Query
	ExecCommand(command Command)
}
