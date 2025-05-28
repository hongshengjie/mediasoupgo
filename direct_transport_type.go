package mediasoupgo

import (
	"mediasoupgo/internal/events"
)

type DirectTransportOptions struct {
	MaxMessageSize uint32
	AppData        AppData
}

type DirectTransportDump struct {
	BaseTransportDump
}

type DirectTransportStat struct {
	BaseTransportStats
	Type string
}

type DirectTransportEvents struct {
	TransportEvents
	RTCP *events.Event1[[]byte]
}

type DirectTransportObserver events.EventEmmiter[DirectTransportObserverEvents]

type DirectTransportObserverEvents struct {
	TransportObserverEvents
	RTCP events.Event1[[]byte]
}

// DirectTransport interface definition
type DirectTransport interface {
	events.EventEmmiter[DirectTransportEvents]
	Transport
	// Embeds all methods from Transport interface

	// Transport type
	// Override: always returns "direct"
	Type() string

	// Observer
	// Override: returns DirectTransportObserver
	Observer() DirectTransportObserver

	// Dump DirectTransport
	// Override
	Dump() (DirectTransportDump, error)

	// Get DirectTransport stats
	// Override
	GetStats() ([]DirectTransportStat, error)

	// NO-OP method in DirectTransport
	// Override
	Connect() error

	// Send RTCP packet
	SendRtcp(rtcpPacket []byte)
}
