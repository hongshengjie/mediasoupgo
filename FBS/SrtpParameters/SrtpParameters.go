// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package SrtpParameters

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SrtpParametersT struct {
	CryptoSuite SrtpCryptoSuite `json:"crypto_suite"`
	KeyBase64 string `json:"key_base64"`
}

func (t *SrtpParametersT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	keyBase64Offset := flatbuffers.UOffsetT(0)
	if t.KeyBase64 != "" {
		keyBase64Offset = builder.CreateString(t.KeyBase64)
	}
	SrtpParametersStart(builder)
	SrtpParametersAddCryptoSuite(builder, t.CryptoSuite)
	SrtpParametersAddKeyBase64(builder, keyBase64Offset)
	return SrtpParametersEnd(builder)
}

func (rcv *SrtpParameters) UnPackTo(t *SrtpParametersT) {
	t.CryptoSuite = rcv.CryptoSuite()
	t.KeyBase64 = string(rcv.KeyBase64())
}

func (rcv *SrtpParameters) UnPack() *SrtpParametersT {
	if rcv == nil {
		return nil
	}
	t := &SrtpParametersT{}
	rcv.UnPackTo(t)
	return t
}

type SrtpParameters struct {
	_tab flatbuffers.Table
}

func GetRootAsSrtpParameters(buf []byte, offset flatbuffers.UOffsetT) *SrtpParameters {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SrtpParameters{}
	x.Init(buf, n+offset)
	return x
}

func FinishSrtpParametersBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsSrtpParameters(buf []byte, offset flatbuffers.UOffsetT) *SrtpParameters {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SrtpParameters{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedSrtpParametersBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *SrtpParameters) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SrtpParameters) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SrtpParameters) CryptoSuite() SrtpCryptoSuite {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return SrtpCryptoSuite(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *SrtpParameters) MutateCryptoSuite(n SrtpCryptoSuite) bool {
	return rcv._tab.MutateByteSlot(4, byte(n))
}

func (rcv *SrtpParameters) KeyBase64() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func SrtpParametersStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func SrtpParametersAddCryptoSuite(builder *flatbuffers.Builder, cryptoSuite SrtpCryptoSuite) {
	builder.PrependByteSlot(0, byte(cryptoSuite), 0)
}
func SrtpParametersAddKeyBase64(builder *flatbuffers.Builder, keyBase64 flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(keyBase64), 0)
}
func SrtpParametersEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}