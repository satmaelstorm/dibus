package dibus

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type busTestSuite struct {
	suite.Suite
}

func TestBus(t *testing.T) {
	suite.Run(t, new(busTestSuite))
}

func (bts *busTestSuite) TestBuild() {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriber)

	bts.Assert().Equal(4, len(bus.subscribers))
}

func (bts *busTestSuite) TestQuery() {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriber)

	query := &tQuery{additional: 11}
	queryResult, _ := ExecQueryWrapper[int](bus, query)

	bts.Assert().Equal(11, queryResult)
}

func (bts *busTestSuite) TestCommand() {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriber)

	command := &tCommand{
		Val: 1,
	}
	bus.ExecCommand(command)
	query := &tQuery{}
	queryResult, _ := ExecQueryWrapper[int](bus, query)

	bts.Assert().Equal(1, queryResult)
}

func (bts *busTestSuite) TestMultiQuery() {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriber)

	query1 := &tQuery{additional: 1}
	query2 := &tQuery{additional: 10}
	query3 := &tQuery{additional: 100}

	results := bus.ExecMultiQuery(query3, query2, query1)

	bts.Assert().Equal(1, results[2].(int))
	bts.Assert().Equal(10, results[1].(int))
	bts.Assert().Equal(100, results[0].(int))
}

func (bts *busTestSuite) TestBusAfterInitialization() {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriberWithInitialization, provideTSubscriber)

	subscribers, ok := bus.subscribers[(&BusInitializedCommand{}).Name()]
	bts.Assert().True(ok)
	bts.Assert().NotNil(subscribers)
	bts.Assert().Equal(1, len(subscribers))
	bts.Assert().IsType(&tSubscriberWithInitialization{}, subscribers[0])
	bts.Assert().IsType(&tSubscriber{}, subscribers[0].(*tSubscriberWithInitialization).dep)
}

// Fuzz tests

func FuzzQuery(f *testing.F) {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriber)
	f.Add(10)
	f.Fuzz(func(t *testing.T, i int) {
		query := &tQuery{additional: i}
		queryResult, err := ExecQueryWrapper[int](bus, query)
		if err != nil {
			t.Errorf("err is not nil in FuzzQuery\n")
		}
		if queryResult != i {
			t.Errorf("expected %d, but got %d in FuzzQuery\n", i, queryResult)
		}
	})
}

// Benchmarks

func BenchmarkQuery(b *testing.B) {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriber)
	query := &tQuery{additional: 10}
	var qr int
	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qr, err = ExecQueryWrapper[int](bus, query)
		if err != nil || qr != 10 {
			b.Errorf("Unexcpected results")
		}
	}
}

func BenchmarkBuild1Subscriber(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
		bus.Build(provideTSubscriber)
	}
}

func BenchmarkBuild10Subscriber(b *testing.B) {
	providers := make([]SubscriberProvider, 10)
	for i := 0; i < 10; i++ {
		providers[i] = provideTSubscriber
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
		bus.Build(providers...)
	}
}

func BenchmarkBuild100Subscriber(b *testing.B) {
	providers := make([]SubscriberProvider, 100)
	for i := 0; i < 100; i++ {
		providers[i] = provideTSubscriber
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
		bus.Build(providers...)
	}
}

// make benchmark for tQuery

// Test Subscribers and Events

func provideTSubscriber(ctx context.Context, bus Bus) SubscriberForBuild {
	return &tSubscriber{
		val: 0,
		bus: bus,
	}
}

type tSubscriber struct {
	val int
	bus Bus
}

func (t *tSubscriber) ProcessQuery(query Query) QueryResult {
	switch q := query.(type) {
	case *tQuery:
		return t.val + q.additional
	case *tQueryTSubscriber:
		return t
	default:
		panic("unknown query type")
	}
}

func (t *tSubscriber) ProcessCommand(command Command) {
	switch c := command.(type) {
	case *tCommand:
		t.val = c.Val
	default:
		panic("unknown command type")
	}
}

func (t *tSubscriber) GetBuildOptions() SubscriberOptions {
	return SubscriberOptions{
		Order:            0,
		ImStoppedChannel: nil,
		SupportedEvents:  []Event{new(tQuery), new(tCommand), new(tQueryTSubscriber)},
	}
}

type tQuery struct {
	AbstractQuery
	additional int
}

func (t *tQuery) Name() EventName {
	return FormEventName(t)
}

type tCommand struct {
	AbstractCommand
	Val int
}

func (t *tCommand) Name() EventName {
	return FormEventName(t)
}

type tQueryTSubscriber struct {
	AbstractQuery
}

func (t *tQueryTSubscriber) Name() EventName {
	return FormEventName(t)
}

func provideTSubscriberWithInitialization(ctx context.Context, bus Bus) SubscriberForBuild {
	return &tSubscriberWithInitialization{
		bus: bus,
	}
}

type tSubscriberWithInitialization struct {
	dep *tSubscriber
	bus Bus
}

func (t *tSubscriberWithInitialization) ProcessQuery(query Query) QueryResult {
	switch query.(type) {
	default:
		panic("can't handle queries")
	}
}

func (t *tSubscriberWithInitialization) ProcessCommand(command Command) {
	switch command.(type) {
	case *BusInitializedCommand:
		t.dep, _ = ExecQueryWrapper[*tSubscriber](t.bus, &tQueryTSubscriber{})
	default:
		panic("unknown command type")
	}
}

func (t *tSubscriberWithInitialization) GetBuildOptions() SubscriberOptions {
	return SubscriberOptions{
		Order:            100,
		ImStoppedChannel: nil,
		SupportedEvents:  []Event{new(BusInitializedCommand)},
	}
}
