package dibus

type BusInitializedCommand struct {
	AbstractCommand
}

func (bi *BusInitializedCommand) Name() EventName {
	return FormEventName(bi)
}
