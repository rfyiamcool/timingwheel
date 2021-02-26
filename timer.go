package timewheel

import (
	"math/rand"
	"sync"
	"time"
)

var (
	count = 5
	tws   = []*TimingWheel{}

	once sync.Once
)

func init() {
	for i := 0; i < count; i++ {
		tw := NewTimingWheel(1*time.Second, 600) // 10 minite
		tw.Start()
		tws = append(tws, tw)
	}
}

func SetDefaultTimeingWheels(obj []*TimingWheel) {
	for _, tw := range tws {
		tw.Stop()
	}

	tws = obj

	for _, tw := range tws {
		tw.Start()
	}
}

func Sleep(timeout time.Duration) {
	After(timeout)
}

func After(timeout time.Duration) chan struct{} {
	n := rand.Intn(count) // safe array bound
	return tws[n].After(timeout)
}

func Stop() {
	once.Do(func() {
		for _, tw := range tws {
			tw.Stop()
		}
	})
}
