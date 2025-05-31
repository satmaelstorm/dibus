# diBus
Dependency Injection Bus

[![License](https://img.shields.io/badge/license-MIT-blue)]()

## Table of Contents
- [Core Concepts](#core-concepts)
- [Bus](#bus)
- [Subscribers](#subscribers)
- [Events](#events)
    - [Commands](#commands)
    - [Queries](#queries)
- [Building the Bus](#how-to-build-bus)
- [Execution](#executing-commands-and-queries)
- [Lifecycle](#built-in-commands)
- [Shutdown](#shutdown-process)
- [Dependencies](#what-about-dependencies)

A fresh approach to dependency injection: everything is simply Subscribers and Events.
Following the Command-Query Separation principle:
* Commands perform actions (no return value)
* Queries retrieve data (no side effects)

The event bus runs entirely within your application—no need for external systems like RabbitMQ or Kafka.

Keep your code decoupled without the overhead of distributed systems.

At the core of this approach are three simple concepts that work together to keep your code clean and decoupled:
* Bus
* Subscribers
* Events (Queries and Commands)

All code examples are from the `bus_test.go` test — there you can check how it all works.

## Bus
Creating a bus is straightforward - just provide a context and configuration:
```go
bus := dibus.NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
```

Current configuration options include:
* AwaitForGracefulStop - timeout duration for graceful shutdown

Always interact with the bus through the `dibus.Bus` interface in your application code.

## Subscribers
A subscriber must implement the `dibus.SubscriberForBuild` interface. For example:
```go
type tSubscriber struct {
	val int
	bus Bus
}

func (t *tSubscriber) ProcessQuery(query dibus.Query) dibus.QueryResult {
	switch q := query.(type) {
	case *tQuery:
		return t.val + q.additional
	case *tQueryTSubscriber:
		return t
	default:
	}
}

func (t *tSubscriber) ProcessCommand(command dibus.Command) {
	switch c := command.(type) {
	case *tCommand:
		t.val = c.Val
	default:
	}
}

func (t *tSubscriber) GetBuildOptions() dibus.SubscriberOptions {
	return dibus.SubscriberOptions{
		Order:            0,
		ImStoppedChannel: nil,
		SupportedEvents:  []Event{new(tQuery), new(tCommand), new(tQueryTSubscriber)},
	}
}
```

`SubscriberOptions` contains the subscriber's configuration:
* `Order` - execution order priority (lower values execute first)
* `ImStoppedChannel` - channel for notifying the Bus about graceful shutdown completion (used for Graceful Shutdown)
* `SupportedEvents` - list of supported event types

## Events 
Events are simple DTOs (Data Transfer Objects), but they must implement their respective interfaces.
### Command
The `dibus.Command` interface requires:
* `Name() EventName` - returns the event name
* `IsStopPropagation() bool` - determines whether to stop processing subsequent commands

For easier command creation, you can use:
* The `dibus.AbstractCommand` base type
* The helper function `dibus.FormEventName`

Example:
```go
type tCommand struct {
	dibus.AbstractCommand
	Val int
}

func (t *tCommand) Name() dibus.EventName {
	return dibus.FormEventName(t)
}
```

### Query
The `dibus.Query` interface requires:
* `Name() EventName` - returns the event name
* `SetExecuted()` - used by the bus to mark the query as executed

For easier query creation, you can use:
* The `dibus.AbstractQuery` base type
* The helper function `dibus.FormEventName`

Example:
```go
type tQuery struct {
	dibus.AbstractQuery
	additional int
}

func (t *tQuery) Name() dibus.EventName {
	return dibus.FormEventName(t)
}
```

## How to Build Bus
Now that you have your Bus and Subscribers ready, it's time to assemble the bus and optionally launch your application.
### Preparing Subscriber Providers
Before building the bus, you need to create a Provider for each subscriber. A Provider is a function matching the `dibus.SubscriberProvider` type:
```go
func provideTSubscriber(ctx context.Context, bus dibus.Bus) diBus.SubscriberForBuild {
	return &tSubscriber{
		val: 0,
		bus: bus,
	}
}
```

### Building the Bus
You can now assemble the bus with your providers:
```go
bus := dibus.NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
bus.Build(provideTSubscriberWithInitialization, provideTSubscriber)
```

### Framework Mode (BuildAndRun)
If you want to both build the bus and immediately launch your application (dibus framework mode):
```go
bus := dibus.NewApplicationBus(context.Background(), BusOptions{AwaitForGracefulStop: time.Second})
bus.BuildAndRun(provideTSubscriberWithInitialization, provideTSubscriber)
```

#### Key Difference
The only difference between `Build` and `BuildAndRun` is that `BuildAndRun` will block the main thread until receiving an interrupt signal.

Notes:
* All providers are registered during the build phase
* `BuildAndRun` is recommended for typical application use
* The bus becomes operational immediately after building

## Executing Commands and Queries
### Command Execution
Execute commands using: 
```go
bus.ExecCommand(dibus.Command)
```
### Query Execution Options
You have several ways to execute queries:
1. Basic execution:
```go
bus.ExecQuery(dibus.Query) dibus.QueryResult
```

2. Type-safe helper (recommended):
```go
dibus.ExecQueryWrapper[T any](bus dibus.Bus, q dibus.Query) (T, error)
```

For Example:
```go
queryResult, _ := ExecQueryWrapper[int](bus, query)
t.dep, _ = ExecQueryWrapper[*tSubscriber](t.bus, &tQueryTSubscriber{})
```

3. Parallel execution of multiple queries:
```go
dibus.ExecMultiQuery(queries ...Query) []QueryResult
```

* The results slice order exactly matches the input queries order
* All queries execute concurrently for better performance

### Key Notes:
* Commands are fire-and-forget operations
* Queries return results and support type-safe unwrapping
* Parallel execution maintains request/response ordering
* Always check errors in production code (omitted here for brevity)

## Built-in Commands
The Bus utilizes the following system commands during its lifecycle:
### `dibus.BusInitializedCommand`
* **Triggered when:** The Bus completes its initialization
* **Purpose:** Subscribe to this command to perform additional application initialization tasks

### `dibus.BusStopCommand`
* **Special note:** This is a command the Bus listens to (not publishes)
* **Purpose:** Programmatically stop the Bus instead of using interrupt signals
* **Framework mode only:** Only relevant when using BuildAndRun framework mode

Note: These commands should not be modified or reimplemented - use them as provided by the library.

## Shutdown Process
*(Applies only to Framework Mode)*
1. **Graceful Stop Attempt**
   * Waits for all subscribers with `ImStoppedChannel != nil` to signal completion
   * Subscribers should send a notification via their `ImStoppedChannel` when done processing
2. **Timeout Enforcement**
   * If some subscribers fail to signal completion within the configured AwaitForGracefulStop period:
     * The Bus forcibly terminates
     * No additional waiting occurs
3. **Guaranteed Termination**
   * Regardless of subscriber state, the Bus always shuts down after the timeout
   * Ensures the application won't hang during exit

## What About Dependencies?
While the bus architecture naturally promotes decoupling through events, there are cases where explicit dependencies between subscribers are unavoidable. Here's the recommended approach:
1. **Subscribe to Bus Initialization:**
The dependent subscriber should listen for `dibus.BusInitializedCommand` - this ensures all components are registered before requesting dependencies.
2. **Request Dependencies via Queries:** 
When handling the initialization command, the subscriber sends a specific Query to obtain the required dependency.
3. **Provide Self-Reference:**
The provider subscriber handles this query and returns a reference to itself in the response.

**Lazy Resolution:** Dependencies are resolved only after full bus initialization.

This pattern preserves the flexibility of event-driven architecture while allowing necessary explicit dependencies when absolutely required. The key principle remains: prefer event-based communication, but use dependency queries sparingly when direct interaction is unavoidable.