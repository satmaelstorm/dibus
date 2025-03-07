package dibus

type AbstractQueryExecutor struct {
}

func (a *AbstractQueryExecutor) SupportedEvents() []Event {
	panic("implement me")
}

func (a *AbstractQueryExecutor) ExecQuery(query Query) Query {
	panic("implement me")
}

func (a *AbstractQueryExecutor) ExecCommand(command Command) {
	panic("don't call me")
}

func (a *AbstractQueryExecutor) IamStopChan() <-chan struct{} {
	return nil
}

type AbstractCommandExecutor struct {
}

func (a *AbstractCommandExecutor) SupportedEvents() []Event {
	panic("implement me")
}

func (a *AbstractCommandExecutor) ExecQuery(query Query) Query {
	panic("don't call me")
}

func (a *AbstractCommandExecutor) ExecCommand(command Command) {
	panic("implement me")
}

func (a *AbstractCommandExecutor) IamStopChan() <-chan struct{} {
	return nil
}

type AbstractCommand struct {
	stopPropagation bool
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

type AbstractQuery struct {
	executed bool
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
