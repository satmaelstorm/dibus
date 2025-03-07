package dibus

import "errors"

var ErrNotQuery = errors.New("try ExecQuery not Query")
var ErrNoSubscriber = errors.New("no subscriber for event")
