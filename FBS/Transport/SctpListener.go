// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Common "mediasoupgo/FBS/Common"
)

type SctpListenerT struct {
	StreamIdTable []*FBS__Common.Uint16StringT `json:"stream_id_table"`
}

func (t *SctpListenerT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	streamIdTableOffset := flatbuffers.UOffsetT(0)
	if t.StreamIdTable != nil {
		streamIdTableLength := len(t.StreamIdTable)
		streamIdTableOffsets := make([]flatbuffers.UOffsetT, streamIdTableLength)
		for j := 0; j < streamIdTableLength; j++ {
			streamIdTableOffsets[j] = t.StreamIdTable[j].Pack(builder)
		}
		SctpListenerStartStreamIdTableVector(builder, streamIdTableLength)
		for j := streamIdTableLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(streamIdTableOffsets[j])
		}
		streamIdTableOffset = builder.EndVector(streamIdTableLength)
	}
	SctpListenerStart(builder)
	SctpListenerAddStreamIdTable(builder, streamIdTableOffset)
	return SctpListenerEnd(builder)
}

func (rcv *SctpListener) UnPackTo(t *SctpListenerT) {
	streamIdTableLength := rcv.StreamIdTableLength()
	t.StreamIdTable = make([]*FBS__Common.Uint16StringT, streamIdTableLength)
	for j := 0; j < streamIdTableLength; j++ {
		x := FBS__Common.Uint16String{}
		rcv.StreamIdTable(&x, j)
		t.StreamIdTable[j] = x.UnPack()
	}
}

func (rcv *SctpListener) UnPack() *SctpListenerT {
	if rcv == nil {
		return nil
	}
	t := &SctpListenerT{}
	rcv.UnPackTo(t)
	return t
}

type SctpListener struct {
	_tab flatbuffers.Table
}

func GetRootAsSctpListener(buf []byte, offset flatbuffers.UOffsetT) *SctpListener {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SctpListener{}
	x.Init(buf, n+offset)
	return x
}

func FinishSctpListenerBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSctpListener(buf []byte, offset flatbuffers.UOffsetT) *SctpListener {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SctpListener{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSctpListenerBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SctpListener) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SctpListener) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SctpListener) StreamIdTable(obj *FBS__Common.Uint16String, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *SctpListener) StreamIdTableLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func SctpListenerStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SctpListenerAddStreamIdTable(builder *flatbuffers.Builder, streamIdTable flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(streamIdTable), 0)
}
func SctpListenerStartStreamIdTableVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func SctpListenerEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
