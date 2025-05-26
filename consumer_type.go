package mediasoupgo

import (
	"bytes"
	"mediasoupgo/events"
)

type ConsumerOptions struct {
	ProducerID      string
	RTPCapabilities RtpCapabilities
	Paused          *bool
	MID             *string
	PreferredLayers *ConsumerLayers
	EnableRtx       *bool
	IgnoreDtx       *bool
	Pipe            *bool
	AppData         AppData
}

type ConsumerType string
type ConsumerAppData AppData

const (
	SimpleConsumerType    ConsumerType = "simple"
	SimulcastConsumerType ConsumerType = "simulcast"
	SvcConsumerType       ConsumerType = "svc"
	PipeConsumerType      ConsumerType = "pipe"
)

type ConsumerScore struct {
	Score          byte
	ProducerScore  byte
	ProducerScores []byte
}

type ConsumerLayers struct {
	SpatialLayer  byte
	TemporalLayer *byte
}

type ConsumerDump interface {
	// This is a union type in TypeScript, represented as an interface in Go
	GetType() string
}

type SimpleConsumerDump struct {
	BaseConsumerDump
	Type      string
	RTPStream RtpStreamDump
}

func (s SimpleConsumerDump) GetType() string { return s.Type }

type SimulcastConsumerDump struct {
	BaseConsumerDump
	Type                   string
	RTPStream              RtpStreamDump
	PreferredSpatialLayer  int
	TargetSpatialLayer     int
	CurrentSpatialLayer    int
	PreferredTemporalLayer int
	TargetTemporalLayer    int
	CurrentTemporalLayer   int
}

func (s SimulcastConsumerDump) GetType() string { return s.Type }

type SvcConsumerDump SimulcastConsumerDump

func (s SvcConsumerDump) GetType() string { return s.Type }

type PipeConsumerDump struct {
	BaseConsumerDump
	Type       string
	RTPStreams []RtpStreamDump
}

func (s PipeConsumerDump) GetType() string { return s.Type }

type BaseConsumerDump struct {
	ID                         string
	ProducerID                 string
	Kind                       MediaKind
	RTPParameters              RtpParameters
	ConsumableRTPEodings       []RtpEncodingParameters
	SupportedCodecPayloadTypes []int
	TraceEventTypes            []string
	Paused                     bool
	ProducerPaused             bool
	Priority                   int
}

type RtpStreamDump struct {
	Params    RtpStreamParametersDump
	Score     int
	RTXStream *RtxStreamDump
}

type RtpStreamParametersDump struct {
	EncodingIdx    int
	SSRC           int
	PayloadType    int
	MimeType       string
	ClockRate      int
	RID            *string
	Cname          string
	RTXSSRC        *int
	RTXPayloadType *int
	UseNack        bool
	UsePli         bool
	UseFir         bool
	UseInBandFec   bool
	UseDtx         bool
	SpatialLayers  int
	TemporalLayers int
}

type RtxStreamDump struct {
	Params RtxStreamParameters
}

type RtxStreamParameters struct {
	SSRC        int
	PayloadType int
	MimeType    string
	ClockRate   int
	RRID        *string
	Cname       string
}

type ConsumerStat struct {
	// Assuming RtpStreamSendStats is a struct, placeholder here
	RtpStreamSendStats
}

type ConsumerTraceEventType string

const (
	RTPConsumerTraceEventType      ConsumerTraceEventType = "rtp"
	KeyframeConsumerTraceEventType ConsumerTraceEventType = "keyframe"
	NackConsumerTraceEventType     ConsumerTraceEventType = "nack"
	PliConsumerTraceEventType      ConsumerTraceEventType = "pli"
	FirConsumerTraceEventType      ConsumerTraceEventType = "fir"
)

type ConsumerTraceEventData struct {
	Type      ConsumerTraceEventType
	Timestamp int
	Direction string // "in" or "out"
	Info      any
}

type ConsumerEvents struct {
	TransportClose  struct{}
	ProducerClose   struct{}
	ProducerPause   struct{}
	ProducerResume  struct{}
	Score           *events.Event1[ConsumerScore]
	LayersChange    *events.Event1[ConsumerLayers]
	Trace           *events.Event1[ConsumerTraceEventData]
	RTP             bytes.Buffer
	AtClose         struct{}
	AtProducerClose struct{}
}

type (
	ConsumerObserver       events.EventEmmiter[ConsumerObserverEvents]
	ConsumerObserverEvents struct {
		Close        struct{}
		Pause        struct{}
		Resume       struct{}
		Score        *events.Event1[ConsumerScore]
		LayersChange *events.Event1[ConsumerLayers]
		Trace        *events.Event1[ConsumerTraceEventData]
	}
)

// Consumer interface definition
type Consumer interface {
	// EnhancedEventEmitter methods would be included here if defined
	// For now, we'll focus on Consumer-specific methods
	events.EventEmmiter[ConsumerEvents]
	// Consumer id
	ID() string

	// Associated Producer id
	ProducerID() string

	// Whether the Consumer is closed
	Closed() bool

	// Media kind
	Kind() MediaKind

	// RTP parameters
	RTPParameters() RtpParameters

	// Consumer type
	Type() ConsumerType

	// Whether the Consumer is paused
	Paused() bool

	// Whether the associate Producer is paused
	ProducerPaused() bool

	// Current priority
	Priority() byte

	// Consumer score
	Score() ConsumerScore

	// Preferred video layers
	PreferredLayers() *ConsumerLayers // Using pointer for undefined possibility

	// Current video layers
	CurrentLayers() *ConsumerLayers // Using pointer for undefined possibility

	// App custom data
	AppData() AppData
	SetAppData(appData AppData)

	// Observer
	Observer() ConsumerObserver

	// Close the Consumer
	Close()

	// Transport was closed
	TransportClosed()

	// Dump Consumer
	Dump() (ConsumerDump, error)

	// Get Consumer stats
	GetStats() ([]interface{}, error) // Could be more specific with a union type

	// Pause the Consumer
	Pause() error

	// Resume the Consumer
	Resume() error

	// Set preferred video layers
	SetPreferredLayers(layers ConsumerLayers) error

	// Set priority
	SetPriority(priority byte) error

	// Unset priority
	UnsetPriority() error

	// Request a key frame to the Producer
	RequestKeyFrame() error

	// Enable 'trace' event
	EnableTraceEvent(types []ConsumerTraceEventType) error
}
