package mediasoupgo

import (
	"mediasoupgo/events"
)

type ProducerOptions struct {
	ID                   *string
	Kind                 MediaKind
	RTPParameters        RtpParameters
	Paused               *bool
	KeyFrameRequestDelay *uint32
	AppData              ProducerAppData
}

type (
	ProducerType    string
	ProducerAppData AppData
)

const (
	SimpleProducerType    ProducerType = "simple"
	SimulcastProducerType ProducerType = "simulcast"
	SvcProducerType       ProducerType = "svc"
)

type ProducerScore struct {
	EncodingIdx uint32
	SSRC        uint32
	RID         *string
	Score       byte
}

type ProducerVideoOrientation struct {
	Camera   bool
	Flip     bool
	Rotation int
}

type ProducerDump struct {
	ID              string
	Kind            string
	Type            ProducerType
	RTPParameters   RtpParameters
	RTPMapping      interface{}
	RTPStreams      interface{}
	TraceEventTypes []string
	Paused          bool
}

type ProducerStat struct {
	RtpStreamRecvStats
}

type ProducerTraceEventType string

const (
	RTPProducerTraceEventType      ProducerTraceEventType = "rtp"
	KeyframeProducerTraceEventType ProducerTraceEventType = "keyframe"
	NackProducerTraceEventType     ProducerTraceEventType = "nack"
	PliProducerTraceEventType      ProducerTraceEventType = "pli"
	FirProducerTraceEventType      ProducerTraceEventType = "fir"
	SrProducerTraceEventType       ProducerTraceEventType = "sr"
)

type ProducerTraceEventData struct {
	Type      ProducerTraceEventType
	Timestamp int
	Direction string // "in" or "out"
	Info      any
}

type ProducerEvents struct {
	TransportClose         struct{}
	Score                  *events.Event1[[]*ProducerScore]
	VideoOrientationChange *events.Event1[ProducerVideoOrientation]
	Trace                  *events.Event1[ProducerTraceEventData]
	AtClose                struct{}
}

type ProducerObserver events.EventEmmiter[ProducerObserverEvents]

type ProducerObserverEvents struct {
	Close                  struct{}
	Pause                  struct{}
	Resume                 struct{}
	Score                  *events.Event1[[]*ProducerScore]
	VideoOrientationChange *events.Event1[ProducerVideoOrientation]
	Trace                  *events.Event1[ProducerTraceEventData]
}

// Producer interface definition
type Producer interface {
	// EnhancedEventEmitter methods would be included here if defined
	// For now, we'll focus on Producer-specific methods
	events.EventEmmiter[ProducerEvents]
	// Producer id
	ID() string

	// Whether the Producer is closed
	Closed() bool

	// Media kind
	Kind() MediaKind

	// RTP parameters
	RTPParameters() RtpParameters

	// Producer type
	Type() ProducerType

	// Consumable RTP parameters
	ConsumableRTPParameters() RtpParameters

	// Whether the Producer is paused
	Paused() bool

	// Producer score list
	Score() []*ProducerScore

	// App custom data
	AppData() AppData
	SetAppData(appData AppData)

	// Observer
	Observer() ProducerObserver

	// Close the Producer
	Close()

	// Transport was closed
	TransportClosed()

	// Dump Producer
	Dump() (ProducerDump, error)

	// Get Producer stats
	GetStats() ([]ProducerStat, error)

	// Pause the Producer
	Pause() error

	// Resume the Producer
	Resume() error

	// Enable 'trace' event
	EnableTraceEvent(types []ProducerTraceEventType) error

	// Send RTP packet (just valid for Producers created on a DirectTransport)
	Send(rtpPacket []byte)
}
