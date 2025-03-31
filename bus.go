package dibus

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

type initSubscribersListItem struct {
	sub   SubscriberForBuild
	order int64
}

type BusOptions struct {
	AwaitForGracefulStop time.Duration
}

type ApplicationBus struct {
	subscribers      map[EventName][]Subscriber
	needAwaitStops   []<-chan struct{}
	awaitStopTimeout time.Duration
	ctx              context.Context
	cancel           context.CancelFunc
	signalChannel    chan os.Signal
}

func NewApplicationBus(ctx context.Context, opts BusOptions) *ApplicationBus {
	ctxInner, cancel := context.WithCancel(ctx)
	return &ApplicationBus{
		subscribers:      make(map[EventName][]Subscriber),
		awaitStopTimeout: opts.AwaitForGracefulStop,
		ctx:              ctxInner,
		cancel:           cancel,
		signalChannel:    make(chan os.Signal, 1),
	}
}

func (ab *ApplicationBus) ExecQuery(query Query) Query {
	if subscribers, ok := ab.subscribers[query.Name()]; ok {
		for _, subscriber := range subscribers {
			q := subscriber.ProcessQuery(query)
			q.SetExecuted()
			return q
		}
	}
	return query
}

func (ab *ApplicationBus) ExecCommand(command Command) {
	if subscribers, ok := ab.subscribers[command.Name()]; ok {
		for _, subscriber := range subscribers {
			subscriber.ProcessCommand(command)
			if command.IsStopPropagation() {
				break
			}
		}
	}
}

func (ab *ApplicationBus) Build(providers ...SubscriberProvider) {
	ab.selfSubscribe()
	initList := make([]initSubscribersListItem, len(providers))
	for idx, provider := range providers {
		subscriber := provider(ab.ctx, ab)
		for _, event := range subscriber.SupportedEvents() {
			eventSubscribers := ab.subscribers[event.Name()]
			eventSubscribers = append(eventSubscribers, subscriber)
			ab.subscribers[event.Name()] = eventSubscribers
		}
		initList[idx] = initSubscribersListItem{
			sub:   subscriber,
			order: subscriber.InitOrder(),
		}
		ch := subscriber.IamStopChan()
		if ch != nil {
			ab.needAwaitStops = append(ab.needAwaitStops, ch)
		}
	}
	ab.init(initList)
}

func (ab *ApplicationBus) selfSubscribe() {
	stopEvent := &BusStopCommand{}
	ab.subscribers[stopEvent.Name()] = []Subscriber{ab}
}

func (ab *ApplicationBus) init(list []initSubscribersListItem) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].order < list[j].order
	})
	for _, subscriber := range list {
		subscriber.sub.AfterBusBuild()
	}
}

func (ab *ApplicationBus) Run() {
	signal.Notify(ab.signalChannel, os.Interrupt, syscall.SIGTERM)
	go func(ch <-chan os.Signal) {
		<-ch
		ab.shutdown()
	}(ab.signalChannel)

	<-ab.ctx.Done()
}

func (ab *ApplicationBus) BuildAndRun(providers ...SubscriberProvider) {
	ab.Build(providers...)
	ab.Run()
}

func (ab *ApplicationBus) shutdown() {
	ab.cancel()
	ctx, cancel := context.WithTimeout(context.Background(), ab.awaitStopTimeout)
	defer cancel()
	defer close(ab.signalChannel)

	allStopped := make(chan struct{})

	go func() {
		for _, chnl := range ab.needAwaitStops {
			<-chnl
		}
		allStopped <- struct{}{}
	}()

	select {
	case <-allStopped:
		return
	case <-ctx.Done():
		return
	}
}

func (ab *ApplicationBus) ProcessQuery(query Query) Query {
	return query
}

func (ab *ApplicationBus) ProcessCommand(command Command) {
	switch command.(type) {
	case *BusStopCommand:
		ab.shutdown()
	}
}

func (ab *ApplicationBus) IamStopChan() <-chan struct{} {
	return nil
}
