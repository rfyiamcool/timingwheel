package timingwheel

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
		tw      = New(1*time.Second, 60)
		startAt time.Time
		notify  <-chan struct{}
	)

	tw.Start()

	startAt = time.Now()
	notify = tw.After(0 * time.Second)
	<-notify
	assert.LessOrEqual(t, time.Now().Unix()-startAt.Unix(), int64(1))

	startAt = time.Now()
	notify = tw.After(-1 * time.Second)
	<-notify
	assert.LessOrEqual(t, time.Now().Unix()-startAt.Unix(), int64(1))

	timer.Stop()
}

func TestSleep1s(t *testing.T) {
	var (
		tw = New(1*time.Second, 60)
	)

	tw.Start()
	defer tw.Stop()

	afterAt := time.Now()
	<-tw.After(1 * time.Second)
	<-tw.After(1 * time.Second)
	<-tw.After(1 * time.Second)
	assert.GreaterOrEqual(t, time.Now().Unix()-afterAt.Unix(), int64(3))

	sleepAt := time.Now()
	tw.Sleep(1 * time.Second)
	tw.Sleep(1 * time.Second)
	tw.Sleep(1 * time.Second)
	assert.GreaterOrEqual(t, time.Now().Unix()-sleepAt.Unix(), int64(3))

	sleepAt = time.Now()
	Sleep(1 * time.Second)
	assert.GreaterOrEqual(t, time.Now().Unix()-sleepAt.Unix(), int64(1))
}
