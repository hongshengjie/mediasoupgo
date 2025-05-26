package mediasoupgo

import "mediasoupgo/events"

// ActiveSpeakerObserver interface definition
type ActiveSpeakerObserver interface {
	events.EventEmmiter[ActiveSpeakerObserverEvents]
	RtpObserver
	Type() string

	Observer() ActiveSpeakerObserverObserver
}

type ActiveSpeakerObserverOptions struct {
	Interval int
	AppData  ActiveSpeakerObserverAppData
}

type (
	ActiveSpeakerObserverAppData         = AppData
	ActiveSpeakerObserverDominantSpeaker struct {
		Producer Producer
	}
)

type ActiveSpeakerObserverEvents struct {
	RtpObserverEvents
	Dominantspeaker *events.Event1[ActiveSpeakerObserverDominantSpeaker]
}

type ActiveSpeakerObserverObserver events.EventEmmiter[ActiveSpeakerObserverObserverEvents]

type ActiveSpeakerObserverObserverEvents struct {
	RtpObserverObserverEvents
	Dominantspeaker *events.Event1[ActiveSpeakerObserverDominantSpeaker]
}
