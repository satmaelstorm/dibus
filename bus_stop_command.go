package dibus

type BusStopCommand struct {
	AbstractCommand
}

func (bts *BusStopCommand) Name() EventName {
	return FormEventName(bts)
}
