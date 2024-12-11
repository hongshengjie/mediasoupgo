// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Router

import (
	flatbuffers "github.com/google/flatbuffers/go"

	FBS__Common "mediasoupgo/FBS/Common"
)

type DumpResponseT struct {
	Id string `json:"id"`
	TransportIds []string `json:"transport_ids"`
	RtpObserverIds []string `json:"rtp_observer_ids"`
	MapProducerIdConsumerIds []*FBS__Common.StringStringArrayT `json:"map_producer_id_consumer_ids"`
	MapConsumerIdProducerId []*FBS__Common.StringStringT `json:"map_consumer_id_producer_id"`
	MapProducerIdObserverIds []*FBS__Common.StringStringArrayT `json:"map_producer_id_observer_ids"`
	MapDataProducerIdDataConsumerIds []*FBS__Common.StringStringArrayT `json:"map_data_producer_id_data_consumer_ids"`
	MapDataConsumerIdDataProducerId []*FBS__Common.StringStringT `json:"map_data_consumer_id_data_producer_id"`
}

func (t *DumpResponseT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	idOffset := flatbuffers.UOffsetT(0)
	if t.Id != "" {
		idOffset = builder.CreateString(t.Id)
	}
	transportIdsOffset := flatbuffers.UOffsetT(0)
	if t.TransportIds != nil {
		transportIdsLength := len(t.TransportIds)
		transportIdsOffsets := make([]flatbuffers.UOffsetT, transportIdsLength)
		for j := 0; j < transportIdsLength; j++ {
			transportIdsOffsets[j] = builder.CreateString(t.TransportIds[j])
		}
		DumpResponseStartTransportIdsVector(builder, transportIdsLength)
		for j := transportIdsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(transportIdsOffsets[j])
		}
		transportIdsOffset = builder.EndVector(transportIdsLength)
	}
	rtpObserverIdsOffset := flatbuffers.UOffsetT(0)
	if t.RtpObserverIds != nil {
		rtpObserverIdsLength := len(t.RtpObserverIds)
		rtpObserverIdsOffsets := make([]flatbuffers.UOffsetT, rtpObserverIdsLength)
		for j := 0; j < rtpObserverIdsLength; j++ {
			rtpObserverIdsOffsets[j] = builder.CreateString(t.RtpObserverIds[j])
		}
		DumpResponseStartRtpObserverIdsVector(builder, rtpObserverIdsLength)
		for j := rtpObserverIdsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(rtpObserverIdsOffsets[j])
		}
		rtpObserverIdsOffset = builder.EndVector(rtpObserverIdsLength)
	}
	mapProducerIdConsumerIdsOffset := flatbuffers.UOffsetT(0)
	if t.MapProducerIdConsumerIds != nil {
		mapProducerIdConsumerIdsLength := len(t.MapProducerIdConsumerIds)
		mapProducerIdConsumerIdsOffsets := make([]flatbuffers.UOffsetT, mapProducerIdConsumerIdsLength)
		for j := 0; j < mapProducerIdConsumerIdsLength; j++ {
			mapProducerIdConsumerIdsOffsets[j] = t.MapProducerIdConsumerIds[j].Pack(builder)
		}
		DumpResponseStartMapProducerIdConsumerIdsVector(builder, mapProducerIdConsumerIdsLength)
		for j := mapProducerIdConsumerIdsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(mapProducerIdConsumerIdsOffsets[j])
		}
		mapProducerIdConsumerIdsOffset = builder.EndVector(mapProducerIdConsumerIdsLength)
	}
	mapConsumerIdProducerIdOffset := flatbuffers.UOffsetT(0)
	if t.MapConsumerIdProducerId != nil {
		mapConsumerIdProducerIdLength := len(t.MapConsumerIdProducerId)
		mapConsumerIdProducerIdOffsets := make([]flatbuffers.UOffsetT, mapConsumerIdProducerIdLength)
		for j := 0; j < mapConsumerIdProducerIdLength; j++ {
			mapConsumerIdProducerIdOffsets[j] = t.MapConsumerIdProducerId[j].Pack(builder)
		}
		DumpResponseStartMapConsumerIdProducerIdVector(builder, mapConsumerIdProducerIdLength)
		for j := mapConsumerIdProducerIdLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(mapConsumerIdProducerIdOffsets[j])
		}
		mapConsumerIdProducerIdOffset = builder.EndVector(mapConsumerIdProducerIdLength)
	}
	mapProducerIdObserverIdsOffset := flatbuffers.UOffsetT(0)
	if t.MapProducerIdObserverIds != nil {
		mapProducerIdObserverIdsLength := len(t.MapProducerIdObserverIds)
		mapProducerIdObserverIdsOffsets := make([]flatbuffers.UOffsetT, mapProducerIdObserverIdsLength)
		for j := 0; j < mapProducerIdObserverIdsLength; j++ {
			mapProducerIdObserverIdsOffsets[j] = t.MapProducerIdObserverIds[j].Pack(builder)
		}
		DumpResponseStartMapProducerIdObserverIdsVector(builder, mapProducerIdObserverIdsLength)
		for j := mapProducerIdObserverIdsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(mapProducerIdObserverIdsOffsets[j])
		}
		mapProducerIdObserverIdsOffset = builder.EndVector(mapProducerIdObserverIdsLength)
	}
	mapDataProducerIdDataConsumerIdsOffset := flatbuffers.UOffsetT(0)
	if t.MapDataProducerIdDataConsumerIds != nil {
		mapDataProducerIdDataConsumerIdsLength := len(t.MapDataProducerIdDataConsumerIds)
		mapDataProducerIdDataConsumerIdsOffsets := make([]flatbuffers.UOffsetT, mapDataProducerIdDataConsumerIdsLength)
		for j := 0; j < mapDataProducerIdDataConsumerIdsLength; j++ {
			mapDataProducerIdDataConsumerIdsOffsets[j] = t.MapDataProducerIdDataConsumerIds[j].Pack(builder)
		}
		DumpResponseStartMapDataProducerIdDataConsumerIdsVector(builder, mapDataProducerIdDataConsumerIdsLength)
		for j := mapDataProducerIdDataConsumerIdsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(mapDataProducerIdDataConsumerIdsOffsets[j])
		}
		mapDataProducerIdDataConsumerIdsOffset = builder.EndVector(mapDataProducerIdDataConsumerIdsLength)
	}
	mapDataConsumerIdDataProducerIdOffset := flatbuffers.UOffsetT(0)
	if t.MapDataConsumerIdDataProducerId != nil {
		mapDataConsumerIdDataProducerIdLength := len(t.MapDataConsumerIdDataProducerId)
		mapDataConsumerIdDataProducerIdOffsets := make([]flatbuffers.UOffsetT, mapDataConsumerIdDataProducerIdLength)
		for j := 0; j < mapDataConsumerIdDataProducerIdLength; j++ {
			mapDataConsumerIdDataProducerIdOffsets[j] = t.MapDataConsumerIdDataProducerId[j].Pack(builder)
		}
		DumpResponseStartMapDataConsumerIdDataProducerIdVector(builder, mapDataConsumerIdDataProducerIdLength)
		for j := mapDataConsumerIdDataProducerIdLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(mapDataConsumerIdDataProducerIdOffsets[j])
		}
		mapDataConsumerIdDataProducerIdOffset = builder.EndVector(mapDataConsumerIdDataProducerIdLength)
	}
	DumpResponseStart(builder)
	DumpResponseAddId(builder, idOffset)
	DumpResponseAddTransportIds(builder, transportIdsOffset)
	DumpResponseAddRtpObserverIds(builder, rtpObserverIdsOffset)
	DumpResponseAddMapProducerIdConsumerIds(builder, mapProducerIdConsumerIdsOffset)
	DumpResponseAddMapConsumerIdProducerId(builder, mapConsumerIdProducerIdOffset)
	DumpResponseAddMapProducerIdObserverIds(builder, mapProducerIdObserverIdsOffset)
	DumpResponseAddMapDataProducerIdDataConsumerIds(builder, mapDataProducerIdDataConsumerIdsOffset)
	DumpResponseAddMapDataConsumerIdDataProducerId(builder, mapDataConsumerIdDataProducerIdOffset)
	return DumpResponseEnd(builder)
}

