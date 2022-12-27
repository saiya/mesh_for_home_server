package messagewindow_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
	"time"

	mw "github.com/saiya/mesh_for_home_server/peering/messagewindow"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
	"golang.org/x/sync/errgroup"
)

func TestMessageWindow(t *testing.T) {
	for _, c := range []struct {
		senders        int
		nMsgPerSender  int
		pSendFutureMsg float64 // Probability to send messages in reverse order
	}{
		{senders: 1, nMsgPerSender: 0, pSendFutureMsg: 0.0},
		{senders: 1, nMsgPerSender: 1, pSendFutureMsg: 0.0},
		{senders: 1, nMsgPerSender: 1000, pSendFutureMsg: 0.0},
		{senders: 3, nMsgPerSender: 0, pSendFutureMsg: 0.0},
		{senders: 3, nMsgPerSender: 1, pSendFutureMsg: 0.0},
		{senders: 3, nMsgPerSender: 10000, pSendFutureMsg: 0.1},
		{senders: 10, nMsgPerSender: 100000, pSendFutureMsg: 0.1},
	} {
		t.Run(
			fmt.Sprintf("senders:%d,nMsgPerSender:%d,pSendFutureMsg:%v", c.senders, c.nMsgPerSender, c.pSendFutureMsg),
			func(t *testing.T) {
				messages := make([][]msg, c.senders)
				senderRngs := make([]*rand.Rand, c.senders)
				for sender := 0; sender < c.senders; sender++ {
					messages[sender] = make([]msg, c.nMsgPerSender)
					senderRngs[sender] = rand.New(rand.NewSource(time.Now().UnixNano()))
					for nMsg := 0; nMsg < c.nMsgPerSender; nMsg++ {
						messages[sender][nMsg] = msg{sender, nMsg*c.senders + sender}
					}
				}

				w := mw.NewMessageWindow[int, msg]()

				// Consumer consumes messages
				consumer := consumeAll(w)

				// Senders send all messages
				senders := new(errgroup.Group)
				for sender := 0; sender < c.senders; sender++ {
					msgs := messages[sender]
					rng := senderRngs[sender]
					senders.Go(func() error {
						for cursor := 0; cursor < len(msgs); {
							var i int
							if cursor+10 < len(msgs) && rng.Float64() < c.pSendFutureMsg {
								i = cursor + rng.Intn(10)
							} else {
								i = cursor
								cursor++
							}

							err := w.Send(msgs[i].seq, msgs[i])
							if err != nil {
								return err
							}
						}
						return nil
					})
				}
				if !assert.NoError(t, senders.Wait()) {
					return
				}
				w.Close()

				received := <-consumer
				received.assert(t, c.senders*c.nMsgPerSender)
			},
		)
	}
}

func TestMessageWindowEmpty(t *testing.T) {
	w := mw.NewMessageWindow[int, msg]()

	_, _, err := consumeWithTimeout(w, 1*time.Millisecond)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestMessageWindowSimple(t *testing.T) {
	w := mw.NewMessageWindow[int, interface{}]()

	for i := 0; i <= 3; i++ {
		assert.NoError(t, w.Send(i, nil))
		seq, _, err := consumeWithTimeout(w, 1*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, i, seq)
	}
}

func TestOutOfOrderMessage(t *testing.T) {
	w := mw.NewMessageWindow[int, interface{}]()

	assert.NoError(t, w.Send(0, nil)) // Send #0 (consumer also observes it)
	assert.NoError(t, w.Send(0, nil)) // Send #0 again
	seq, _, err := consumeWithTimeout(w, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, seq)

	assert.NoError(t, w.Send(2, nil)) // Send #2 (consumer awaits #1)
	_, _, err = consumeWithTimeout(w, 1*time.Millisecond)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	assert.NoError(t, w.Send(0, nil)) // Send #0 again
	assert.NoError(t, w.Send(2, nil)) // Send #2 again

	assert.NoError(t, w.Send(1, nil)) // Send #1 (consumer observes #1 and #2)
	for i := 1; i <= 2; i++ {
		seq, _, err = consumeWithTimeout(w, 0)
		assert.NoError(t, err)
		assert.Equal(t, i, seq)
	}
	_, _, err = consumeWithTimeout(w, 1*time.Millisecond)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestMultipleClose(t *testing.T) {
	w := mw.NewMessageWindow[int, interface{}]()

	w.Close()
	w.Close()
	w.Close()

	_, msg, err := consumeWithTimeout(w, 1*time.Millisecond)
	assert.Nil(t, msg)
	assert.ErrorIs(t, err, io.EOF)

	assert.ErrorIs(t, w.Send(0, ""), mw.ErrClosedMessageWindow)
}

func TestConcurrentClose(t *testing.T) {
	w := mw.NewMessageWindow[int, interface{}]()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.Close()
		}()
	}
	wg.Wait()

	_, msg, err := consumeWithTimeout(w, 1*time.Millisecond)
	assert.Nil(t, msg)
	assert.ErrorIs(t, err, io.EOF)

	assert.ErrorIs(t, w.Send(0, ""), mw.ErrClosedMessageWindow)
}

type msg struct {
	sender int
	seq    int
}

func (m *msg) String() string {
	return fmt.Sprintf("msg{ sender=%d, seq=%d }", m.sender, m.seq)
}

func consumeWithTimeout[S constraints.Integer, T any](w mw.MessageWindow[S, T], timeout time.Duration) (S, T, error) {
	ctx, ctxClose := context.WithTimeout(context.Background(), timeout)
	defer ctxClose()
	return w.Consume(ctx)
}

type consumerResult struct {
	msgs      []msg
	anomaries map[string]interface{}
}

func (result *consumerResult) assert(t *testing.T, nMsgs int) {
	assert.Empty(t, result.anomaries)

	if assert.Equal(t, nMsgs, len(result.msgs)) {
		for i := 0; i < nMsgs; i++ {
			msg := result.msgs[i]
			if !assert.Equal(t, i, msg.seq) {
				break
			}
		}
	}
}

func consumeAll(w mw.MessageWindow[int, msg]) chan consumerResult {
	result := make(chan consumerResult)
	go func() {
		msgs := make([]msg, 0)
		anomaries := make(map[string]interface{})
		defer func() { result <- consumerResult{msgs, anomaries} }()

		prevMsg := msg{sender: -1, seq: -1}
		for {
			seq, msg, err := consumeWithTimeout(w, 10*time.Second)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				} else {
					anomaries["consumer encountered error"] = err
					continue
				}
			}
			if seq != msg.seq {
				anomaries["consumer found sequence unmatch"] = fmt.Sprintf("{ seq=%d, msg=%s }", seq, msg.String())
			}
			if prevMsg.seq >= msg.seq {
				anomaries["consumer encountered duplicate or reversed msg"] = fmt.Sprintf("{ previous msg=%s, msg=%s }", prevMsg.String(), msg.String())
			}
			if seq != len(msgs) {
				anomaries["consumer encountered incorrect sequence ID"] = fmt.Sprintf("{ expected=%d, msg=%s }", len(msgs), msg.String())
			}
			msgs = append(msgs, msg)
			prevMsg = msg
		}
	}()
	return result
}
