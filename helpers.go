package dibus

import (
	"reflect"
)

func FormEventName(caller any) EventName {
	ref := reflect.TypeOf(caller).Elem()
	return EventName(ref.PkgPath() + "/" + ref.Name())
}

func ExecQueryWrapper[T any](bus Bus, q Query) (T, error) {
	qr := bus.ExecQuery(q)
	if qr == nil {
		return *new(T), ErrSubscriberResultTypeMismatch
	}
	r, ok := qr.(T)
	if !ok {
		return r, ErrSubscribersForQueryNotFound
	}
	return r, nil
}
