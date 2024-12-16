package meidsoupgo

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"

	"mediasoupgo/FBS/Log"
	"mediasoupgo/FBS/Message"
	"mediasoupgo/FBS/Notification"
	"mediasoupgo/FBS/Request"
	"mediasoupgo/FBS/Response"
)

const (
	MESSAGE_MAX_LEN = 4194308
	PAYLOAD_MAX_LEN = 4194304
)

type sent struct {
	notify   chan struct{}
	id       uint32
	method   Request.Method
	response *Response.ResponseT
}

type Channel struct {
	producerSocket *os.File
	consumerSocket *os.File
	r              *bufio.Reader
	closed         atomic.Bool
	nextId         atomic.Uint32
	pid            uint32
	sentsMutex     sync.RWMutex
	sents          map[uint32]*sent
}

func NewChannel(producerWriter, consumerReader *os.File, pid uint32) *Channel {
	c := &Channel{
		pid:            pid,
		producerSocket: producerWriter,
		consumerSocket: consumerReader,
		r:              bufio.NewReader(consumerReader),

		sents: make(map[uint32]*sent),
	}
	go c.readLoop()
	return c
}

func (c *Channel) Close() {
	if c.closed.CompareAndSwap(false, true) {
		c.producerSocket.Close()
		c.consumerSocket.Close()
	}
}

func (c *Channel) Notify(event Notification.Event, body *Notification.BodyT, handleId string) error {
	m := &Message.MessageT{
		Data: &Message.BodyT{Type: Message.BodyNotification, Value: &Notification.NotificationT{
			HandlerId: handleId,
			Event:     event,
			Body:      body,
		}},
	}

	b := flatbuffers.NewBuilder(0)
	b.FinishSizePrefixed(m.Pack(b))
	_, err := c.producerSocket.Write(b.FinishedBytes())
	if err != nil {
		return err
	}
	return nil
}

func (c *Channel) Request(method Request.Method, body *Request.BodyT, handleId string) (*Response.ResponseT, error) {
	if c.nextId.Load() < 4294967295 {
	} else {
		c.nextId.Store(1)
	}

	var id uint32
	id = c.nextId.Add(1)
	m := &Message.MessageT{
		Data: &Message.BodyT{Type: Message.BodyRequest, Value: &Request.RequestT{
			Id:        id,
			HandlerId: handleId,
			Method:    method,
			Body:      body,
		}},
	}

	b := flatbuffers.NewBuilder(0)
	b.FinishSizePrefixed(m.Pack(b))

	notify := make(chan struct{})
	s := &sent{id: id, method: method, notify: notify}
	c.addsent(s)
	_, err := c.producerSocket.Write(b.FinishedBytes())
	if err != nil {
		return nil, err
	}
	select {
	case <-notify:
		c.removesent(id)
		return s.response, nil
	case <-time.After(time.Second):
		return nil, errors.New("time out")
	}
}

func (c *Channel) addsent(s *sent) {
	c.sentsMutex.Lock()
	c.sents[s.id] = s
	c.sentsMutex.Unlock()
}

func (c *Channel) getsent(id uint32) *sent {
	c.sentsMutex.RLock()
	x := c.sents[id]
	c.sentsMutex.RUnlock()
	return x
}

func (c *Channel) removesent(id uint32) {
	c.sentsMutex.Lock()
	delete(c.sents, id)
	c.sentsMutex.Unlock()
}

func (c *Channel) readLoop() error {
	defer func() {
		c.Close()
	}()
	for {
		if c.closed.Load() {
			return nil
		}
		l, err := c.r.Peek(4)
		if err != nil {
			return err
		}
		length := binary.LittleEndian.Uint32(l)
		if length > PAYLOAD_MAX_LEN {
			return errors.New("playload is too big")
		}
		data := make([]byte, length+4)
		_, err = io.ReadFull(c.r, data)
		if err != nil {
			return err
		}
		msg := Message.GetSizePrefixedRootAsMessage(data, 0)
		msgT := msg.UnPack()

		switch msgT.Data.Type {
		case Message.BodyLog:
			c.processLog(c.pid, msgT.Data.Value.(*Log.LogT))
		case Message.BodyNotification:
			c.processNotification(msgT.Data.Value.(*Notification.NotificationT))
		case Message.BodyResponse:
			c.processRespone(msgT.Data.Value.(*Response.ResponseT))
		}
	}
}

func (c *Channel) processLog(pid uint32, msgT *Log.LogT) {
}

func (c *Channel) processNotification(notification *Notification.NotificationT) {
}

func (c *Channel) processRespone(response *Response.ResponseT) {
	s := c.getsent(response.Id)
	if s != nil {
		s.response = response
		close(s.notify)
	}
}
