package main

import (
	"context"
	"fmt"
	"time"
)

type futureInterface interface {
	Cancel()
	Cancelled() bool
	Done() bool
	Result() (interface{}, error)
	ResultUntil(d time.Duration) (interface{}, bool, error)
	doneCallBack(func(interface{}) (interface{}, error)) futureInterface
}

type futureStruct struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	result     chan resultData
}

type resultData struct {
	val interface{}
	err error
}

func (f *futureStruct) Cancel() {
	f.cancelFunc()
}

func (f *futureStruct) Cancelled() bool {
	select {
	case <-f.ctx.Done():
		return true
	default:
		return false
	}
}

func (f *futureStruct) Done() bool {
	select {
	case <-f.ctx.Done():
		return true
	default:
		return false
	}
}

func (f *futureStruct) Result() (interface{}, error) {
	result := <-f.result
	return result.val, result.err
}

func (f *futureStruct) ResultUntil(d time.Duration) (interface{}, bool, error) {
	select {
	case result := <-f.result:
		return result.val, false, result.err
	case <-time.After(d):
		return nil, true, nil
	case <-f.ctx.Done():
	}
	return nil, false, nil
}

func (f *futureStruct) doneCallBack(next func(interface{}) (interface{}, error)) futureInterface {
	nextCtx, nextCancel := context.WithCancel(f.ctx)
	nextFuture := &futureStruct{
		ctx:        nextCtx,
		cancelFunc: nextCancel,
		result:     make(chan resultData),
	}

	go func() {
		val, err := f.Result()
		if !f.Cancelled() && err == nil {
			nextResult, nextErr := next(val)
			nextFuture.result <- resultData{val: nextResult, err: nextErr}
		} else {
			nextFuture.result <- resultData{val: val, err: err}
		}
		close(nextFuture.result)
	}()

	return nextFuture
}

// New creates a new Future that wraps the provided function.
func New(inFunc func() (interface{}, error)) futureInterface {
	ctx, cancel := context.WithCancel(context.Background())
	future := &futureStruct{
		ctx:        ctx,
		cancelFunc: cancel,
		result:     make(chan resultData),
	}

	go func() {
		val, err := inFunc()
		future.result <- resultData{val: val, err: err}
		close(future.result)
	}()

	return future
}

// Main function to test above functions.
func main() {
	var tempVal = 200
	longTimeFunc := func(tempVal int) (int, error) {
		time.Sleep(5 * time.Second)
		return tempVal * 2, nil
	}

	// Start a new instance of the future implementation
	f := New(func() (interface{}, error) {
		return longTimeFunc(tempVal)
	})

	// Checking cancel call
	go func() {
		time.Sleep(2 * time.Second)
		f.Cancel()
	}()

	result, err := f.Result()
	fmt.Println(result, err, f.Cancelled())

	// Checking get call
	g := New(func() (interface{}, error) {
		return longTimeFunc(tempVal)
	})
	gResult, gErr := g.Result()

	fmt.Println(g.Done(), gResult, gErr, g.Cancelled())
}
