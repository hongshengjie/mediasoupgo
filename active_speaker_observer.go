package mediasoupgo

import (
	fbsactivespeakerobserver "mediasoupgo/FBS/ActiveSpeakerObserver"
	"mediasoupgo/FBS/Notification"
	"mediasoupgo/events"
)

var _ ActiveSpeakerObserver = &activeSpeakerObserverImpl{}

type activeSpeakerObserverImpl struct {
	*rtpObserverImpl
	observer events.EventEmmiter[ActiveSpeakerObserverObserverEvents]
	events.EventEmmiter[ActiveSpeakerObserverEvents]
}

func NewActiveSpeakerObserver(internal RtpObserverInternal, c *Channel, appData RtpObserverAppData, getProducerById func(producerId string) Producer) ActiveSpeakerObserver {
	a := &activeSpeakerObserverImpl{
		observer:     events.New[ActiveSpeakerObserverObserverEvents](),
		EventEmmiter: events.New[ActiveSpeakerObserverEvents](),
	}
	aa := NewRtpObserver(internal, c, appData, getProducerById,
		func(en events.EventName, roe RtpObserverEvents) {
			a.Emit(en, ActiveSpeakerObserverEvents{RtpObserverEvents: roe})
		},
		func(en events.EventName, rooe RtpObserverObserverEvents) {
			a.observer.Emit(en, ActiveSpeakerObserverObserverEvents{RtpObserverObserverEvents: rooe})
		})
	a.rtpObserverImpl = aa
	a.handleWorkerNotifications()
	a.handleListenerError()
	return a
}

// RtpObserver type
// Override: always returns "activespeaker"
func (a *activeSpeakerObserverImpl) Type() string {
	return "activespeaker"
}

// Observer
// Override: returns ActiveSpeakerObserverObserver
func (a *activeSpeakerObserverImpl) Observer() ActiveSpeakerObserverObserver {
	return a.observer

}
func (a *activeSpeakerObserverImpl) handleWorkerNotifications() {
	a.channel.On(events.EventName(a.rtpObserverId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventACTIVESPEAKEROBSERVER_DOMINANT_SPEAKER:
			value := arg.Body.Value.(fbsactivespeakerobserver.DominantSpeakerNotificationT)
			producer := a.getProducerById(value.ProducerId)
			if producer == nil {
				break
			}
			a.Emit("dominantspeaker", ActiveSpeakerObserverEvents{Dominantspeaker: events.NewEvent1(ActiveSpeakerObserverDominantSpeaker{Producer: producer})})

			a.observer.Emit("dominantspeaker", ActiveSpeakerObserverObserverEvents{
				Dominantspeaker: events.NewEvent1(ActiveSpeakerObserverDominantSpeaker{
					Producer: producer,
				}),
			})
		}
	})
}

func (a *activeSpeakerObserverImpl) handleListenerError() {}
