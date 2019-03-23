package backoff

import (
	"errors"
	"fmt"
	"time"
)

// Ticker holds a channel that delivers ticks that backoff exponentially.
// The ticker must be initialised with NewTicker/MustNewTicker -
// these are used to set the starting value, the backoff factor and a maximum value.
type Ticker struct {
	C <-chan struct{}

	// c is the internal version of C. It holds the same channel, but allows this package to write to it as well.
	c        chan struct{}
	min      time.Duration
	max      time.Duration
	factor   float32
	current  time.Duration
	shutChan chan struct{}
}

// MustNewTicker does the same as NewTicker except that it panics if provided with invalid inputs.
func MustNewTicker(min time.Duration, max time.Duration, factor float32) *Ticker {
	ticker, err := NewTicker(min, max, factor)
	if err != nil {
		panic("backoff timer init error: " + err.Error())
	}
	return ticker
}

// NewTicker returns a new Ticker containing a channel that will be pulsed after each backoff-time duration.
// The ticker starts immediately.
// The backoff time starts at 'min' and increases by 'factor' each time, until it reaches 'max'.
// To specify no max value, 'max' can be set to zero.
// 'min' must be > 0, 'max' must be >= 'min' and 'factor' must be >=- 1.
func NewTicker(min time.Duration, max time.Duration, factor float32) (*Ticker, error) {
	if min <= 0 {
		return nil, errors.New(fmt.Sprintf("non-positive min backoff interval of %v specified for ticker", min))
	}
	if max < min && max != 0 {
		return nil, errors.New(fmt.Sprintf("non-positive max backoff interval of %v specified for ticker", max))
	}
	if factor <= 1 {
		return nil, errors.New("zero backoff factor specified for ticker")
	}

	// The channel has a 1-element buffer.
	// If the client doesn't collect the tick, then subsequent ticks are dropped until the client starts reading.
	c := make(chan struct{}, 1)
	shutChan := make(chan struct{})
	ticker := &Ticker{
		C:        c,
		c:        c,
		min:      min,
		max:      max,
		current:  min,
		factor:   factor,
		shutChan: shutChan,
	}

	go func() {
		ticker.run()
	}()

	return ticker, nil
}

func (t *Ticker) run() {
	defer close(t.c)
	for {
		select {
		case <-t.shutChan:
			return
		case <-time.After(t.current):
		}
		t.tick()
		newTime := time.Duration(float32(t.current) * t.factor)
		if t.max > 0 && newTime > t.max {
			newTime = t.max
		}
		t.current = newTime
	}
}

func (t *Ticker) tick() {
	select {
	case t.c <- struct{}{}:
	default:
	}
}

// Stop is used to stop the ticker.
// Once the ticker is stopped, it cannot be restarted - a new ticker must be instantiated.
func (t *Ticker) Stop() {
	t.shutChan <- struct{}{}
}
