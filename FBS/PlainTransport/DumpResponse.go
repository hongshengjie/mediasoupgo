// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package PlainTransport

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__SrtpParameters "mediasoupgo/FBS/SrtpParameters"
	FBS__Transport "mediasoupgo/FBS/Transport"
)

type DumpResponseT struct {
	Base *FBS__Transport.DumpT `json:"base"`
	RtcpMux bool `json:"rtcp_mux"`
	Comedia bool `json:"comedia"`
	Tuple *FBS__Transport.TupleT `json:"tuple"`
	RtcpTuple *FBS__Transport.TupleT `json:"rtcp_tuple"`
	SrtpParameters *FBS__SrtpParameters.SrtpParametersT `json:"srtp_parameters"`
}

func (t *DumpResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	baseOffset := t.Base.Pack(builder)
	tupleOffset := t.Tuple.Pack(builder)
	rtcpTupleOffset := t.RtcpTuple.Pack(builder)
	srtpParametersOffset := t.SrtpParameters.Pack(builder)
	DumpResponseStart(builder)
	DumpResponseAddBase(builder, baseOffset)
	DumpResponseAddRtcpMux(builder, t.RtcpMux)
	DumpResponseAddComedia(builder, t.Comedia)
	DumpResponseAddTuple(builder, tupleOffset)
	DumpResponseAddRtcpTuple(builder, rtcpTupleOffset)
	DumpResponseAddSrtpParameters(builder, srtpParametersOffset)
	return DumpResponseEnd(builder)
}

func (rcv *DumpResponse) UnPackTo(t *DumpResponseT) {
	t.Base = rcv.Base(nil).UnPack()
	t.RtcpMux = rcv.RtcpMux()
	t.Comedia = rcv.Comedia()
	t.Tuple = rcv.Tuple(nil).UnPack()
	t.RtcpTuple = rcv.RtcpTuple(nil).UnPack()
	t.SrtpParameters = rcv.SrtpParameters(nil).UnPack()
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

func (rcv *DumpResponse) RtcpMux() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *DumpResponse) MutateRtcpMux(n bool) bool {
	return rcv._tab.MutateBoolSlot(6, n)
}

func (rcv *DumpResponse) Comedia() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *DumpResponse) MutateComedia(n bool) bool {
	return rcv._tab.MutateBoolSlot(8, n)
}

func (rcv *DumpResponse) Tuple(obj *FBS__Transport.Tuple) *FBS__Transport.Tuple {
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

func (rcv *DumpResponse) RtcpTuple(obj *FBS__Transport.Tuple) *FBS__Transport.Tuple {
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

func (rcv *DumpResponse) SrtpParameters(obj *FBS__SrtpParameters.SrtpParameters) *FBS__SrtpParameters.SrtpParameters {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FBS__SrtpParameters.SrtpParameters)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func DumpResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func DumpResponseAddBase(builder *flatbuffers.Builder, base flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(base), 0)
}
func DumpResponseAddRtcpMux(builder *flatbuffers.Builder, rtcpMux bool) {
	builder.PrependBoolSlot(1, rtcpMux, false)
}
func DumpResponseAddComedia(builder *flatbuffers.Builder, comedia bool) {
	builder.PrependBoolSlot(2, comedia, false)
}
func DumpResponseAddTuple(builder *flatbuffers.Builder, tuple flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(tuple), 0)
}
func DumpResponseAddRtcpTuple(builder *flatbuffers.Builder, rtcpTuple flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(rtcpTuple), 0)
}
func DumpResponseAddSrtpParameters(builder *flatbuffers.Builder, srtpParameters flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(srtpParameters), 0)
}
func DumpResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
