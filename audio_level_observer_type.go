package mediasoupgo

import "mediasoupgo/events"

// AudioLevelObserver interface definition
type AudioLevelObserver interface {
	events.EventEmmiter[AudioLevelObserverEvents]
	RtpObserver
	Type() string

	Observer() AudioLevelObserverObserver
}

type AudioLevelObserverOptions struct {
	MaxEntries uint16
	Threshold  uint8
	Interval   uint16
	AppData    AudioLevelObserverAppData
}

type (
	AudioLevelObserverAppData = AppData
	AudioLevelObserverVolume  struct {
		Producer Producer
		Volume   int
	}
)

type AudioLevelObserverEvents struct {
	RtpObserverEvents
	Volumes *events.Event1[[]*AudioLevelObserverVolume]
	Silence struct{}
}

type AudioLevelObserverObserver events.EventEmmiter[AudioLevelObserverObserverEvents]

type AudioLevelObserverObserverEvents struct {
	RtpObserverObserverEvents
	Volumes *events.Event1[[]*AudioLevelObserverVolume]
	Silence struct{}
}
