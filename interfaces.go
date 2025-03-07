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

type Subscriber interface {
	SupportedEvents() []Event
	ExecQuery(query Query) Query
	ExecCommand(command Command)
	IamStopChan() <-chan struct{}
}

type Bus interface {
	ExecQuery(query Query) Query
	ExecCommand(command Command)
}
