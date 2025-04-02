package dibus

type AbstractQueryExecutor struct {
}

func (a *AbstractQueryExecutor) ProcessQuery(query Query) QueryResult {
	panic("implement me")
}

func (a *AbstractQueryExecutor) ProcessCommand(command Command) {
	panic("don't call me")
}

func (a *AbstractQueryExecutor) GetBuildOptions() SubscriberOptions {
	return SubscriberOptions{
		InitOrder:             0,
		ImStoppedChannel:      nil,
		AfterBusBuildCallback: nil,
		SupportedEvents:       nil,
	}
}
