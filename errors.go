package dibus

import "errors"

var ErrSubscribersForQueryNotFound = errors.New("subscribers for query not found")
var ErrSubscriberResultTypeMismatch = errors.New("subscriber result type mismatch")
