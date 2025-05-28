package mediasoupgo

import (
	"log/slog"
	directtransport "mediasoupgo/internal/FBS/DirectTransport"
	"mediasoupgo/internal/FBS/Notification"
	"mediasoupgo/internal/FBS/Request"
	transport "mediasoupgo/internal/FBS/Transport"
	"mediasoupgo/internal/events"
	"strings"
)

var _ DirectTransport = &directTransportImpl{}

type DirectTransportData struct {
	sctpParameters *SctpParameters
}

type directTransportImpl struct {
	*transportImpl
	data     *DirectTransportData
	observer events.EventEmmiter[DirectTransportObserverEvents]
	events.EventEmmiter[DirectTransportEvents]
}

func NewDirectTransport(
	options *DirectTransportData,
	id TransportInternal,
	channel *Channel,
	appData TransportAppData,
	getRouterRtpCapabilities func() RtpCapabilities,
	getProducerById func(string) Producer,
	getDataProducerById func(string) DataProducer,
) DirectTransport {

	observer := events.New[DirectTransportObserverEvents]()
	eventEmmiter := events.New[DirectTransportEvents]()

	fn1 := func(en events.EventName, te TransportEvents) {
		eventEmmiter.Emit(en, DirectTransportEvents{TransportEvents: te})
	}
	fn2 := func(en events.EventName, toe TransportObserverEvents) {
		observer.Emit(en, DirectTransportObserverEvents{TransportObserverEvents: toe})
	}
	d := &directTransportImpl{
		data:         options,
		observer:     observer,
		EventEmmiter: eventEmmiter,
	}
	ti := NewTransport(id, channel, appData,
		getRouterRtpCapabilities, getProducerById, getDataProducerById,
		fn1, fn2, "direct")
	d.transportImpl = ti
	d.handleWorkerNotifications()
	d.handleListenerError()
	return d
}

// Transport type
// Override: always returns "direct"
func (d *directTransportImpl) Type() (_ string) {
	return "direct"
}

func (d *directTransportImpl) Close() {
	if d.closed.Load() {
		return
	}
	d.transportImpl.Close()
}

func (d *directTransportImpl) RouterClosed() {
	if d.closed.Load() {
		return
	}
	d.transportImpl.RouterClosed()
}

// Observer
// Override: returns DirectTransportObserver
func (d *directTransportImpl) Observer() (_ DirectTransportObserver) {
	return d.observer
}

// Dump DirectTransport
// Override
func (w *directTransportImpl) Dump() (_ DirectTransportDump, _ error) {
	_, err := w.channel.Request(Request.MethodTRANSPORT_DUMP, &Request.BodyT{Type: Request.BodyNONE}, w.transportId)
	if err != nil {
		return DirectTransportDump{}, err
	}
	return DirectTransportDump{}, nil
}

// Get DirectTransport stats
// Override
func (w *directTransportImpl) GetStats() (_ []DirectTransportStat, _ error) {
	_, err := w.channel.Request(Request.MethodTRANSPORT_GET_STATS, nil, w.transportId)
	if err != nil {
		return nil, err
	}
	return []DirectTransportStat{}, nil
}

// NO-OP method in DirectTransport
// Override
func (w *directTransportImpl) Connect() (_ error) {
	slog.Debug("connect ")
	return nil
}

// Send RTCP packet
func (w *directTransportImpl) SendRtcp(rtcpPacket []byte) {
	notify := directtransport.RtcpNotificationT{Data: rtcpPacket}
	w.channel.Notify(Notification.EventTRANSPORT_SEND_RTCP, &Notification.BodyT{Type: Notification.BodyTransport_SendRtcpNotification, Value: notify}, w.transportId)
}
func (w *directTransportImpl) handleWorkerNotifications() {
	w.channel.On(events.EventName(w.transportId), func(arg *Notification.NotificationT) {

		switch arg.Event {
		case Notification.EventDIRECTTRANSPORT_RTCP:
			if w.closed.Load() {
				break
			}
			value := arg.Body.Value.(*directtransport.RtcpNotificationT)
			w.Emit("rtcp", DirectTransportEvents{RTCP: events.NewEvent1(value.Data)})

		case Notification.EventTRANSPORT_TRACE:
			value := arg.Body.Value.(*transport.TraceNotificationT)

			trace := &TransportTraceEventData{
				Type:      TransportTraceEventType(strings.ToLower(value.Type.String())),
				Timestamp: value.Timestamp,
				Direction: value.Direction.String(),
				Info:      value.Info,
			}

			w.Emit("trace", DirectTransportEvents{TransportEvents: TransportEvents{Trace: events.NewEvent1(*trace)}})
			w.observer.Emit("trace", DirectTransportObserverEvents{TransportObserverEvents: TransportObserverEvents{Trace: events.NewEvent1(*trace)}})
		}

	})
}
func (w *directTransportImpl) handleListenerError() {}
