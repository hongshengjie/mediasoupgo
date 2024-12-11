package main

import (
	"bufio"
	"fmt"
	"mediasoupgo/FBS/Message"
	"mediasoupgo/FBS/Request"
	"os"
	"os/exec"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

func main() {
	cmd := exec.Command("mediasoup-worker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil
	producerReader, producerWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	consumerReader, consumerWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	cmd.ExtraFiles = []*os.File{producerReader, consumerWriter}

	cmd.Env = []string{"MEDIASOUP_VERSION=" + "3.15.2"}
	go func() {
		time.Sleep(time.Second * 2)
		b := flatbuffers.NewBuilder(0)
		m := &Message.MessageT{
			Data: &Message.BodyT{Type: Message.BodyRequest, Value: &Request.RequestT{
				Id:        1,
				Method:    Request.MethodWORKER_DUMP,
				HandlerId: "fsdfsdfsdf",
				Body:      &Request.BodyT{Type: Request.BodyNONE},
			}}}
		b.FinishSizePrefixed(m.Pack(b))
		fmt.Println("write len of data:", len(b.FinishedBytes()))
		producerWriter.Write(b.FinishedBytes())
	}()
	go func() {
		b := bufio.NewReader(consumerReader)
		for {
			bb := make([]byte, 10240)
			n, err := b.Read(bb)
			fmt.Println("read", string(bb), "len", n, err)
		}

	}()
	cmd.Start()
	cmd.Wait()
	time.Sleep(time.Hour)
}