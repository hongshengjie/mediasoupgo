// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package DirectTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "mediasoupgo/FBS/Transport"
)

type DumpResponseT struct {
	Base *FBS__Transport.DumpT `json:"base"`
}

func (t *DumpResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	baseOffset := t.Base.Pack(builder)
	DumpResponseStart(builder)
	DumpResponseAddBase(builder, baseOffset)
	return DumpResponseEnd(builder)
}

func (rcv *DumpResponse) UnPackTo(t *DumpResponseT) {
	t.Base = rcv.Base(nil).UnPack()
}

func (rcv *DumpResponse) UnPack() *DumpResponseT {
	if rcv == nil {
		return nil
	}
	t := &DumpResponseT{}
	rcv.UnPackTo(t)
	return t
}

type DumpResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsDumpResponse(buf []byte, offset flatbuffers.UOffsetT) *DumpResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DumpResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishDumpResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsDumpResponse(buf []byte, offset flatbuffers.UOffsetT) *DumpResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &DumpResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedDumpResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *DumpResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DumpResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DumpResponse) Base(obj *FBS__Transport.Dump) *FBS__Transport.Dump {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__Transport.Dump)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func DumpResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func DumpResponseAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func DumpResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}