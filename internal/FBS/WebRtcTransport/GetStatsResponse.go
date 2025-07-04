// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WebRtcTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "mediasoupgo/internal/FBS/Transport"
)

type GetStatsResponseT struct {
	Base *FBS__Transport.StatsT `json:"base"`
	IceRole IceRole `json:"ice_role"`
	IceState IceState `json:"ice_state"`
	IceSelectedTuple *FBS__Transport.TupleT `json:"ice_selected_tuple"`
	DtlsState DtlsState `json:"dtls_state"`
}

func (t *GetStatsResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	baseOffset := t.Base.Pack(builder)
	iceSelectedTupleOffset := t.IceSelectedTuple.Pack(builder)
	GetStatsResponseStart(builder)
	GetStatsResponseAddBase(builder, baseOffset)
	GetStatsResponseAddIceRole(builder, t.IceRole)
	GetStatsResponseAddIceState(builder, t.IceState)
	GetStatsResponseAddIceSelectedTuple(builder, iceSelectedTupleOffset)
	GetStatsResponseAddDtlsState(builder, t.DtlsState)
	return GetStatsResponseEnd(builder)
}

func (rcv *GetStatsResponse) UnPackTo(t *GetStatsResponseT) {
	t.Base = rcv.Base(nil).UnPack()
	t.IceRole = rcv.IceRole()
	t.IceState = rcv.IceState()
	t.IceSelectedTuple = rcv.IceSelectedTuple(nil).UnPack()
	t.DtlsState = rcv.DtlsState()
}

func (rcv *GetStatsResponse) UnPack() *GetStatsResponseT {
	if rcv == nil {
		return nil
	}
	t := &GetStatsResponseT{}
	rcv.UnPackTo(t)
	return t
}

type GetStatsResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsGetStatsResponse(buf []byte, offset flatbuffers.UOffsetT) *GetStatsResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &GetStatsResponse{}
	x.Init(buf, n+offset)
	return x
}

func FinishGetStatsResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsGetStatsResponse(buf []byte, offset flatbuffers.UOffsetT) *GetStatsResponse {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &GetStatsResponse{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedGetStatsResponseBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *GetStatsResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *GetStatsResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *GetStatsResponse) Base(obj *FBS__Transport.Stats) *FBS__Transport.Stats {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__Transport.Stats)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *GetStatsResponse) IceRole() IceRole {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return IceRole(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *GetStatsResponse) MutateIceRole(n IceRole) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *GetStatsResponse) IceState() IceState {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return IceState(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *GetStatsResponse) MutateIceState(n IceState) bool {
	return rcv._tab.MutateByteSlot(8, byte(n))
}

func (rcv *GetStatsResponse) IceSelectedTuple(obj *FBS__Transport.Tuple) *FBS__Transport.Tuple {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
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

func (rcv *GetStatsResponse) DtlsState() DtlsState {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return DtlsState(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *GetStatsResponse) MutateDtlsState(n DtlsState) bool {
	return rcv._tab.MutateByteSlot(12, byte(n))
}

func GetStatsResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func GetStatsResponseAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func GetStatsResponseAddIceRole(builder *flatbuffers.Builder, iceRole IceRole) {
	builder.PrependByteSlot(1, byte(iceRole), 0)
}
func GetStatsResponseAddIceState(builder *flatbuffers.Builder, iceState IceState) {
	builder.PrependByteSlot(2, byte(iceState), 0)
}
func GetStatsResponseAddIceSelectedTuple(builder *flatbuffers.Builder, iceSelectedTuple flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(iceSelectedTuple), 0)
}
func GetStatsResponseAddDtlsState(builder *flatbuffers.Builder, dtlsState DtlsState) {
	builder.PrependByteSlot(4, byte(dtlsState), 0)
}
func GetStatsResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
