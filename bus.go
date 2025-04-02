package dibus

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

type initSubscribersListItem struct {
	sub      SubscriberForBuild
	order    int64
	callback func()
}

type multiQueryResultTransport struct {
	idx int
	qr  QueryResult
}

// BusOptions - options for Build Bus
type BusOptions struct {
	AwaitForGracefulStop time.Duration
}

// ApplicationBus - realization of Bus
type ApplicationBus struct {
	subscribers      map[EventName][]Subscriber
	needAwaitStops   []<-chan struct{}
	awaitStopTimeout time.Duration
	ctx              context.Context
	cancel           context.CancelFunc
	signalChannel    chan os.Signal
}

// NewApplicationBus - ApplicationBus constructor
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

// ExecQuery executes a query by finding its corresponding subscribers and processing it
func (ab *ApplicationBus) ExecQuery(query Query) QueryResult {
	if subscribers, ok := ab.subscribers[query.Name()]; ok {
		for _, subscriber := range subscribers {
			qr := subscriber.ProcessQuery(query)
			query.SetExecuted()
			return qr
		}
	}
	return nil
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

func (ab *ApplicationBus) ExecMultiQuery(queries ...Query) []QueryResult {
	cnt := len(queries)
	if cnt == 0 {
		return nil
	}

	resultChan := make(chan multiQueryResultTransport, cnt)

	wg := new(sync.WaitGroup)
	wg.Add(cnt)

	for i, query := range queries {
		go func(idx int, q Query) {
			defer wg.Done()
			resultChan <- multiQueryResultTransport{
				idx: idx,
				qr:  ab.ExecQuery(q),
			}
		}(i, query)
	}

	wg.Wait()
	close(resultChan)

	results := make([]QueryResult, cnt)
	for r := range resultChan {
		results[r.idx] = r.qr
	}

	return results
}

func (ab *ApplicationBus) Build(providers ...SubscriberProvider) {
	ab.selfSubscribe()
	initList := make([]initSubscribersListItem, len(providers))
	for idx, provider := range providers {
		subscriber := provider(ab.ctx, ab)
		opts := subscriber.GetBuildOptions()
		if opts.SupportedEvents != nil && len(opts.SupportedEvents) > 0 {
			for _, event := range opts.SupportedEvents {
				eventSubscribers := ab.subscribers[event.Name()]
				eventSubscribers = append(eventSubscribers, subscriber)
				ab.subscribers[event.Name()] = eventSubscribers
			}
		}

		initList[idx] = initSubscribersListItem{
			sub:      subscriber,
			order:    opts.InitOrder,
			callback: opts.AfterBusBuildCallback,
		}

		ch := opts.ImStoppedChannel
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
		if nil != subscriber.callback {
			subscriber.callback()
		}
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

func (ab *ApplicationBus) ProcessQuery(query Query) QueryResult {
	return query
}

func (ab *ApplicationBus) ProcessCommand(command Command) {
	switch command.(type) {
	case *BusStopCommand:
		ab.shutdown()
	}
}
