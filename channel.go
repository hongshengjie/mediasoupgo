package mediasoupgo

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"unsafe"

	flatbuffers "github.com/google/flatbuffers/go"

	"mediasoupgo/FBS/Log"
	fbsMessage "mediasoupgo/FBS/Message"
	"mediasoupgo/FBS/Notification"
	"mediasoupgo/FBS/Request"
	"mediasoupgo/FBS/Response"
	"mediasoupgo/events"
)

const intWidth int = int(unsafe.Sizeof(0))

var byteOrder binary.ByteOrder

func ByteOrder() binary.ByteOrder { return byteOrder }

func init() {
	i := int(0x1)
	if v := (*[intWidth]byte)(unsafe.Pointer(&i)); v[0] == 0 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}

const (
	MESSAGE_MAX_LEN = 4194308
	PAYLOAD_MAX_LEN = 4194304
)

type Channel struct {
	producerSocket *os.File
	consumerSocket *os.File
	r              *bufio.Reader
	spawnDone      chan struct{}
	closed         atomic.Bool
	nextId         atomic.Uint32
	pid            int
	sentsMutex     sync.RWMutex
	sents          map[uint32]*sent
	events.EventEmmiter[*Notification.NotificationT]
}

func NewChannel(producerWriter, consumerReader *os.File) (*Channel, chan struct{}) {
	spawnDone := make(chan struct{})
	c := &Channel{
		producerSocket: producerWriter,
		consumerSocket: consumerReader,
		r:              bufio.NewReader(consumerReader),
		spawnDone:      spawnDone,
		sents:          make(map[uint32]*sent),
		EventEmmiter:   events.New[*Notification.NotificationT](),
	}
	go c.readLoop()
	return c, spawnDone
}

func (c *Channel) SetPid(pid int) {
	c.pid = pid
}

func (c *Channel) Close() {
	if c.closed.CompareAndSwap(false, true) {
		c.producerSocket.Close()
		c.consumerSocket.Close()
	}
}

func (c *Channel) Notify(
	event Notification.Event,
	body *Notification.BodyT,
	handleId string,
) error {
	m := &fbsMessage.MessageT{
		Data: &fbsMessage.BodyT{
			Type: fbsMessage.BodyNotification,
			Value: &Notification.NotificationT{
				HandlerId: handleId,
				Event:     event,
				Body:      body,
			},
		},
	}

	b := flatbuffers.NewBuilder(0)
	b.FinishSizePrefixed(m.Pack(b))
	_, err := c.producerSocket.Write(b.FinishedBytes())
	if err != nil {
		return err
	}
	return nil
}

func (c *Channel) Request(
	method Request.Method,
	body *Request.BodyT,
	handleId string,
) (*Response.ResponseT, error) {
	if c.nextId.Load() < 4294967295 {
	} else {
		c.nextId.Store(1)
	}

	var id uint32
	id = c.nextId.Add(1)
	m := &fbsMessage.MessageT{
		Data: &fbsMessage.BodyT{
			Type: fbsMessage.BodyRequest,
			Value: &Request.RequestT{
				Id:        id,
				HandlerId: handleId,
				Method:    method,
				Body:      body,
			},
		},
	}
	d, _ := json.Marshal(m)
	slog.Debug("request", "data", string(d))
	b := flatbuffers.NewBuilder(0)
	b.FinishSizePrefixed(m.Pack(b))

	notify := make(chan struct{})
	s := &sent{id: id, method: method, notify: notify}
	c.addsent(s)
	data := b.FinishedBytes()
	_, err := c.producerSocket.Write(data)
	if err != nil {
		return nil, err
	}
	select {
	case <-notify:
		c.removesent(id)
		return s.response, nil

		// case <-time.After(time.Second * 10):
		// 	return nil, errors.New("time out")
	}
}

func (c *Channel) readLoop() error {
	defer func() {
		slog.Info("readLoop end")
		c.Close()
	}()
	for {
		if c.closed.Load() {
			return nil
		}
		l, err := c.r.Peek(4)
		if err != nil {
			slog.Error("peek error", "error", err)
			return err
		}
		length := byteOrder.Uint32(l)
		if length > PAYLOAD_MAX_LEN {
			return errors.New("playload is too big")
		}
		data := make([]byte, length+4)
		_, err = io.ReadFull(c.r, data)
		if err != nil {
			return err
		}
		msg := fbsMessage.GetSizePrefixedRootAsMessage(data, 0)
		msgT := msg.UnPack()
		switch msgT.Data.Type {
		case fbsMessage.BodyLog:
			c.processLog(c.pid, msgT.Data.Value.(*Log.LogT))
		case fbsMessage.BodyNotification:
			c.processNotification(msgT.Data.Value.(*Notification.NotificationT))
		case fbsMessage.BodyResponse:
			c.processRespone(msgT.Data.Value.(*Response.ResponseT))
		default:
			slog.Debug("unknow message")
		}
	}
}

func (c *Channel) processRespone(response *Response.ResponseT) {
	s := c.getsent(response.Id)
	if s != nil {
		s.response = response
		close(s.notify)
	}
}

func (c *Channel) processNotification(notification *Notification.NotificationT) {
	switch notification.Event {
	case Notification.EventWORKER_RUNNING:
		close(c.spawnDone)
	default:
		c.Emit(events.EventName(notification.HandlerId), notification)

	}
}

func (c *Channel) processLog(pid int, msgT *Log.LogT) {
	switch msgT.Data[0] {
	case 'D':
		slog.Debug(msgT.Data, slog.Int("pid", pid))
	case 'W':
		slog.Warn(msgT.Data, slog.Int("pid", pid))
	case 'E':
		slog.Error(msgT.Data, slog.Int("pid", pid))
	case 'X':
		slog.Info(msgT.Data, slog.Int("pid", pid))

	}
}

type sent struct {
	notify   chan struct{}
	id       uint32
	method   Request.Method
	response *Response.ResponseT
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
