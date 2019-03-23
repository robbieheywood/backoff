package main

import (
	"fmt"
	"testing"
	"time"
)

// In these tests, we generally use short times to avoid the testcases taking too long.
func TestRunningTicker(t *testing.T) {
	testcases := []struct {
		min    time.Duration
		max    time.Duration
		factor float32
	}{
		{min: 1 * time.Millisecond, max: 5 * time.Millisecond, factor: 1.1},
		{min: 2 * time.Millisecond, max: 200 * time.Millisecond, factor: 2},
		{min: 1 * time.Millisecond, max: 100 * time.Millisecond, factor: 5},
		{min: 2 * time.Millisecond, max: 2 * time.Millisecond, factor: 10},
		{min: 1 * time.Microsecond, max: 0, factor: 2},
	}

	for _, test := range testcases {
		t.Run(fmt.Sprintf("Test ticker with min '%v', max '%v' & factor '%v'", test.min, test.max, test.factor),
			func(t *testing.T) {
				ticker, err := NewTicker(test.min, test.max, test.factor)
				if err != nil {
					t.Fatalf("failed to create ticker with min '%v', max '%v' & factor '%v': %v", test.min, test.max, test.factor, err)
				}

				waitTime := test.min
				for i := 0; i < 10; i++ {
					earlyChan := time.After(waitTime - 5*time.Millisecond)
					lateChan := time.After(waitTime + 5*time.Millisecond)
					// Test that the ticker fires within 1ms of the expected time
					start := time.Now()
					select {
					case <-ticker.C:
						t.Fatalf("timer fired too early - fired after %v, expected after %v", time.Since(start), waitTime)
					case <-earlyChan:
					}
					select {
					case <-ticker.C:
					case <-lateChan:
						t.Fatalf("timer fired too late - fired after %v, expected after %v", time.Since(start), waitTime)
					}

					waitTime = time.Duration(float32(waitTime) * test.factor)
					if waitTime > test.max {
						waitTime = test.max
					}
				}

				ticker.Stop()
				_, ok := <-ticker.C
				if ok {
					t.Fatalf("")
				}
			})

	}
}

func TestTickerErrors(t *testing.T) {
	testcases := []struct {
		min    time.Duration
		max    time.Duration
		factor float32
	}{
		{min: 1 * time.Microsecond, max: 2 * time.Microsecond, factor: 0},
		{min: 1 * time.Microsecond, max: 2 * time.Microsecond, factor: -1},
		{min: 1 * time.Microsecond, max: 2 * time.Microsecond, factor: 0.9},
		{min: 0 * time.Microsecond, max: 2 * time.Microsecond, factor: 2},
		{min: -1 * time.Microsecond, max: 2 * time.Microsecond, factor: 2},
		{min: 1 * time.Microsecond, max: -1 * time.Microsecond, factor: 2},
		{min: 4 * time.Microsecond, max: 3 * time.Microsecond, factor: 2},
		{min: -1 * time.Microsecond, max: -1 * time.Microsecond, factor: -1},
	}

	for _, test := range testcases {
		t.Run(fmt.Sprintf("Test creating ticker with min '%v', max '%v' & factor '%v - expecting error", test.min, test.max, test.factor),
			func(t *testing.T) {
				_, err := NewTicker(test.min, test.max, test.factor)
				if err == nil {
					t.Fatalf("expected error when creating ticker with min '%v', max '%v' & factor '%v", test.min, test.max, test.factor)
				}
			})
	}
}
