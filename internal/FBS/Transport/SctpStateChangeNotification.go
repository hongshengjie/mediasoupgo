// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__SctpAssociation "mediasoupgo/internal/FBS/SctpAssociation"
)

type SctpStateChangeNotificationT struct {
	SctpState FBS__SctpAssociation.SctpState `json:"sctp_state"`
}

func (t *SctpStateChangeNotificationT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	SctpStateChangeNotificationStart(builder)
	SctpStateChangeNotificationAddSctpState(builder, t.SctpState)
	return SctpStateChangeNotificationEnd(builder)
}

func (rcv *SctpStateChangeNotification) UnPackTo(t *SctpStateChangeNotificationT) {
	t.SctpState = rcv.SctpState()
}

func (rcv *SctpStateChangeNotification) UnPack() *SctpStateChangeNotificationT {
	if rcv == nil {
		return nil
	}
	t := &SctpStateChangeNotificationT{}
	rcv.UnPackTo(t)
	return t
}

type SctpStateChangeNotification struct {
	_tab flatbuffers.Table
}

func GetRootAsSctpStateChangeNotification(buf []byte, offset flatbuffers.UOffsetT) *SctpStateChangeNotification {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SctpStateChangeNotification{}
	x.Init(buf, n+offset)
	return x
}

func FinishSctpStateChangeNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSctpStateChangeNotification(buf []byte, offset flatbuffers.UOffsetT) *SctpStateChangeNotification {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SctpStateChangeNotification{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSctpStateChangeNotificationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SctpStateChangeNotification) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SctpStateChangeNotification) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SctpStateChangeNotification) SctpState() FBS__SctpAssociation.SctpState {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return FBS__SctpAssociation.SctpState(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *SctpStateChangeNotification) MutateSctpState(n FBS__SctpAssociation.SctpState) bool {
	return rcv._tab.MutateByteSlot(4, byte(n))
}

func SctpStateChangeNotificationStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SctpStateChangeNotificationAddSctpState(builder *flatbuffers.Builder, sctpState FBS__SctpAssociation.SctpState) {
	builder.PrependByteSlot(0, byte(sctpState), 0)
}
func SctpStateChangeNotificationEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
