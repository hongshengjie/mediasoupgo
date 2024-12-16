package meidsoupgo

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"mediasoupgo/FBS/Message"
	"mediasoupgo/FBS/Notification"
	"mediasoupgo/FBS/Response"
	"mediasoupgo/FBS/Worker"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
)

type sentdata struct {
	notify chan struct{}
	reqId  uint32
}
type Channel struct {
	producerWriter *os.File
	consumerReader *os.File
	r              *bufio.Reader
}

func NewChannel(pw, cr *os.File) *Channel {
	c := &Channel{
		producerWriter: pw,
		consumerReader: cr,
		r:              bufio.NewReader(cr),
	}
	return c
}
func (c *Channel) Send(m *Message.MessageT) error {
	b := flatbuffers.NewBuilder(0)
	b.FinishSizePrefixed(m.Pack(b))
	_, err := c.producerWriter.Write(b.FinishedBytes())
	return err
}
func (c *Channel) readLoop() error {

	for {
		l, err := c.r.Peek(4)
		if err != nil {
			fmt.Println(err)
		}

		length := binary.LittleEndian.Uint32(l)
		fmt.Println("data length", length)

		data := make([]byte, length+4)
		n, err := io.ReadFull(c.r, data)
		if err != nil {
			fmt.Println(err, n)
		}
		msg := Message.GetSizePrefixedRootAsMessage(data, 0)
		msgT := msg.UnPack()
		fmt.Println(msgT.Data, msg.DataType())
		switch msgT.Data.Type {
		case Message.BodyLog:
		case Message.BodyNotification:
		case Message.BodyResponse:
		}
	}
}

func processLog(msgT *Message.MessageT) {

}
func processNotification(msgT *Message.MessageT) {

	notify, ok := msgT.Data.Value.(*Notification.NotificationT)
	fmt.Println(notify, ok)
}

func processRespone(msgT *Message.MessageT) {

	resp, ok := msgT.Data.Value.(*Response.ResponseT)
	if ok {
		switch resp.Body.Type {
		case Response.BodyWorker_DumpResponse:
			d, ok := resp.Body.Value.(*Worker.DumpResponseT)
			fmt.Println(d, ok)
		case Response.BodyWorker_ResourceUsageResponse:
			d, ok := resp.Body.Value.(*Worker.ResourceUsageResponseT)
			fmt.Println(d, ok)
		}
	}
}
