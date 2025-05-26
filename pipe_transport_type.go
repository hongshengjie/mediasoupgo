package mediasoupgo

import "mediasoupgo/events"

type PipeTransportOptions struct {
	EnableSctp         *bool
	NumSctpStreams     *NumSctpStreams
	MaxSctpMessageSize *int
	SctpSendBufferSize *int
	EnableRtx          *bool
	EnableSrtp         *bool
	AppData            AppData
	ListenInfo         *TransportListenInfo
	ListenIp           *interface{}
	Port               *int
}

type PipeConsumerOptions struct {
	ProducerID string
	AppData    AppData
}

type PipeTransportDump struct {
	BaseTransportDump
	Tuple          TransportTuple
	RTX            bool
	SRTPParameters *SrtpParameters
}

type PipeTransportStat struct {
	BaseTransportStats
	Type  string
	Tuple TransportTuple
}

type PipeTransportEvents struct {
	TransportEvents
	SctpStateChange *events.Event1[SctpState]
}

type PipeTransportObserver events.EventEmmiter[PipeTransportObserverEvents]

type PipeTransportObserverEvents struct {
	TransportObserverEvents
	SctpStateChange *events.Event1[SctpState]
}

type PipeTransportConnectParams struct {
	IP             string
	Port           uint16
	SRTPParameters *SrtpParameters // Optional field using pointer
}

// PipeTransport interface definition
type PipeTransport interface {
	events.EventEmmiter[PipeTransportEvents]
	Transport
	// Transport type
	// Override: always returns "pipe"
	Type() string

	// Observer
	// Override: returns PipeTransportObserver
	Observer() PipeTransportObserver

	// PipeTransport tuple
	Tuple() TransportTuple

	// SCTP parameters
	SCTPParameters() *SctpParameters // Using pointer for undefined possibility

	// SCTP state
	SCTPState() *SctpState // Using pointer for undefined possibility

	// SRTP parameters
	SRTPParameters() *SrtpParameters // Using pointer for undefined possibility

	// Dump PipeTransport
	// Override
	Dump() (PipeTransportDump, error)

	// Get PipeTransport stats
	// Override
	GetStats() ([]PipeTransportStat, error)

	// Provide the PipeTransport remote parameters
	// Override
	Connect(params PipeTransportConnectParams) error

	// Create a pipe Consumer
	// Override
	Consume(options *ConsumerOptions) (Consumer, error)
}
