package dibus

type AbstractCommandExecutor struct {
}

func (a *AbstractCommandExecutor) ProcessQuery(query Query) Query {
	panic("don't call me")
}

func (a *AbstractCommandExecutor) ProcessCommand(command Command) {
	panic("implement me")
}

func (a *AbstractCommandExecutor) GetBuildOptions() SubscriberOptions {
	return SubscriberOptions{
		InitOrder:             0,
		ImStoppedChannel:      nil,
		AfterBusBuildCallback: nil,
		SupportedEvents:       nil,
	}
}
