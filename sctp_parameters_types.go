package mediasoupgo

type SctpCapabilities struct {
	NumStreams NumSctpStreams
}

type NumSctpStreams struct {
	OS  int // Initially requested number of outgoing SCTP streams.
	MIS int // Maximum number of incoming SCTP streams.
}

type SctpParameters struct {
	Port           int // Must always equal 5000.
	OS             int // Initially requested number of outgoing SCTP streams.
	MIS            int // Maximum number of incoming SCTP streams.
	MaxMessageSize int // Maximum allowed size for SCTP messages.
}

type SctpStreamParameters struct {
	StreamId          int   // SCTP stream id.
	Ordered           *bool // Whether data messages must be received in order. Default true.
	MaxPacketLifeTime *int  // When ordered is false indicates the time (in milliseconds) after which a SCTP packet will stop being retransmitted.
	MaxRetransmits    *int  // When ordered is false indicates the maximum number of times a packet will be retransmitted.
}

type SctpParametersDump struct {
	Port               int  // Port
	OS                 int  // OS
	MIS                int  // MIS
	MaxMessageSize     int  // MaxMessageSize
	SendBufferSize     int  // SendBufferSize
	SctpBufferedAmount int  // sctpBufferedAmount
	IsDataChannel      bool // isDataChannel
}