func (rcv *DumpResponse) UnPackTo(t *DumpResponseT) {
	t.Id = string(rcv.Id())
	transportIdsLength := rcv.TransportIdsLength()
	t.TransportIds = make([]string, transportIdsLength)
	for j := 0; j < transportIdsLength; j++ {
		t.TransportIds[j] = string(rcv.TransportIds(j))
	}
	rtpObserverIdsLength := rcv.RtpObserverIdsLength()
	t.RtpObserverIds = make([]string, rtpObserverIdsLength)
	for j := 0; j < rtpObserverIdsLength; j++ {
		t.RtpObserverIds[j] = string(rcv.RtpObserverIds(j))
	}
	mapProducerIdConsumerIdsLength := rcv.MapProducerIdConsumerIdsLength()
	t.MapProducerIdConsumerIds = make([]*FBS__Common.StringStringArrayT, mapProducerIdConsumerIdsLength)
	for j := 0; j < mapProducerIdConsumerIdsLength; j++ {
		x := FBS__Common.StringStringArray{}
		rcv.MapProducerIdConsumerIds(&x, j)
		t.MapProducerIdConsumerIds[j] = x.UnPack()
	}
	mapConsumerIdProducerIdLength := rcv.MapConsumerIdProducerIdLength()
	t.MapConsumerIdProducerId = make([]*FBS__Common.StringStringT, mapConsumerIdProducerIdLength)
	for j := 0; j < mapConsumerIdProducerIdLength; j++ {
		x := FBS__Common.StringString{}
		rcv.MapConsumerIdProducerId(&x, j)
		t.MapConsumerIdProducerId[j] = x.UnPack()
	}
	mapProducerIdObserverIdsLength := rcv.MapProducerIdObserverIdsLength()
	t.MapProducerIdObserverIds = make([]*FBS__Common.StringStringArrayT, mapProducerIdObserverIdsLength)
	for j := 0; j < mapProducerIdObserverIdsLength; j++ {
		x := FBS__Common.StringStringArray{}
		rcv.MapProducerIdObserverIds(&x, j)
		t.MapProducerIdObserverIds[j] = x.UnPack()
	}
	mapDataProducerIdDataConsumerIdsLength := rcv.MapDataProducerIdDataConsumerIdsLength()
	t.MapDataProducerIdDataConsumerIds = make([]*FBS__Common.StringStringArrayT, mapDataProducerIdDataConsumerIdsLength)
	for j := 0; j < mapDataProducerIdDataConsumerIdsLength; j++ {
		x := FBS__Common.StringStringArray{}
		rcv.MapDataProducerIdDataConsumerIds(&x, j)
		t.MapDataProducerIdDataConsumerIds[j] = x.UnPack()
	}
	mapDataConsumerIdDataProducerIdLength := rcv.MapDataConsumerIdDataProducerIdLength()
	t.MapDataConsumerIdDataProducerId = make([]*FBS__Common.StringStringT, mapDataConsumerIdDataProducerIdLength)
	for j := 0; j < mapDataConsumerIdDataProducerIdLength; j++ {
		x := FBS__Common.StringString{}
		rcv.MapDataConsumerIdDataProducerId(&x, j)
		t.MapDataConsumerIdDataProducerId[j] = x.UnPack()
	}
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

func (rcv *DumpResponse) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *DumpResponse) TransportIds(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *DumpResponse) TransportIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) RtpObserverIds(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *DumpResponse) RtpObserverIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) MapProducerIdConsumerIds(obj *FBS__Common.StringStringArray, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *DumpResponse) MapProducerIdConsumerIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) MapConsumerIdProducerId(obj *FBS__Common.StringString, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *DumpResponse) MapConsumerIdProducerIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) MapProducerIdObserverIds(obj *FBS__Common.StringStringArray, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *DumpResponse) MapProducerIdObserverIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) MapDataProducerIdDataConsumerIds(obj *FBS__Common.StringStringArray, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *DumpResponse) MapDataProducerIdDataConsumerIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *DumpResponse) MapDataConsumerIdDataProducerId(obj *FBS__Common.StringString, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *DumpResponse) MapDataConsumerIdDataProducerIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func DumpResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(8)
}
func DumpResponseAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func DumpResponseAddTransportIds(builder *flatbuffers.Builder, transportIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(transportIds), 0)
}
func DumpResponseStartTransportIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddRtpObserverIds(builder *flatbuffers.Builder, rtpObserverIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(rtpObserverIds), 0)
}
func DumpResponseStartRtpObserverIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddMapProducerIdConsumerIds(builder *flatbuffers.Builder, mapProducerIdConsumerIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(mapProducerIdConsumerIds), 0)
}
func DumpResponseStartMapProducerIdConsumerIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddMapConsumerIdProducerId(builder *flatbuffers.Builder, mapConsumerIdProducerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(mapConsumerIdProducerId), 0)
}
func DumpResponseStartMapConsumerIdProducerIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddMapProducerIdObserverIds(builder *flatbuffers.Builder, mapProducerIdObserverIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(mapProducerIdObserverIds), 0)
}
func DumpResponseStartMapProducerIdObserverIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddMapDataProducerIdDataConsumerIds(builder *flatbuffers.Builder, mapDataProducerIdDataConsumerIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(mapDataProducerIdDataConsumerIds), 0)
}
func DumpResponseStartMapDataProducerIdDataConsumerIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseAddMapDataConsumerIdDataProducerId(builder *flatbuffers.Builder, mapDataConsumerIdDataProducerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(mapDataConsumerIdDataProducerId), 0)
}
func DumpResponseStartMapDataConsumerIdDataProducerIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func DumpResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
