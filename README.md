# timingwheel

golang fast timingwheel

## Usage

```go
package main

import (
	"time"

	"github.com/rfyiamcool/timingwheel"
)

func main() {
	tw := timingwheel.NewTimingWheel(1*time.Second, 60)
	tw.Start()

	tw.Sleep(1 * time.Second)

	select {
	case <-tw.After(1 * time.Second):
	}
}
```