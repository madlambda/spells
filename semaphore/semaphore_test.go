package semaphore_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/semaphore"
)

func TestSemaphore(t *testing.T) {
	sizes := []uint{1, 2, 3, 1000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			s := semaphore.New(size)
			// WHY: each test function releases all acquires
			// we use this to test if the semaphore will reset to
			// its initial state properly after all releases.
			testSemaphore(t, s, size)
			testSemaphore(t, s, size)
		})
	}

}

func TestSemaphoreSizeCantBeZero(t *testing.T) {
	defer assertPanic(t, "expected panic creating semaphore with size 0")
	semaphore.New(0)
}

func TestSemaphoreCantReleaseSameAcquireTwice(t *testing.T) {
	defer assertPanic(t, "Expected panic releasing semaphore twice")

	s := semaphore.New(1)
	release, err := s.Acquire(context.Background())

	assert.NoError(t, err)
	release()
	release()
}

func assertPanic(t *testing.T, errmsg string) {
	r := recover()
	if r == nil {
		t.Error(errmsg)
	}
	if _, ok := r.(runtime.Error); ok {
		t.Errorf("Unexpected runtime error[%s]", r)
	}
}

func testSemaphore(t *testing.T, s semaphore.S, size uint) {

	releases := []semaphore.Release{}

	for i := uint(0); i < size; i++ {
		release, err := s.Acquire(context.Background())
		assert.NoError(t, err)
		releases = append(releases, release)
	}

	timeoutWorked := make(chan error)
	go func() {
		timeout := 100 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := s.Acquire(ctx)
		if err == nil {
			timeoutWorked <- errors.New("expected timeout error, got none")
			return
		}
		if ctx.Err() != context.DeadlineExceeded {
			timeoutWorked <- fmt.Errorf("expected deadline exceeded on acquire context, got[%s]", ctx.Err())
			return
		}
		timeoutWorked <- nil
	}()

	err := <-timeoutWorked
	assert.NoError(t, err)

	releases[0]()
	releases = releases[1:]

	r, err := s.Acquire(context.Background())
	assert.NoError(t, err)

	releases = append(releases, r)

	for _, release := range releases {
		release()
	}
}
