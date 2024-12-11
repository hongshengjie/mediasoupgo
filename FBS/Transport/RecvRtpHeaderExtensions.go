// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RecvRtpHeaderExtensionsT struct {
	Mid *byte `json:"mid"`
	Rid *byte `json:"rid"`
	Rrid *byte `json:"rrid"`
	AbsSendTime *byte `json:"abs_send_time"`
	TransportWideCc01 *byte `json:"transport_wide_cc01"`
}

func (t *RecvRtpHeaderExtensionsT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	RecvRtpHeaderExtensionsStart(builder)
	if t.Mid != nil {
		RecvRtpHeaderExtensionsAddMid(builder, *t.Mid)
	}
	if t.Rid != nil {
		RecvRtpHeaderExtensionsAddRid(builder, *t.Rid)
	}
	if t.Rrid != nil {
		RecvRtpHeaderExtensionsAddRrid(builder, *t.Rrid)
	}
	if t.AbsSendTime != nil {
		RecvRtpHeaderExtensionsAddAbsSendTime(builder, *t.AbsSendTime)
	}
	if t.TransportWideCc01 != nil {
		RecvRtpHeaderExtensionsAddTransportWideCc01(builder, *t.TransportWideCc01)
	}
	return RecvRtpHeaderExtensionsEnd(builder)
}

func (rcv *RecvRtpHeaderExtensions) UnPackTo(t *RecvRtpHeaderExtensionsT) {
	t.Mid = rcv.Mid()
	t.Rid = rcv.Rid()
	t.Rrid = rcv.Rrid()
	t.AbsSendTime = rcv.AbsSendTime()
	t.TransportWideCc01 = rcv.TransportWideCc01()
}

func (rcv *RecvRtpHeaderExtensions) UnPack() *RecvRtpHeaderExtensionsT {
	if rcv == nil {
		return nil
	}
	t := &RecvRtpHeaderExtensionsT{}
	rcv.UnPackTo(t)
	return t
}

type RecvRtpHeaderExtensions struct {
	_tab flatbuffers.Table
}

func GetRootAsRecvRtpHeaderExtensions(buf []byte, offset flatbuffers.UOffsetT) *RecvRtpHeaderExtensions {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RecvRtpHeaderExtensions{}
	x.Init(buf, n+offset)
	return x
}

func FinishRecvRtpHeaderExtensionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsRecvRtpHeaderExtensions(buf []byte, offset flatbuffers.UOffsetT) *RecvRtpHeaderExtensions {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RecvRtpHeaderExtensions{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedRecvRtpHeaderExtensionsBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *RecvRtpHeaderExtensions) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RecvRtpHeaderExtensions) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RecvRtpHeaderExtensions) Mid() *byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		v := rcv._tab.GetByte(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RecvRtpHeaderExtensions) MutateMid(n byte) bool {
	return rcv._tab.MutateByteSlot(4, n)
}

func (rcv *RecvRtpHeaderExtensions) Rid() *byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		v := rcv._tab.GetByte(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RecvRtpHeaderExtensions) MutateRid(n byte) bool {
	return rcv._tab.MutateByteSlot(6, n)
}

func (rcv *RecvRtpHeaderExtensions) Rrid() *byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		v := rcv._tab.GetByte(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RecvRtpHeaderExtensions) MutateRrid(n byte) bool {
	return rcv._tab.MutateByteSlot(8, n)
}

func (rcv *RecvRtpHeaderExtensions) AbsSendTime() *byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		v := rcv._tab.GetByte(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RecvRtpHeaderExtensions) MutateAbsSendTime(n byte) bool {
	return rcv._tab.MutateByteSlot(10, n)
}

func (rcv *RecvRtpHeaderExtensions) TransportWideCc01() *byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		v := rcv._tab.GetByte(o + rcv._tab.Pos)
		return &v
	}
	return nil
}

func (rcv *RecvRtpHeaderExtensions) MutateTransportWideCc01(n byte) bool {
	return rcv._tab.MutateByteSlot(12, n)
}

func RecvRtpHeaderExtensionsStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func RecvRtpHeaderExtensionsAddMid(builder *flatbuffers.Builder, mid byte) {
	builder.PrependByte(mid)
	builder.Slot(0)
}
func RecvRtpHeaderExtensionsAddRid(builder *flatbuffers.Builder, rid byte) {
	builder.PrependByte(rid)
	builder.Slot(1)
}
func RecvRtpHeaderExtensionsAddRrid(builder *flatbuffers.Builder, rrid byte) {
	builder.PrependByte(rrid)
	builder.Slot(2)
}
func RecvRtpHeaderExtensionsAddAbsSendTime(builder *flatbuffers.Builder, absSendTime byte) {
	builder.PrependByte(absSendTime)
	builder.Slot(3)
}
func RecvRtpHeaderExtensionsAddTransportWideCc01(builder *flatbuffers.Builder, transportWideCc01 byte) {
	builder.PrependByte(transportWideCc01)
	builder.Slot(4)
}
func RecvRtpHeaderExtensionsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}