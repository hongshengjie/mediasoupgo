package mediasoupgo

type SctpCapabilities struct {
	NumStreams NumSctpStreams
}

type NumSctpStreams struct {
	OS  uint16
	MIS uint16
}

type SctpParameters struct {
	Port           uint16
	OS             uint16
	MIS            uint16
	MaxMessageSize uint32
}

type SctpStreamParameters struct {
	StreamID          uint16  `json:"streamId"`
	Ordered           *bool   `json:"ordered"`
	MaxPacketLifeTime *uint16 `json:"maxPacketLifeTime"`
	MaxRetransmits    *uint16 `json:"maxRetransmits"`
}

type SctpParametersDump struct {
	Port               int
	OS                 int
	MIS                int
	MaxMessageSize     int
	SendBufferSize     int
	SCTPBufferedAmount int
	IsDataChannel      bool
}
