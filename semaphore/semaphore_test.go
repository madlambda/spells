package semaphore_test

import (
	"testing"
	"github.com/madlambda/spells/semaphore"
)

func TestSemaphoreSizeCantBeZero(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic on sempahore with size 0")
        }
    }()
    
    semaphore.New(0)
}