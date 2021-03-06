package hubsub

import (
	"errors"
	"sync"
)

// DefaultBufferCap subscriber buffer size.
var DefaultBufferCap = 32

// NewHub returns a new instance of the hub.
func NewHub() *Hub {
	return &Hub{
		subs: make(map[uint]*subscribe),
	}
}

type Hub struct {
	rw        sync.RWMutex
	subs      map[uint]*subscribe // subID -> channel
	lastSubID uint
}

type subscribe struct {
	ch   chan Message
	meta map[string]string
}

// Publish publish a message to the hub for matching subscriptions.
//
// Filter function must be not-bloking.
func (h *Hub) Publish(in string, filterFn func(meta map[string]string) bool) error {
	if in == "" || filterFn == nil {
		return errors.New("interrupted publication - invalid input args")
	}

	return h.publish(in, filterFn)
}

func (h *Hub) publish(in string, matcherFn func(meta map[string]string) bool) error {
	if in == "" {
		return errors.New("interrupted publication - invalid input args")
	}

	h.rw.RLock()

	// Sending message if matched
	var toUnsubscribe []uint // TODO: get slice from pool?
	for subID, sub := range h.subs {
		if sub == nil {
			toUnsubscribe = append(toUnsubscribe, subID)
			continue
		}
		if !matcherFn(sub.meta) {
			continue
		}
		select {
		case sub.ch <- in:
		default:
			// Sub can't keep up. Will be closed.
			toUnsubscribe = append(toUnsubscribe, subID)
		}
	}

	h.rw.RUnlock()

	h.Unsubscribe(toUnsubscribe...)

	return nil
}

// Subscribe returns new subscription.
//
// Channel will be cloed if the client is slow (buffer will be filled).
func (h *Hub) Subscribe(meta map[string]string) (uint, <-chan Message) {
	ch := make(chan Message, DefaultBufferCap)

	h.rw.Lock()

	h.lastSubID++

	subID := h.lastSubID

	sub := &subscribe{
		ch:   ch,
		meta: meta,
	}
	h.subs[subID] = sub

	h.rw.Unlock()

	return subID, ch
}

// Unsubscribe unsubscribe for specified sub IDs.
func (h *Hub) Unsubscribe(subIDs ...uint) {
	h.rw.Lock()

	for _, subID := range subIDs {
		sub, ok := h.subs[subID]
		if ok {
			delete(h.subs, subID)
			close(sub.ch)
			continue
		}
	}

	h.rw.Unlock()
}

// Message this is the message container.
type Message interface{}
