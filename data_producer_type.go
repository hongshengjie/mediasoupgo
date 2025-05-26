package mediasoupgo

import (
	"mediasoupgo/events"
)

type DataProducerOptions struct {
	ID                   *string
	SCTPStreamParameters *SctpStreamParameters
	Label                *string
	Protocol             *string
	Paused               *bool
	AppData              AppData
}

type DataProducerType string

const (
	SCTPDataProducerType   DataProducerType = "sctp"
	DirectDataProducerType DataProducerType = "direct"
)

type DataProducerDump struct {
	ID                   string
	Paused               bool
	Type                 DataProducerType
	SCTPStreamParameters *SctpStreamParameters
	Label                string
	Protocol             string
}

type DataProducerStat struct {
	Type             string
	Timestamp        int
	Label            string
	Protocol         string
	MessagesReceived int
	BytesReceived    int
}

type DataProducerEvents struct {
	TransportClose struct{}
	AtClose        struct{}
}

type DataProducerObserver events.EventEmmiter[DataProducerObserverEvents]

type DataProducerObserverEvents struct {
	Close  struct{}
	Pause  struct{}
	Resume struct{}
}

// DataProducer struct definition
type DataProducer interface {
	// EnhancedEventEmitter methods would be included here if defined
	// For now, we'll focus on DataProducer-specific methods
	events.EventEmmiter[DataProducerEvents]
	// DataProducer id
	ID() string

	// Whether the DataProducer is closed
	Closed() bool

	// DataProducer type
	Type() DataProducerType

	// SCTP stream parameters
	SCTPStreamParameters() *SctpStreamParameters // Using pointer for undefined possibility

	// DataChannel label
	Label() string

	// DataChannel protocol
	Protocol() string

	// Whether the DataProducer is paused
	Paused() bool

	// App custom data
	AppData() AppData
	SetAppData(appData AppData)

	// Observer
	Observer() DataProducerObserver

	// Close the DataProducer
	Close()

	// Transport was closed
	TransportClosed()

	// Dump DataProducer
	Dump() (DataProducerDump, error)

	// Get DataProducer stats
	GetStats() ([]DataProducerStat, error)

	// Pause the DataProducer
	Pause() error

	// Resume the DataProducer
	Resume() error

	// Send data (just valid for DataProducers created on a DirectTransport)
	Send(message []byte, isString bool, subchannels []uint16, requiredSubchannel uint16)
}
