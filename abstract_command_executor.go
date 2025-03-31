package dibus

type AbstractCommandExecutor struct {
}

func (a *AbstractCommandExecutor) InitOrder() int64 {
	return 0
}

func (a *AbstractCommandExecutor) AfterBusBuild() {
	//pass
}

func (a *AbstractCommandExecutor) SupportedEvents() []Event {
	panic("implement me")
}

func (a *AbstractCommandExecutor) ProcessQuery(query Query) Query {
	panic("don't call me")
}

func (a *AbstractCommandExecutor) ProcessCommand(command Command) {
	panic("implement me")
}

func (a *AbstractCommandExecutor) IamStopChan() <-chan struct{} {
	return nil
}
