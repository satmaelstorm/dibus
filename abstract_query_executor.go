package dibus

type AbstractQueryExecutor struct {
}

func (a *AbstractQueryExecutor) InitOrder() int64 {
	return 0
}

func (a *AbstractQueryExecutor) AfterBusBuild() {
	//pass
}

func (a *AbstractQueryExecutor) SupportedEvents() []Event {
	panic("implement me")
}

func (a *AbstractQueryExecutor) ProcessQuery(query Query) Query {
	panic("implement me")
}

func (a *AbstractQueryExecutor) ProcessCommand(command Command) {
	panic("don't call me")
}

func (a *AbstractQueryExecutor) IamStopChan() <-chan struct{} {
	return nil
}
