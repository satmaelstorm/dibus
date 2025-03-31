package dibus

import "reflect"

type AbstractQuery struct {
	executed bool
	name     EventName
}

func (a *AbstractQuery) Name() EventName {
	if a.name == "" {
		ref := reflect.TypeOf(a)
		a.name = EventName(ref.PkgPath() + "/" + ref.Name())
	}
	return a.name
}

func (a *AbstractQuery) SetExecuted() {
	a.executed = true
}

func (a *AbstractQuery) IsExecuted() bool {
	return a.executed
}
