package cqrs

import (
	"context"
	"fmt"
	"reflect"
)

// CommandBus dispatches commands to handlers.
type CommandBus interface {
	Register(cmd any, handler any)
	Dispatch(ctx context.Context, cmd any) error
}

// InMemoryCommandBus is a simple in-memory implementation.
type InMemoryCommandBus struct {
	handlers map[reflect.Type]reflect.Value
}

func NewCommandBus() *InMemoryCommandBus {
	return &InMemoryCommandBus{handlers: make(map[reflect.Type]reflect.Value)}
}

func (b *InMemoryCommandBus) Register(cmd any, handler any) {
	cmdType := reflect.TypeOf(cmd)
	// handler should be func(ctx, cmd) error
	b.handlers[cmdType] = reflect.ValueOf(handler)
}

func (b *InMemoryCommandBus) Dispatch(ctx context.Context, cmd any) error {
	cmdType := reflect.TypeOf(cmd)
	handler, ok := b.handlers[cmdType]
	if !ok {
		return fmt.Errorf("no handler registered for command: %s", cmdType)
	}

	args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(cmd)}
	results := handler.Call(args)

	if len(results) > 0 {
		errVal := results[0]
		if !errVal.IsNil() {
			return errVal.Interface().(error)
		}
	}
	return nil
}

// QueryBus dispatches queries to handlers.
type QueryBus interface {
	Register(query any, handler any)
	Dispatch(ctx context.Context, query any) (any, error)
}

type InMemoryQueryBus struct {
	handlers map[reflect.Type]reflect.Value
}

func NewQueryBus() *InMemoryQueryBus {
	return &InMemoryQueryBus{handlers: make(map[reflect.Type]reflect.Value)}
}

func (b *InMemoryQueryBus) Register(query any, handler any) {
	queryType := reflect.TypeOf(query)
	b.handlers[queryType] = reflect.ValueOf(handler)
}

func (b *InMemoryQueryBus) Dispatch(ctx context.Context, query any) (any, error) {
	queryType := reflect.TypeOf(query)
	handler, ok := b.handlers[queryType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for query: %s", queryType)
	}

	args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(query)}
	results := handler.Call(args)

	// Result is (Response, error)
	var res any
	if len(results) > 0 {
		res = results[0].Interface()
	}

	if len(results) > 1 {
		errVal := results[1]
		if !errVal.IsNil() {
			return nil, errVal.Interface().(error)
		}
	}
	return res, nil
}
