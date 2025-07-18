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
		return *new(T), ErrSubscribersForQueryNotFound
	}
	r, ok := qr.(T)
	if !ok {
		return r, ErrSubscriberResultTypeMismatch
	}
	return r, nil
}

func MustExecQueryWrapper[T any](bus Bus, q Query) T {
	qr := bus.ExecQuery(q)
	if qr == nil {
		panic(ErrSubscribersForQueryNotFound)
	}
	r, ok := qr.(T)
	if !ok {
		panic(ErrSubscriberResultTypeMismatch)
	}
	return r
}
