// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Worker

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Transport "mediasoupgo/internal/FBS/Transport"
)

type CreateWebRtcServerRequestT struct {
	WebRtcServerId string `json:"web_rtc_server_id"`
	ListenInfos []*FBS__Transport.ListenInfoT `json:"listen_infos"`
}

func (t *CreateWebRtcServerRequestT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	webRtcServerIdOffset := flatbuffers.UOffsetT(0)
	if t.WebRtcServerId != "" {
		webRtcServerIdOffset = builder.CreateString(t.WebRtcServerId)
	}
	listenInfosOffset := flatbuffers.UOffsetT(0)
	if t.ListenInfos != nil {
		listenInfosLength := len(t.ListenInfos)
		listenInfosOffsets := make([]flatbuffers.UOffsetT, listenInfosLength)
		for j := 0; j < listenInfosLength; j++ {
			listenInfosOffsets[j] = t.ListenInfos[j].Pack(builder)
		}
		CreateWebRtcServerRequestStartListenInfosVector(builder, listenInfosLength)
		for j := listenInfosLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(listenInfosOffsets[j])
		}
		listenInfosOffset = builder.EndVector(listenInfosLength)
	}
	CreateWebRtcServerRequestStart(builder)
	CreateWebRtcServerRequestAddWebRtcServerId(builder, webRtcServerIdOffset)
	CreateWebRtcServerRequestAddListenInfos(builder, listenInfosOffset)
	return CreateWebRtcServerRequestEnd(builder)
}

func (rcv *CreateWebRtcServerRequest) UnPackTo(t *CreateWebRtcServerRequestT) {
	t.WebRtcServerId = string(rcv.WebRtcServerId())
	listenInfosLength := rcv.ListenInfosLength()
	t.ListenInfos = make([]*FBS__Transport.ListenInfoT, listenInfosLength)
	for j := 0; j < listenInfosLength; j++ {
		x := FBS__Transport.ListenInfo{}
		rcv.ListenInfos(&x, j)
		t.ListenInfos[j] = x.UnPack()
	}
}

func (rcv *CreateWebRtcServerRequest) UnPack() *CreateWebRtcServerRequestT {
	if rcv == nil {
		return nil
	}
	t := &CreateWebRtcServerRequestT{}
	rcv.UnPackTo(t)
	return t
}

type CreateWebRtcServerRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsCreateWebRtcServerRequest(buf []byte, offset flatbuffers.UOffsetT) *CreateWebRtcServerRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CreateWebRtcServerRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishCreateWebRtcServerRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsCreateWebRtcServerRequest(buf []byte, offset flatbuffers.UOffsetT) *CreateWebRtcServerRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &CreateWebRtcServerRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedCreateWebRtcServerRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *CreateWebRtcServerRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CreateWebRtcServerRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CreateWebRtcServerRequest) WebRtcServerId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *CreateWebRtcServerRequest) ListenInfos(obj *FBS__Transport.ListenInfo, j int) bool {
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

func (rcv *CreateWebRtcServerRequest) ListenInfosLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func CreateWebRtcServerRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func CreateWebRtcServerRequestAddWebRtcServerId(builder *flatbuffers.Builder, webRtcServerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(webRtcServerId), 0)
}
func CreateWebRtcServerRequestAddListenInfos(builder *flatbuffers.Builder, listenInfos flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(listenInfos), 0)
}
func CreateWebRtcServerRequestStartListenInfosVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func CreateWebRtcServerRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
