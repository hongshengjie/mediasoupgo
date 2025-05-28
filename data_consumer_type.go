package mediasoupgo

import (
	"mediasoupgo/internal/events"
)

type DataConsumerOptions struct {
	DataProducerID    string
	Ordered           *bool
	MaxPacketLifeTime *uint16
	MaxRetransmits    *uint16
	Paused            *bool
	Subchannels       []uint16
	AppData           AppData
}

type DataConsumerType string

const (
	SCTPDataConsumerType   DataConsumerType = "sctp"
	DirectDataConsumerType DataConsumerType = "direct"
)

type DataConsumerDump struct {
	ID                         string
	Paused                     bool
	DataProducerPaused         bool
	Subchannels                []uint16
	DataProducerID             string
	Type                       DataConsumerType
	SCTPStreamParameters       *SctpStreamParameters
	Label                      string
	Protocol                   string
	BufferedAmountLowThreshold int
}

type DataConsumerStat struct {
	Type           string
	Timestamp      int
	Label          string
	Protocol       string
	MessagesSent   int
	BytesSent      int
	BufferedAmount int
}
type MessageItem struct {
	Ppid int
	Data []byte
}
type DataConsumerEvents struct {
	TransportClose      struct{}
	DataProducerClose   struct{}
	DataProducerPause   struct{}
	DataProducerResume  struct{}
	Message             *events.Event1[*MessageItem]
	SCTPSendBufferFull  struct{}
	BufferedAmountLow   uint32
	AtClose             struct{}
	AtDataProducerClose struct{}
}

type DataConsumerObserver events.EventEmmiter[DataConsumerObserverEvents]

type DataConsumerObserverEvents struct {
	Close  struct{}
	Pause  struct{}
	Resume struct{}
}

// DataConsumer interface definition
type DataConsumer interface {
	// EnhancedEventEmitter methods would be included here if defined
	// For now, we'll focus on DataConsumer-specific methods
	events.EventEmmiter[DataConsumerEvents]
	// DataConsumer id
	ID() string

	// Associated DataProducer id
	DataProducerID() string

	// Whether the DataConsumer is closed
	Closed() bool

	// DataConsumer type
	Type() DataConsumerType

	// SCTP stream parameters
	SCTPStreamParameters() *SctpStreamParameters // Using pointer for undefined possibility

	// DataChannel label
	Label() string

	// DataChannel protocol
	Protocol() string

	// Whether the DataConsumer is paused
	Paused() bool

	// Whether the associate DataProducer is paused
	DataProducerPaused() bool

	// Get current subchannels this data consumer is subscribed to
	Subchannels() []uint16

	// App custom data
	AppData() AppData
	SetAppData(appData AppData)

	// Observer
	Observer() DataConsumerObserver

	// Close the DataConsumer
	Close()

	// Transport was closed
	TransportClosed()

	// Dump DataConsumer
	Dump() (DataConsumerDump, error)

	// Get DataConsumer stats
	GetStats() ([]DataConsumerStat, error)

	// Pause the DataConsumer
	Pause() error

	// Resume the DataConsumer
	Resume() error

	// Set buffered amount low threshold
	SetBufferedAmountLowThreshold(threshold int) error

	// Get buffered amount size
	GetBufferedAmount() (int, error)

	// Send a message
	Send(message []byte, isString bool) error // message can be string or []byte

	// Set subchannels
	SetSubchannels(subchannels []uint16) error

	// Add a subchannel
	AddSubchannel(subchannel int) error

	// Remove a subchannel
	RemoveSubchannel(subchannel int) error
}
