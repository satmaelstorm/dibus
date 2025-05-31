package dibus

type AbstractCommandExecutor struct {
}

func (a *AbstractCommandExecutor) ProcessQuery(query Query) QueryResult {
	panic("don't call me")
}

func (a *AbstractCommandExecutor) ProcessCommand(command Command) {
	panic("implement me")
}

func (a *AbstractCommandExecutor) GetBuildOptions() SubscriberOptions {
	return SubscriberOptions{
		Order:            0,
		ImStoppedChannel: nil,
		SupportedEvents:  nil,
	}
}
