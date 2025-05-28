package mediasoupgo

import (
	"mediasoupgo/internal/events"
)

// TransportType and other types remain unchanged
type TransportType string

const (
	WebRTCTransportType TransportType = "webrtc"
	PlainTransportType  TransportType = "plain"
	PipeTransportType   TransportType = "pipe"
	DirectTransportType TransportType = "direct"
)

type TransportListenInfo struct {
	Protocol         TransportProtocol
	IP               string
	AnnouncedIP      *string
	AnnouncedAddress *string
	Port             *uint16
	PortRange        *TransportPortRange
	Flags            *TransportSocketFlags
	SendBufferSize   *uint32
	RecvBufferSize   *uint32
}

type TransportListenIp struct {
	IP          string
	AnnouncedIP *string
}

type TransportProtocol string

const (
	UDPTransportProtocol TransportProtocol = "udp"
	TCPTransportProtocol TransportProtocol = "tcp"
)

type TransportPortRange struct {
	Min uint16
	Max uint16
}

type TransportSocketFlags struct {
	IPV6Only     bool
	UDPReusePort bool
}

type TransportTuple struct {
	LocalIP      string
	LocalAddress string
	LocalPort    uint16
	RemoteIP     *string
	RemotePort   *uint16
	Protocol     TransportProtocol
}

type SctpState string

const (
	NewSctpState        SctpState = "new"
	ConnectingSctpState SctpState = "connecting"
	ConnectedSctpState  SctpState = "connected"
	FailedSctpState     SctpState = "failed"
	ClosedSctpState     SctpState = "closed"
)

type RtpListenerDump struct {
	SSRCTable []struct {
		Key   int
		Value string
	}
	MIDTable []struct {
		Key   int
		Value string
	}
	RIDTable []struct {
		Key   int
		Value string
	}
}

type SctpListenerDump struct {
	StreamIDTable []struct {
		Key   int
		Value string
	}
}

type RecvRtpHeaderExtensions struct {
	MID               *int
	RID               *int
	RRID              *int
	AbsSendTime       *int
	TransportWideCc01 *int
}

type BaseTransportDump struct {
	ID                string
	ProducerIDs       []string
	ConsumerIDs       []string
	MapSSRCConsumerID []struct {
		Key   int
		Value string
	}
	MapRTXSSRCConsumerID []struct {
		Key   int
		Value string
	}
	RecvRTPHeaderExtensions RecvRtpHeaderExtensions
	RTPListener             RtpListenerDump
	MaxMessageSize          int
	DataProducerIDs         []string
	DataConsumerIDs         []string
	SCTPParameters          *SctpParameters
	SCTPState               *SctpState
	SCTPListener            *SctpListenerDump
	TraceEventTypes         []string
}

type BaseTransportStats struct {
	TransportID              string
	Timestamp                int
	SCTPState                *SctpState
	BytesReceived            int
	RecvBitrate              int
	BytesSent                int
	SendBitrate              int
	RTPBytesReceived         int
	RTPRecvBitrate           int
	RTPBytesSent             int
	RTPSendBitrate           int
	RTXBytesReceived         int
	RTXRecvBitrate           int
	RTXBytesSent             int
	RTXSendBitrate           int
	ProbationBytesSent       int
	ProbationSendBitrate     int
	AvailableOutgoingBitrate *int
	AvailableIncomingBitrate *int
	MaxIncomingBitrate       *int
	MaxOutgoingBitrate       *int
	MinOutgoingBitrate       *int
	RTPPacketLossReceived    *int
	RTPPacketLossSent        *int
}

type TransportTraceEventType string

const (
	ProbationTransportTraceEventType TransportTraceEventType = "probation"
	BWETransportTraceEventType       TransportTraceEventType = "bwe"
)

type TransportTraceEventData struct {
	Type      TransportTraceEventType
	Timestamp uint64
	Direction string
	Info      any
}

// Supporting types (these would need to be defined based on your specific needs)
type (
	TransportAppData AppData

	TransportObserver events.EventEmmiter[TransportObserverEvents]

	TransportObserverEvents struct {
		Close           struct{}
		Newproducer     *events.Event1[Producer]
		Newconsumer     *events.Event1[Consumer]
		Newdataproducer *events.Event1[DataProducer]
		Newdataconsumer *events.Event1[DataConsumer]
		Trace           *events.Event1[TransportTraceEventData]
	}
	TransportEvents struct {
		Routerclose       struct{}
		Listrnserverclose struct{}
		Trace             *events.Event1[TransportTraceEventData]
		Close             struct{}

		Newproducer       *events.Event1[Producer]
		Producerclose     *events.Event1[Producer]
		Newdataproducer   *events.Event1[DataProducer]
		Dataproducerclose *events.Event1[DataProducer]
		Listenserverclose struct{}
	}
)

type Transport interface {
	ID() string
	Closed() bool
	AppData() TransportAppData
	SetAppData(appData TransportAppData)
	Close()
	RouterClosed()
	ListenServerClosed()
	SetMaxIncomingBitrate(bitrate int) error
	SetMaxOutgoingBitrate(bitrate int) error
	SetMinOutgoingBitrate(bitrate int) error
	Produce(options *ProducerOptions) (Producer, error)
	Consume(options *ConsumerOptions) (Consumer, error)
	ProduceData(options *DataProducerOptions) (DataProducer, error)
	ConsumeData(options *DataConsumerOptions) (DataConsumer, error)
	EnableTraceEvent(types []TransportTraceEventType) error
}
