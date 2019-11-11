// Package semaphore provides a simple implementation of a semaphore (hence the name, duh).
// Semaphores may be useful to enforce a upper bound on how much
// concurrent units executes some algorithm that is very heavy on memory and/or CPU.
//
// One example is when you have an http server that can accept hundreds
// of concurrent requests but you have some critical section of your code that
// cant handle that much concurrency at once (but you don't want to serialize either).
//
// There are other ways to handle this too using channels, which approach to use will
// be your judgement call. If you choose the semaphore here you got one.
package semaphore

import (
	"context"
	"fmt"
	"sync"
)

// S is a semaphore instance. Always use the New function
// to create semaphores. Using a unproperly initialized
// semaphore may cause mayhem on your code.
type S chan struct{}

// Release is used to release a previous call to S.Acquire
type Release func()

// New creates a new semaphore with the given size.
// Passing 0 as size is a moronic programming mistake and will
// result in a panic due to its moronicness.
func New(size uint) S {
	if size == 0 {
		panic("semaphore.New:cant create a semaphore with size 0")
	}
	return S(make(chan struct{}, size))
}

// Acquire will acquire the semaphore, if the semaphore is
// already full Acquire will wait until some other goroutine
// releases the semaphore or until the given context timeouts.
//
// If the context timeouts or is cancelled a non nil error is returned.
// On success a release function is returned, this function can be used
// only to release the correspondent call to Acquire, calling the
// returned Release function twice will lock your goroutine forever since
// this is a very stupid thing to do (better than releasing the semaphore twice
// I think, at least the bugged code is the one who will get screwed).
//
// Never calling release is also a terrible idea since this may cause
// starvation of resources if the semaphore is used to provide controlled
// access to some resource (usually an expensive one).
func (s S) Acquire(ctx context.Context) (Release, error) {
	select {
	case s <- struct{}{}:
		{
			r := &sync.Mutex{}
			released := false
			return func() {
				r.Lock()
				defer r.Unlock()

				if released {
					panic("released semaphore twice for the same Acquire")
				}
				released = true
				<-s
			}, nil
		}
	case <-ctx.Done():
		{
			return func() {}, fmt.Errorf("error[%s] waiting for semaphore", ctx.Err())
		}
	}
}
