// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Common "mediasoupgo/internal/FBS/Common"
)

type RtpListenerT struct {
	SsrcTable []*FBS__Common.Uint32StringT `json:"ssrc_table"`
	MidTable []*FBS__Common.StringStringT `json:"mid_table"`
	RidTable []*FBS__Common.StringStringT `json:"rid_table"`
}

func (t *RtpListenerT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	ssrcTableOffset := flatbuffers.UOffsetT(0)
	if t.SsrcTable != nil {
		ssrcTableLength := len(t.SsrcTable)
		ssrcTableOffsets := make([]flatbuffers.UOffsetT, ssrcTableLength)
		for j := 0; j < ssrcTableLength; j++ {
			ssrcTableOffsets[j] = t.SsrcTable[j].Pack(builder)
		}
		RtpListenerStartSsrcTableVector(builder, ssrcTableLength)
		for j := ssrcTableLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(ssrcTableOffsets[j])
		}
		ssrcTableOffset = builder.EndVector(ssrcTableLength)
	}
	midTableOffset := flatbuffers.UOffsetT(0)
	if t.MidTable != nil {
		midTableLength := len(t.MidTable)
		midTableOffsets := make([]flatbuffers.UOffsetT, midTableLength)
		for j := 0; j < midTableLength; j++ {
			midTableOffsets[j] = t.MidTable[j].Pack(builder)
		}
		RtpListenerStartMidTableVector(builder, midTableLength)
		for j := midTableLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(midTableOffsets[j])
		}
		midTableOffset = builder.EndVector(midTableLength)
	}
	ridTableOffset := flatbuffers.UOffsetT(0)
	if t.RidTable != nil {
		ridTableLength := len(t.RidTable)
		ridTableOffsets := make([]flatbuffers.UOffsetT, ridTableLength)
		for j := 0; j < ridTableLength; j++ {
			ridTableOffsets[j] = t.RidTable[j].Pack(builder)
		}
		RtpListenerStartRidTableVector(builder, ridTableLength)
		for j := ridTableLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(ridTableOffsets[j])
		}
		ridTableOffset = builder.EndVector(ridTableLength)
	}
	RtpListenerStart(builder)
	RtpListenerAddSsrcTable(builder, ssrcTableOffset)
	RtpListenerAddMidTable(builder, midTableOffset)
	RtpListenerAddRidTable(builder, ridTableOffset)
	return RtpListenerEnd(builder)
}

func (rcv *RtpListener) UnPackTo(t *RtpListenerT) {
	ssrcTableLength := rcv.SsrcTableLength()
	t.SsrcTable = make([]*FBS__Common.Uint32StringT, ssrcTableLength)
	for j := 0; j < ssrcTableLength; j++ {
		x := FBS__Common.Uint32String{}
		rcv.SsrcTable(&x, j)
		t.SsrcTable[j] = x.UnPack()
	}
	midTableLength := rcv.MidTableLength()
	t.MidTable = make([]*FBS__Common.StringStringT, midTableLength)
	for j := 0; j < midTableLength; j++ {
		x := FBS__Common.StringString{}
		rcv.MidTable(&x, j)
		t.MidTable[j] = x.UnPack()
	}
	ridTableLength := rcv.RidTableLength()
	t.RidTable = make([]*FBS__Common.StringStringT, ridTableLength)
	for j := 0; j < ridTableLength; j++ {
		x := FBS__Common.StringString{}
		rcv.RidTable(&x, j)
		t.RidTable[j] = x.UnPack()
	}
}

func (rcv *RtpListener) UnPack() *RtpListenerT {
	if rcv == nil {
		return nil
	}
	t := &RtpListenerT{}
	rcv.UnPackTo(t)
	return t
}

type RtpListener struct {
	_tab flatbuffers.Table
}

func GetRootAsRtpListener(buf []byte, offset flatbuffers.UOffsetT) *RtpListener {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RtpListener{}
	x.Init(buf, n+offset)
	return x
}

func FinishRtpListenerBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRtpListener(buf []byte, offset flatbuffers.UOffsetT) *RtpListener {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RtpListener{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRtpListenerBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *RtpListener) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RtpListener) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RtpListener) SsrcTable(obj *FBS__Common.Uint32String, j int) bool {
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

func (rcv *RtpListener) SsrcTableLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RtpListener) MidTable(obj *FBS__Common.StringString, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *RtpListener) MidTableLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RtpListener) RidTable(obj *FBS__Common.StringString, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *RtpListener) RidTableLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RtpListenerStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func RtpListenerAddSsrcTable(builder *flatbuffers.Builder, ssrcTable flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(ssrcTable), 0)
}
func RtpListenerStartSsrcTableVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RtpListenerAddMidTable(builder *flatbuffers.Builder, midTable flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(midTable), 0)
}
func RtpListenerStartMidTableVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RtpListenerAddRidTable(builder *flatbuffers.Builder, ridTable flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(ridTable), 0)
}
func RtpListenerStartRidTableVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RtpListenerEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
