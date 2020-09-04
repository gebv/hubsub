# hubsub

Low-level tool for managing subscriptions.

Popular of use:
* websocket, grpc streaming (only server publishing) // TODO: example code

# Quick example

```go
package main

func main() {
    h := NewHub()

	subID, ch := h.Subscribe(MapFromPairs("user", "Bob"))

	go func() {
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
}
```

[On the playground](https://play.golang.org/)

## TODO

* external fast index for matching
