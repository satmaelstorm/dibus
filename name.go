package dibus

import "reflect"

func FormEventName(caller any) EventName {
	ref := reflect.TypeOf(caller).Elem()
	return EventName(ref.PkgPath() + "/" + ref.Name())
}
