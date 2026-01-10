package main

import (
	"testing"
	"time"

	"github.com/juju/mutex/v2"
)

// TestSingleInstanceMutex tests that we can acquire and release the mutex
func TestSingleInstanceMutex(t *testing.T) {
	spec := mutex.Spec{
		Name:    "click-guardian-test-single-instance", // Valid name for testing
		Clock:   realClock{},                           // Use real-time clock
		Delay:   100 * time.Millisecond,                // Short polling interval for testing
		Timeout: 500 * time.Millisecond,                // Short timeout for testing
	}

	// Acquire the mutex
	releaser, err := mutex.Acquire(spec)
	if err != nil {
		t.Fatalf("Failed to acquire mutex: %v", err)
	}

	// Try to acquire the mutex again - this should fail
	_, err = mutex.Acquire(spec)
	if err == nil {
		t.Fatal("Expected to fail when acquiring mutex again, but succeeded")
	}

	// Release the mutex
	releaser.Release()

	// Now we should be able to acquire it again
	releaser2, err := mutex.Acquire(spec)
	if err != nil {
		t.Fatalf("Failed to acquire mutex after releasing: %v", err)
	}
	defer releaser2.Release()
}
