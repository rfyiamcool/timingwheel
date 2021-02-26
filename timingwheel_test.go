package timewheel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAbnormalCase(t *testing.T) {
	timer := time.AfterFunc(5*time.Second, func() {
		panic("timeout")
	})
	var (
		tw      = NewTimingWheel(1*time.Second, 60)
		startAt time.Time
		notify  <-chan struct{}
	)

	tw.Start()

	startAt = time.Now()
	notify = tw.After(0 * time.Second)
	<-notify
	assert.LessOrEqual(t, time.Now().Unix()-startAt.Unix(), int64(2))

	startAt = time.Now()
	notify = tw.After(-1 * time.Second)
	<-notify
	assert.LessOrEqual(t, time.Now().Unix()-startAt.Unix(), int64(3))

	timer.Stop()
}
