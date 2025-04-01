package dibus

type AbstractCommand struct {
	stopPropagation bool
	name            EventName
}

func (a *AbstractCommand) Name() EventName {
	panic("implement me")
}

func (a *AbstractCommand) IsStopPropagation() bool {
	return a.stopPropagation
}

func (a *AbstractCommand) StopPropagation() {
	a.stopPropagation = true
}
