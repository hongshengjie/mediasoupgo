// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package PipeTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "mediasoupgo/FBS/Transport"
)

type ConnectResponseT struct {
	Tuple *FBS__Transport.TupleT `json:"tuple"`
}

func (t *ConnectResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	tupleOffset := t.Tuple.Pack(builder)
	ConnectResponseStart(builder)
	ConnectResponseAddTuple(builder, tupleOffset)
	return ConnectResponseEnd(builder)
}

func (rcv *ConnectResponse) UnPackTo(t *ConnectResponseT) {
	t.Tuple = rcv.Tuple(nil).UnPack()
}

func (rcv *ConnectResponse) UnPack() *ConnectResponseT {
	if rcv == nil {
		return nil
	}
	t := &ConnectResponseT{}
	rcv.UnPackTo(t)
	return t
}

type ConnectResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsConnectResponse(buf []byte, offset flatbuffers.UOffsetT) *ConnectResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConnectResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishConnectResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConnectResponse(buf []byte, offset flatbuffers.UOffsetT) *ConnectResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConnectResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConnectResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConnectResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConnectResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConnectResponse) Tuple(obj *FBS__Transport.Tuple) *FBS__Transport.Tuple {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__Transport.Tuple)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func ConnectResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ConnectResponseAddTuple(builder *flatbuffers.Builder, tuple flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(tuple), 0)
}
func ConnectResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
