package messagewindow

import (
	"context"
	"errors"
	"io"
	"sync"

	"golang.org/x/exp/constraints"
)

// ErrClosedMessageWindow is an error that senders oberve if window have been closed.
var ErrClosedMessageWindow = errors.New("Message window have been closed")

// MessageWindow is an object that bridge incoming messages to a consumer, with following guarantees:
// - Message sender won't be blocked always
// - Consumer will observe sorted message sequence
type MessageWindow[S constraints.Integer, T any] interface {
	// Send message into this window.
	// This method won't block and concurrency safe.
	//
	// Sequence number must start with 0.
	//
	// If this window have been closed by the consumer, returns ErrClosedMessageWindow
	Send(sequence S, msg T) error

	// Consume message from this window.
	//
	// Because the aim of this window is to provide sequential message stream, only one goroutine should consume messages.
	// If multiple goroutines consume, the result is not guaranteed.
	//
	// If this window have been closed (= end of sequence), this method returns io.EOF.
	// Also, this method block until next message, context cancelation, or close.
	Consume(ctx context.Context) (S, T, error)

	// Close means the end of message sequence (if sender call) or means receiver's will to not receive message anymore.
	// At least one of sender or receiver must call Close to avoid resource leakage.
	// Concurrent / duplicated call of Close() is fine.
	Close()
}

type messageWindowImpl[S constraints.Integer, T any] struct {
	notify chan struct{}

	lock      sync.Mutex
	closed    error
	seqOffset S
	available []bool
	window    []T
}

func NewMessageWindow[S constraints.Integer, T any]() MessageWindow[S, T] {
	return &messageWindowImpl[S, T]{
		notify: make(chan struct{}, 1),
	}
}

func (w *messageWindowImpl[S, T]) Close() {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.closed != nil {
		return
	}

	w.closed = io.EOF
	close(w.notify)
}

func (w *messageWindowImpl[S, T]) Send(sequence S, msg T) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.closed != nil {
		return ErrClosedMessageWindow
	}

	if sequence < w.seqOffset { // The message have been already consumed by consumer
		return nil
	}
	i := int(sequence - w.seqOffset)            // Always >= 0, and expecting not so too huge number...
	if i < len(w.available) && w.available[i] { // Duplicated message have been already pushed
		// We always take ealiest one to keep consistency with the `sequence < w.offset` case
		return nil
	}

	for len(w.available) <= i {
		w.available = append(w.available, false)
		var empty T
		w.window = append(w.window, empty)
	}
	w.available[i] = true
	w.window[i] = msg
	w.pushed()
	return nil
}

func (w *messageWindowImpl[S, T]) pushed() {
	select {
	case w.notify <- struct{}{}:
		return
	default:
		return
	}
}

func (w *messageWindowImpl[S, T]) Consume(ctx context.Context) (S, T, error) {
	var seq S
	var msg T

	consumeIfAvailable := func() (bool, error) {
		w.lock.Lock()
		defer w.lock.Unlock()

		if len(w.available) == 0 {
			// Iff no message queued and window closed, retruns closed error
			return false, w.closed
		} else if !w.available[0] {
			// Even if the window have been closed, we must deliver messages that have been already pushed
			// Hence, here we should not return closed (EOF) error to let consumer continue.
			return false, nil
		} else {
			seq = w.seqOffset
			msg = w.window[0]

			w.seqOffset++
			w.available = w.available[1:]
			w.window = w.window[1:]
			return true, nil
		}
	}

	// If multiple messages have been sent, window holds them but w.notify queue holds <=1 notification.
	// Hence this consumer checks remaning messages in the window first.
	if available, err := consumeIfAvailable(); available {
		return seq, msg, err
	}

	for {
		select {
		case <-ctx.Done():
			return seq, msg, ctx.Err()
		case <-w.notify:
			// Even if the window have been closed, we must deliver messages that have been already sent.
			// Hence, here we check closed flag later.
			if available, err := consumeIfAvailable(); available || err != nil {
				return seq, msg, err
			}
			// "empty" notification can happen if...
			// - message pushed out-of-order
			// - the 1st consumeIfAvailable() call consumes message just after sender updated window.
			continue
		}
	}
}
