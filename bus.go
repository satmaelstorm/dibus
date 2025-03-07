package dibus

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ApplicationBus struct {
	subscribers      map[EventName][]Subscriber
	needAwaitStops   []<-chan struct{}
	awaitStopTimeout time.Duration
	ctx              context.Context
	cancel           context.CancelFunc
}

func NewApplicationBus(ctx context.Context, awaitForGracefulStop time.Duration) *ApplicationBus {
	ctxInner, cancel := context.WithCancel(ctx)
	return &ApplicationBus{
		subscribers:      make(map[EventName][]Subscriber),
		awaitStopTimeout: awaitForGracefulStop,
		ctx:              ctxInner,
		cancel:           cancel,
	}
}

func (ab *ApplicationBus) ExecQuery(query Query) Query {
	if subscribers, ok := ab.subscribers[query.Name()]; ok {
		for _, subscriber := range subscribers {
			q := subscriber.ExecQuery(query)
			q.SetExecuted()
			return q
		}
	}
	return query
}

func (ab *ApplicationBus) ExecCommand(command Command) {
	if subscribers, ok := ab.subscribers[command.Name()]; ok {
		for _, subscriber := range subscribers {
			subscriber.ExecCommand(command)
			if command.IsStopPropagation() {
				break
			}
		}
	}
}

func (ab *ApplicationBus) Build(providers ...SubscriberProvider) {
	for _, provider := range providers {
		subscriber := provider(ab.ctx, ab)
		for _, event := range subscriber.SupportedEvents() {
			eventSubscribers := ab.subscribers[event.Name()]
			eventSubscribers = append(eventSubscribers, subscriber)
			ab.subscribers[event.Name()] = eventSubscribers
		}
		ch := subscriber.IamStopChan()
		if ch != nil {
			ab.needAwaitStops = append(ab.needAwaitStops, ch)
		}
	}
}

func (ab *ApplicationBus) Run() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func(ch <-chan os.Signal) {
		<-ch
		ab.cancel()

		ctx, cancel := context.WithTimeout(context.Background(), ab.awaitStopTimeout)
		defer cancel()

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
	}(signalChannel)

	<-ab.ctx.Done()
}

func (ab *ApplicationBus) BuildAndRun(providers ...SubscriberProvider) {
	ab.Build(providers...)
	ab.Run()
}
