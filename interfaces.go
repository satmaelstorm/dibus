package dibus

// Event - just Event
type Event interface {
	Name() EventName
}

// Query - semantic type. Query can perform read-only operations
// Processor of Query MUST return this Query and MAY set in Query result
type Query interface {
	Event
	SetExecuted()
}

// Command - semantic type. Commands can perform write (change) operations
// and returns nothing.
type Command interface {
	Event
	IsStopPropagation() bool
}

// SubscriberForBuild - uses for Build only
type SubscriberForBuild interface {
	Subscriber

	// GetBuildOptions returns options for Bus.Build()
	GetBuildOptions() SubscriberOptions
}

// Subscriber - it is Subscriber and nothing more
type Subscriber interface {
	// ProcessQuery - execute query, set query result and return this (or another, if needed) Query
	ProcessQuery(query Query) Query

	// ProcessCommand - execute command
	ProcessCommand(command Command)
}

// Bus - interface for ApplicationBus
type Bus interface {
	ExecQuery(query Query) Query
	ExecCommand(command Command)
	ExecMultiQuery(queries ...Query) []Query
}
