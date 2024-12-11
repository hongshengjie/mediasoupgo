// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Integer32ArrayT struct {
	Value []int32 `json:"value"`
}

func (t *Integer32ArrayT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	valueOffset := flatbuffers.UOffsetT(0)
	if t.Value != nil {
		valueLength := len(t.Value)
		Integer32ArrayStartValueVector(builder, valueLength)
		for j := valueLength - 1; j >= 0; j-- {
			builder.PrependInt32(t.Value[j])
		}
		valueOffset = builder.EndVector(valueLength)
	}
	Integer32ArrayStart(builder)
	Integer32ArrayAddValue(builder, valueOffset)
	return Integer32ArrayEnd(builder)
}

func (rcv *Integer32Array) UnPackTo(t *Integer32ArrayT) {
	valueLength := rcv.ValueLength()
	t.Value = make([]int32, valueLength)
	for j := 0; j < valueLength; j++ {
		t.Value[j] = rcv.Value(j)
	}
}

func (rcv *Integer32Array) UnPack() *Integer32ArrayT {
	if rcv == nil {
		return nil
	}
	t := &Integer32ArrayT{}
	rcv.UnPackTo(t)
	return t
}

type Integer32Array struct {
	_tab flatbuffers.Table
}

func GetRootAsInteger32Array(buf []byte, offset flatbuffers.UOffsetT) *Integer32Array {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Integer32Array{}
	x.Init(buf, n+offset)
	return x
}

func FinishInteger32ArrayBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsInteger32Array(buf []byte, offset flatbuffers.UOffsetT) *Integer32Array {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Integer32Array{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedInteger32ArrayBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Integer32Array) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Integer32Array) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Integer32Array) Value(j int) int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *Integer32Array) ValueLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Integer32Array) MutateValue(j int, n int32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func Integer32ArrayStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func Integer32ArrayAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(value), 0)
}
func Integer32ArrayStartValueVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func Integer32ArrayEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
