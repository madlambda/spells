package semaphore_test

import (
	"fmt"
	"time"
	"testing"
	"context"
	
	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/semaphore"
)

func TestSemaphore(t *testing.T) {
	sizes := []uint{ 1, 2, 3, 1000 }
	
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
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic on semaphore with size 0")
        }
    }()
    
    semaphore.New(0)
}

func TestSemaphoreCantReleaseSameAcquireTwice(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic releasing semaphore twice")
        }
    }()

	s := semaphore.New(1)
	release, err := s.Acquire(context.Background())
	
	assert.NoError(t, err)
	release()
	release()
}

func testSemaphore(t *testing.T, s semaphore.S, size uint) {
	
	releases := []semaphore.Release{}
	
	for i := uint(0); i < size; i++ {
		release, err := s.Acquire(context.Background())
		assert.NoError(t, err)
		releases = append(releases, release)
	}
	
	timeoutWorked := make(chan struct{})
	go func() {
		timeout := 100 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		
		_, err := s.Acquire(ctx)
		assert.Error(t, err)
		if ctx.Err() != context.DeadlineExceeded {
			t.Fatalf("expected deadline exceeded on acquire context, got[%s]", ctx.Err())
		}
		timeoutWorked <- struct{}{}
	}()
	
	<-timeoutWorked
	
	releases[0]()
	releases = releases[1:]
	
	r, err := s.Acquire(context.Background())
	assert.NoError(t, err)
	
	releases = append(releases, r)
	
	for _, release := range releases {
		release()
	}
}