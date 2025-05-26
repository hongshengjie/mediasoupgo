package mediasoupgo

import (
	"errors"
	"mediasoupgo/FBS/Request"
	router "mediasoupgo/FBS/Router"
	rtpobserver "mediasoupgo/FBS/RtpObserver"
	"mediasoupgo/events"
	"sync/atomic"
)

var _ RtpObserver = &rtpObserverImpl{}

type RtpObserverInternal struct {
	RouterInternal
	rtpObserverId string
}

type rtpObserverImpl struct {
	RtpObserverInternal
	channel *Channel
	closed  atomic.Bool
	paused  atomic.Bool

	appData         RtpObserverAppData
	getProducerById func(produceId string) Producer

	emitEvent         func(events.EventName, RtpObserverEvents)
	emitObserverEvent func(events.EventName, RtpObserverObserverEvents)
}

func NewRtpObserver(internal RtpObserverInternal, c *Channel, appData RtpObserverAppData, getProducerById func(producerId string) Producer, emitEvent func(events.EventName, RtpObserverEvents), emitObserverEvent func(events.EventName, RtpObserverObserverEvents)) *rtpObserverImpl {
	r := &rtpObserverImpl{
		RtpObserverInternal: internal,
		channel:             c,
		closed:              atomic.Bool{},
		paused:              atomic.Bool{},
		appData:             appData,
		getProducerById:     getProducerById,
		emitEvent:           emitEvent,
		emitObserverEvent:   emitObserverEvent,
	}
	return r
}
func (r *rtpObserverImpl) ID() (_ string) {
	return r.rtpObserverId
}

func (r *rtpObserverImpl) Closed() (_ bool) {
	return r.closed.Load()
}

func (r *rtpObserverImpl) Paused() (_ bool) {
	return r.paused.Load()
}

func (r *rtpObserverImpl) AppData() (_ AppData) {
	return AppData(r.appData)
}

func (r *rtpObserverImpl) SetAppData(appData AppData) {
	r.appData = RtpObserverAppData(appData)
}

func (r *rtpObserverImpl) Close() {
	if r.closed.Load() {
		return
	}
	r.closed.Store(true)
	r.channel.RemoveAllListeners(events.EventName(r.rtpObserverId))
	r.channel.Request(Request.MethodROUTER_CLOSE_RTPOBSERVER, &Request.BodyT{
		Type:  Request.BodyRouter_CloseRtpObserverRequest,
		Value: &router.CloseRtpObserverRequestT{RtpObserverId: r.rtpObserverId},
	}, r.routerId)
	r.emitEvent("@close", RtpObserverEvents{AtClose: struct{}{}})
	r.emitObserverEvent("close", RtpObserverObserverEvents{Close: struct{}{}})
}

func (r *rtpObserverImpl) RouterClosed() {
	if r.closed.Load() {
		return
	}
	r.closed.Store(true)
	r.channel.RemoveAllListeners(events.EventName(r.rtpObserverId))
	r.emitEvent("routerclose", RtpObserverEvents{Routerclose: struct{}{}})
	r.emitObserverEvent("routerclose", RtpObserverObserverEvents{Close: struct{}{}})
}

func (r *rtpObserverImpl) Pause() (_ error) {
	wasPaused := r.paused.Load()

	r.channel.Request(Request.MethodRTPOBSERVER_PAUSE, nil, r.rtpObserverId)
	r.paused.Store(true)
	if !wasPaused {
		r.emitObserverEvent("pause", RtpObserverObserverEvents{Pause: struct{}{}})
	}
	return nil
}

func (r *rtpObserverImpl) Resume() (_ error) {
	wapPaused := r.paused.Load()

	r.channel.Request(Request.MethodRTPOBSERVER_RESUME, nil, r.rtpObserverId)
	r.paused.Store(false)
	if wapPaused {
		r.emitObserverEvent("resume", RtpObserverObserverEvents{Resume: struct{}{}})
	}
	return nil
}

func (r *rtpObserverImpl) AddProducer(params struct{ ProducerID string }) (_ error) {
	producer := r.getProducerById(params.ProducerID)
	if producer == nil {
		return errors.New(" producer not found")
	}
	r.channel.Request(Request.MethodRTPOBSERVER_ADD_PRODUCER, &Request.BodyT{
		Type:  Request.BodyRtpObserver_AddProducerRequest,
		Value: &rtpobserver.AddProducerRequestT{ProducerId: params.ProducerID},
	}, r.rtpObserverId)
	r.emitObserverEvent("addproducer", RtpObserverObserverEvents{Addproducer: events.NewEvent1(producer)})
	return nil
}

func (r *rtpObserverImpl) RemoveProducer(params struct{ ProducerID string }) (_ error) {
	producer := r.getProducerById(params.ProducerID)
	if producer == nil {
		return errors.New("producer not found")
	}
	r.channel.Request(Request.MethodRTPOBSERVER_REMOVE_PRODUCER, &Request.BodyT{
		Type:  Request.BodyRtpObserver_RemoveProducerRequest,
		Value: &rtpobserver.RemoveProducerRequestT{ProducerId: params.ProducerID},
	}, r.rtpObserverId)
	r.emitObserverEvent("removeproducer", RtpObserverObserverEvents{Removeproducer: events.NewEvent1(producer)})
	return nil
}
