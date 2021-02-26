package timingwheel

import (
	"sync"
	"time"
)

var (
	nullChan = make(chan struct{})
)

func init() {
	close(nullChan)
}

type TimingWheel struct {
	sync.Mutex

	interval   time.Duration
	maxTimeout time.Duration

	ticker *time.Ticker
	quit   chan struct{}
	once   sync.Once

	// ring position
	pos  int
	ring []chan struct{}
}

func New(interval time.Duration, buckets int) *TimingWheel {
	tw := new(TimingWheel)

	tw.interval = interval
	tw.quit = make(chan struct{})
	tw.pos = 0
	tw.maxTimeout = time.Duration(interval * (time.Duration(buckets)))
	tw.ring = make([]chan struct{}, buckets)

	for i := range tw.ring {
		tw.ring[i] = make(chan struct{})
	}

	tw.ticker = time.NewTicker(interval)

	return tw
}

func (tw *TimingWheel) Start() {
	tw.once.Do(func() {
		go tw.run()
	})
}

func (tw *TimingWheel) Stop() {
	close(tw.quit)
}

func (tw *TimingWheel) Sleep(timeout time.Duration) {
	tw.After(timeout)
}

func (tw *TimingWheel) After(timeout time.Duration) chan struct{} {
	if timeout >= tw.maxTimeout {
		panic("timeout too much, over maxtimeout")
	}
	if timeout < 0 {
		return nullChan
	}

	index := int(timeout / tw.interval)
	if index > 0 {
		index--
	}

	tw.Lock()
	defer tw.Unlock()

	index = (tw.pos + index) % len(tw.ring)
	return tw.ring[index]
}

func (tw *TimingWheel) run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.onTicker()
		case <-tw.quit:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimingWheel) onTicker() {
	tw.Lock()
	oldChan := tw.ring[tw.pos]
	tw.ring[tw.pos] = make(chan struct{})
	tw.pos = (tw.pos + 1) % len(tw.ring)
	tw.Unlock()

	close(oldChan)
}
