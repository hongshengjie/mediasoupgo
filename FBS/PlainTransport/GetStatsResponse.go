// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package PlainTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "mediasoupgo/FBS/Transport"
)

type GetStatsResponseT struct {
	Base *FBS__Transport.StatsT `json:"base"`
	RtcpMux bool `json:"rtcp_mux"`
	Comedia bool `json:"comedia"`
	Tuple *FBS__Transport.TupleT `json:"tuple"`
	RtcpTuple *FBS__Transport.TupleT `json:"rtcp_tuple"`
}

func (t *GetStatsResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	baseOffset := t.Base.Pack(builder)
	tupleOffset := t.Tuple.Pack(builder)
	rtcpTupleOffset := t.RtcpTuple.Pack(builder)
	GetStatsResponseStart(builder)
	GetStatsResponseAddBase(builder, baseOffset)
	GetStatsResponseAddRtcpMux(builder, t.RtcpMux)
	GetStatsResponseAddComedia(builder, t.Comedia)
	GetStatsResponseAddTuple(builder, tupleOffset)
	GetStatsResponseAddRtcpTuple(builder, rtcpTupleOffset)
	return GetStatsResponseEnd(builder)
}

func (rcv *GetStatsResponse) UnPackTo(t *GetStatsResponseT) {
	t.Base = rcv.Base(nil).UnPack()
	t.RtcpMux = rcv.RtcpMux()
	t.Comedia = rcv.Comedia()
	t.Tuple = rcv.Tuple(nil).UnPack()
	t.RtcpTuple = rcv.RtcpTuple(nil).UnPack()
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

func (rcv *GetStatsResponse) RtcpMux() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *GetStatsResponse) MutateRtcpMux(n bool) bool {
	return rcv._tab.MutateBoolSlot(6, n)
}

func (rcv *GetStatsResponse) Comedia() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *GetStatsResponse) MutateComedia(n bool) bool {
	return rcv._tab.MutateBoolSlot(8, n)
}

func (rcv *GetStatsResponse) Tuple(obj *FBS__Transport.Tuple) *FBS__Transport.Tuple {
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

func (rcv *GetStatsResponse) RtcpTuple(obj *FBS__Transport.Tuple) *FBS__Transport.Tuple {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
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

func GetStatsResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func GetStatsResponseAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func GetStatsResponseAddRtcpMux(builder *flatbuffers.Builder, rtcpMux bool) {
	builder.PrependBoolSlot(1, rtcpMux, false)
}
func GetStatsResponseAddComedia(builder *flatbuffers.Builder, comedia bool) {
	builder.PrependBoolSlot(2, comedia, false)
}
func GetStatsResponseAddTuple(builder *flatbuffers.Builder, tuple flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(tuple), 0)
}
func GetStatsResponseAddRtcpTuple(builder *flatbuffers.Builder, rtcpTuple flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(rtcpTuple), 0)
}
func GetStatsResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}