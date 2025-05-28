package mediasoupgo

import (
	"mediasoupgo/events"
)

type PlainTransportOptions struct {
	RTCPMux            *bool
	Comedia            *bool
	EnableSctp         *bool
	NumSctpStreams     *NumSctpStreams
	MaxSctpMessageSize *uint32
	SctpSendBufferSize *uint32
	EnableSrtp         *bool
	SRTPCryptoSuite    *SrtpCryptoSuite
	AppData            AppData
	ListenInfo         *TransportListenInfo
	RTCPListenInfo     *TransportListenInfo
	ListenIp           *string
	Port               *uint16
}

type PlainTransportDump struct {
	BaseTransportDump
	RTCPMux        bool
	Comedia        bool
	Tuple          TransportTuple
	RTCPTuple      *TransportTuple
	SRTPParameters *SrtpParameters
}

type PlainTransportStat struct {
	BaseTransportStats
	Type      string
	RTCPMux   bool
	Comedia   bool
	Tuple     TransportTuple
	RTCPTuple *TransportTuple
}

type PlainTransportEvents struct {
	TransportEvents
	Tuple           *events.Event1[TransportTuple]
	RTCPTuple       *events.Event1[TransportTuple]
	SctpStateChange *events.Event1[SctpState]
}

type PlainTransportObserver events.EventEmmiter[PlainTransportObserverEvents]

type PlainTransportObserverEvents struct {
	TransportObserverEvents
	Tuple           *events.Event1[TransportTuple]
	RTCPTuple       *events.Event1[TransportTuple]
	SctpStateChange *events.Event1[SctpState]
}

type PlainTransportConnectParams struct {
	IP             *string         // Optional field using pointer
	Port           *uint16         // Optional field using pointer
	RTCPPort       *uint16         // Optional field using pointer
	SRTPParameters *SrtpParameters // Optional field using pointer
}

type (
	// PlainTransport interface definition
	PlainTransport interface {
		events.EventEmmiter[PlainTransportEvents]
		Transport
		// Transport type
		// Override: always returns "plain"
		Type() string

		// Observer
		// Override: returns PlainTransportObserver
		Observer() PlainTransportObserver

		// PlainTransport tuple
		Tuple() TransportTuple

		// PlainTransport RTCP tuple
		RTCPTuple() *TransportTuple // Using pointer for undefined possibility

		// SCTP parameters
		SCTPParameters() *SctpParameters // Using pointer for undefined possibility

		// SCTP state
		SCTPState() SctpState // Using pointer for undefined possibility

		// SRTP parameters
		SRTPParameters() *SrtpParameters // Using pointer for undefined possibility

		// Dump PlainTransport
		// Override
		Dump() (PlainTransportDump, error)

		// Get PlainTransport stats
		// Override
		GetStats() ([]PlainTransportStat, error)

		// Provide the PlainTransport remote parameters
		// Override
		Connect(params *PlainTransportConnectParams) error
	}
)
