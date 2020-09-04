# hubsub

![CI](https://github.com/gebv/hubsub/workflows/CI/badge.svg?branch=master)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gebv/hubsub)](https://pkg.go.dev/github.com/gebv/hubsub)
[![Go Report Card](https://goreportcard.com/badge/github.com/gebv/hubsub)](https://goreportcard.com/report/github.com/gebv/hubsub)

Low-level tool for managing subscriptions.

Popular of use:
- websocket, grpc streaming (only server publishing)

# Quick example

```go
package main

import (
	"fmt"
	"github.com/gebv/hubsub"
	"sync"
)

func main() {
	h := hubsub.NewHub()

	subID, ch := h.Subscribe(map[string]string{"user": "Bob"})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, ok := <-ch
			if !ok {
				return
			}
			fmt.Println("message received:", msg)
		}
	}()

	var (
		filterBob = func(meta map[string]string) bool {
			return meta["user"] == "Bob"
		}
		filterAlice = func(meta map[string]string) bool {
			return meta["user"] == "Alice"
		}
	)

	h.Publish("msg1", filterBob)
	h.Publish("msg2", filterAlice)
	h.Publish("msg3", filterBob)

	h.Unsubscribe(subID)

	h.Publish("msg4", filterAlice)
	h.Publish("msg5", filterBob)

	// Output:
	// message received: msg1
	// message received: msg3

	wg.Wait()
}
```

[On the playground](https://play.golang.org/p/2TMRB_yasJ7)

## TODO

- [ ] add example code for popular uses
- [ ] external fast index for matching
- [ ] prometheus collector
