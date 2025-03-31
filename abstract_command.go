package dibus

import "reflect"

type AbstractCommand struct {
	stopPropagation bool
	name            EventName
}

func (a *AbstractCommand) Name() EventName {
	if a.name == "" {
		ref := reflect.TypeOf(a)
		a.name = EventName(ref.PkgPath() + "_" + ref.Name())
	}
	return a.name
}

func (a *AbstractCommand) IsStopPropagation() bool {
	return a.stopPropagation
}

func (a *AbstractCommand) StopPropagation() {
	a.stopPropagation = true
}
