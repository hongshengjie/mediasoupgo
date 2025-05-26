package mediasoupgo

import (
	fbsaudiolevelobserver "mediasoupgo/FBS/AudioLevelObserver"
	"mediasoupgo/FBS/Notification"
	"mediasoupgo/events"
)

var _ AudioLevelObserver = &audioLevelObserverImpl{}

type audioLevelObserverImpl struct {
	*rtpObserverImpl
	observer events.EventEmmiter[AudioLevelObserverObserverEvents]
	events.EventEmmiter[AudioLevelObserverEvents]
}

func NewAudioLevelObserver(internal RtpObserverInternal, c *Channel, appData RtpObserverAppData, getProducerById func(producerId string) Producer) AudioLevelObserver {
	a := &audioLevelObserverImpl{
		observer:     events.New[AudioLevelObserverObserverEvents](),
		EventEmmiter: events.New[AudioLevelObserverEvents](),
	}
	aa := NewRtpObserver(internal, c, appData, getProducerById,
		func(en events.EventName, roe RtpObserverEvents) {
			a.Emit(en, AudioLevelObserverEvents{RtpObserverEvents: roe})
		},
		func(en events.EventName, rooe RtpObserverObserverEvents) {
			a.observer.Emit(en, AudioLevelObserverObserverEvents{RtpObserverObserverEvents: rooe})
		})
	a.rtpObserverImpl = aa
	a.handleWorkerNotifications()
	a.handleListenerError()
	return a
}

// RtpObserver type
// Override: always returns "audiolevel"
func (a *audioLevelObserverImpl) Type() string {
	return "audiolevel"
}

// Observer
// Override: returns AudioLevelObserverObserver
func (a *audioLevelObserverImpl) Observer() AudioLevelObserverObserver {
	return a.observer
}

func (a *audioLevelObserverImpl) handleWorkerNotifications() {
	a.channel.On(events.EventName(a.rtpObserverId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventAUDIOLEVELOBSERVER_VOLUMES:
			value := arg.Body.Value.(*fbsaudiolevelobserver.VolumesNotificationT)

			var volumes []*AudioLevelObserverVolume
			for _, v := range value.Volumes {
				producer := a.getProducerById(v.ProducerId)
				volumes = append(volumes, &AudioLevelObserverVolume{
					Producer: producer,
					Volume:   int(v.Volume),
				})
			}
			if len(volumes) > 0 {
				a.Emit("volumes", AudioLevelObserverEvents{Volumes: events.NewEvent1(volumes)})
				a.observer.Emit("volumes", AudioLevelObserverObserverEvents{Volumes: events.NewEvent1(volumes)})
			}
		case Notification.EventAUDIOLEVELOBSERVER_SILENCE:
			a.Emit("silence", AudioLevelObserverEvents{Silence: struct{}{}})
			a.observer.Emit("silence", AudioLevelObserverObserverEvents{Silence: struct{}{}})
		}
	})
}
func (a *audioLevelObserverImpl) handleListenerError() {}
