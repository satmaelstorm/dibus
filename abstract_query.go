package dibus

type AbstractQuery struct {
	executed bool
	name     EventName
}

func (a *AbstractQuery) Name() EventName {
	panic("implement me")
}

func (a *AbstractQuery) SetExecuted() {
	a.executed = true
}

func (a *AbstractQuery) IsExecuted() bool {
	return a.executed
}
