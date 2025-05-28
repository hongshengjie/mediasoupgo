package mediasoupgo

import "mediasoupgo/internal/events"

// RtpObserverType represents the possible types of RTP observers
type RtpObserverType string

// Constants for RtpObserverType
const (
	AudioLevelRtpObserverType    RtpObserverType = "audiolevel"
	ActiveSpeakerRtpObserverType RtpObserverType = "activespeaker"
)

type RtpObserverAppData AppData

type RtpObserverEvents struct {
	Routerclose struct{}
	AtClose     struct{}
}

type RtpObserverObserverEvents struct {
	Close          struct{}
	Pause          struct{}
	Resume         struct{}
	Addproducer    *events.Event1[Producer]
	Removeproducer *events.Event1[Producer]
}

type RtpObserverObserver events.EventEmmiter[RtpObserverObserverEvents]

type RtpObserver interface {
	ID() string

	Closed() bool

	Paused() bool

	AppData() AppData
	SetAppData(appData AppData)

	Close()

	RouterClosed()

	Pause() error

	Resume() error

	AddProducer(params struct{ ProducerID string }) error

	RemoveProducer(params struct{ ProducerID string }) error
}
