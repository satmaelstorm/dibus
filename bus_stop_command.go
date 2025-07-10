package dibus

type BusStopCommand struct {
	AbstractCommand
	Err error
}

func (bts *BusStopCommand) Name() EventName {
	return FormEventName(bts)
}
