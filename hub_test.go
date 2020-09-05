package hubsub

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Successful(t *testing.T) {
	t.Run("routing", func(t *testing.T) {
		h := NewHub()
		subA, chA := h.Subscribe(MapFromPairs("client", "a"))
		subB, chB := h.Subscribe(MapFromPairs("client", "b"))

		h.Publish("1", FilterFromPairs("client", "a"))
		h.Publish("2", FilterFromPairs("client", "b"))
		h.Publish("3", FilterFromPairs("client", "d"))
		h.Publish("4", FilterFromPairs("client", "b"))
		h.Publish("5", FilterFromPairs("client", "a"))
		h.Unsubscribe(subA, subB)
		gotA := flushStrChan(chA)
		gotB := flushStrChan(chB)

		assert.EqualValues(t, "1\n5\n", gotA)
		assert.EqualValues(t, "2\n4\n", gotB)
	})
	t.Run("subUnsub", func(t *testing.T) {
		h := NewHub()
		h.Publish("1", FilterFromPairs("client", "a"))
		h.Publish("2", FilterFromPairs("client", "b"))

		subA, chA := h.Subscribe(MapFromPairs("client", "a"))
		h.Publish("3", FilterFromPairs("client", "a"))
		h.Publish("4", FilterFromPairs("client", "b"))

		h.Unsubscribe(subA)
		h.Publish("5", FilterFromPairs("client", "a"))
		h.Publish("6", FilterFromPairs("client", "b"))
		gotA := flushStrChan(chA)

		assert.EqualValues(t, "3\n", gotA)
	})
	t.Run("slowClient", func(t *testing.T) {
		stored := DefaultBufferCap
		defer func() {
			DefaultBufferCap = stored
		}()

		DefaultBufferCap = 2

		h := NewHub()

		subA, chA := h.Subscribe(MapFromPairs("client", "a"))
		h.Publish("1", FilterFromPairs("client", "a"))
		h.Publish("2", FilterFromPairs("client", "b"))
		h.Publish("3", FilterFromPairs("client", "a"))
		h.Publish("4", FilterFromPairs("client", "b"))
		h.Publish("5", FilterFromPairs("client", "a"))
		h.Publish("6", FilterFromPairs("client", "b"))
		h.Publish("7", FilterFromPairs("client", "a"))
		h.Publish("8", FilterFromPairs("client", "b"))

		h.Unsubscribe(subA)
		gotA := flushStrChan(chA)

		assert.EqualValues(t, "1\n3\n", gotA)
	})
	t.Run("clientUnsub", func(t *testing.T) {
		h := NewHub()

		subA, chA := h.Subscribe(MapFromPairs("client", "a"))
		h.Publish("1", FilterFromPairs("client", "a"))
		h.Publish("2", FilterFromPairs("client", "a"))
		h.Publish("3", FilterFromPairs("client", "a"))

		h.Unsubscribe(subA)
		h.Publish("4", FilterFromPairs("client", "a"))
		h.Publish("5", FilterFromPairs("client", "a"))
		gotA := flushStrChan(chA)

		assert.EqualValues(t, "1\n2\n3\n", gotA)
	})
}

func Benchmark_SubscriberAndOneMatcher(b *testing.B) {
	rower := func(ch <-chan Message) {
		for {
			_, ok := <-ch
			if !ok {
				return
			}
			if cap(ch) > DefaultBufferCap {
				panic("cap(ch) > DefaultBufferCap")
			}
		}
	}
	filterClient := func(clientID string) func(in map[string]string) bool {
		return func(in map[string]string) bool {
			return in["client"] == clientID
		}
	}
	b.Run("1", func(b *testing.B) {
		h := NewHub()
		for i := 0; i < 1; i++ {
			meta := MapFromPairs("client", fmt.Sprint(i+1))
			_, ch := h.Subscribe(meta)
			go rower(ch)
		}

		filter := filterClient("1")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Publish("1", filter)
		}
	})

	b.Run("10start", func(b *testing.B) {
		h := NewHub()
		for i := 0; i < 10; i++ {
			meta := MapFromPairs("client", fmt.Sprint(i+1))
			_, ch := h.Subscribe(meta)
			go rower(ch)
		}

		filter := filterClient("1")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Publish("1", filter)
		}
	})

	b.Run("10middle", func(b *testing.B) {
		h := NewHub()
		for i := 0; i < 10; i++ {
			meta := MapFromPairs("client", fmt.Sprint(i+1))
			_, ch := h.Subscribe(meta)
			go rower(ch)
		}

		filter := filterClient("5")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Publish("1", filter)
		}
	})

	b.Run("100start", func(b *testing.B) {
		h := NewHub()
		for i := 0; i < 100; i++ {
			meta := MapFromPairs("client", fmt.Sprint(i+1))
			_, ch := h.Subscribe(meta)
			go rower(ch)
		}

		filter := filterClient("1")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Publish("1", filter)
		}
	})

	b.Run("100middle", func(b *testing.B) {
		h := NewHub()
		for i := 0; i < 100; i++ {
			meta := MapFromPairs("client", fmt.Sprint(i+1))
			_, ch := h.Subscribe(meta)
			go rower(ch)
		}

		filter := filterClient("50")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Publish("1", filter)
		}
	})

	b.Run("1000start", func(b *testing.B) {
		h := NewHub()
		for i := 0; i < 1000; i++ {
			meta := MapFromPairs("client", fmt.Sprint(i+1))
			_, ch := h.Subscribe(meta)
			go rower(ch)
		}

		filter := filterClient("1")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Publish("1", filter)
		}
	})

	b.Run("1000middle", func(b *testing.B) {
		h := NewHub()
		for i := 0; i < 1000; i++ {
			meta := MapFromPairs("client", fmt.Sprint(i+1))
			_, ch := h.Subscribe(meta)
			go rower(ch)
		}

		filter := filterClient("500")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Publish("1", filter)
		}
	})
}

func Example() {
	h := NewHub()

	subID, ch := h.Subscribe(MapFromPairs("user", "Bob"))

	defer func() {
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

func flushStrChan(in <-chan Message) string {
	buf := new(bytes.Buffer)
	for {
		msg, ok := <-in
		if !ok {
			break
		}
		fmt.Fprintln(buf, msg)
	}
	return buf.String()
}
