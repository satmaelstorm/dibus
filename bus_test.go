package dibus

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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

	bts.Assert().Equal(3, len(bus.subscribers))
}

func (bts *busTestSuite) TestQuery() {
	bus := NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
	bus.Build(provideTSubscriber)

	query := &tQuery{}
	queryResult, _ := ExecQueryWrapper[int](bus, query)

	bts.Assert().Equal(10, queryResult)
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

	bts.Assert().Equal(11, results[2].(int))
	bts.Assert().Equal(20, results[1].(int))
	bts.Assert().Equal(110, results[0].(int))
}

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
		InitOrder: 0,
		AfterBusBuildCallback: func() {
			t.val = 10
		},
		ImStoppedChannel: nil,
		SupportedEvents:  []Event{new(tQuery), new(tCommand)},
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
