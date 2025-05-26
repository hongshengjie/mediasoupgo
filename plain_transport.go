package mediasoupgo

import (
	"mediasoupgo/FBS/Notification"
	plaintransport "mediasoupgo/FBS/PlainTransport"

	"mediasoupgo/FBS/Request"
	srtpparameters "mediasoupgo/FBS/SrtpParameters"
	transport "mediasoupgo/FBS/Transport"
	"mediasoupgo/events"
	"strings"
)

var _ PlainTransport = &plainTransportImpl{}

type PlainTransportData struct {
	rtcpMux        bool
	comedia        bool
	tuple          TransportTuple
	rtcpTuple      *TransportTuple
	sctpParameters *SctpParameters
	sctpState      SctpState
	srtpParameters *SrtpParameters
}
type plainTransportImpl struct {
	*transportImpl
	data     *PlainTransportData
	observer events.EventEmmiter[PlainTransportObserverEvents]
	events.EventEmmiter[PlainTransportEvents]
}

func NewPlainTransport(
	options *PlainTransportData,
	id TransportInternal,
	channel *Channel,
	appData TransportAppData,
	getRouterRtpCapabilities func() RtpCapabilities,
	getProducerById func(string) Producer,
	getDataProducerById func(string) DataProducer,
) PlainTransport {
	p := &plainTransportImpl{
		data:         options,
		observer:     events.New[PlainTransportObserverEvents](),
		EventEmmiter: events.New[PlainTransportEvents](),
	}

	ti := NewTransport(id, channel, appData,
		getRouterRtpCapabilities, getProducerById, getDataProducerById,
		func(en events.EventName, te TransportEvents) {
			p.EmitEvent(en, te)
		}, func(en events.EventName, toe TransportObserverEvents) {
			p.observer.Emit(en, PlainTransportObserverEvents{TransportObserverEvents: toe})
		}, "plain")
	p.transportImpl = ti
	p.handleWorkerNotifications()
	p.handleListenerError()
	return p
}

// Transport type
// Override: always returns "plain"
func (p *plainTransportImpl) Type() (_ string) {
	return "plain"
}

// Observer
// Override: returns PlainTransportObserver
func (p *plainTransportImpl) Observer() (_ PlainTransportObserver) {
	return p.observer
}

// PlainTransport tuple
func (p *plainTransportImpl) Tuple() (_ TransportTuple) {
	return p.data.tuple
}

// PlainTransport RTCP tuple
func (p *plainTransportImpl) RTCPTuple() (_ *TransportTuple) {
	return p.data.rtcpTuple
}

// SCTP parameters
func (p *plainTransportImpl) SCTPParameters() (_ *SctpParameters) {
	return p.data.sctpParameters
}

// SCTP state
func (p *plainTransportImpl) SCTPState() (_ SctpState) {
	return p.data.sctpState
}

// SRTP parameters
func (p *plainTransportImpl) SRTPParameters() (_ *SrtpParameters) {
	return p.data.srtpParameters
}

func (p *plainTransportImpl) Close() {
	if p.closed.Load() {
		return
	}
	p.data.sctpState = ClosedSctpState
	p.transportImpl.Close()
}
func (p *plainTransportImpl) RouterClose() {
	if p.closed.Load() {
		return
	}
	p.data.sctpState = ClosedSctpState
	p.transportImpl.Close()
}

// Dump PlainTransport
// Override
func (w *plainTransportImpl) Dump() (_ PlainTransportDump, _ error) {
	_, err := w.channel.Request(Request.MethodTRANSPORT_DUMP, &Request.BodyT{Type: Request.BodyNONE}, w.transportId)
	if err != nil {
		return PlainTransportDump{}, err
	}
	return PlainTransportDump{}, nil
}

// Get PlainTransport stats
// Override
func (w *plainTransportImpl) GetStats() (_ []PlainTransportStat, _ error) {
	_, err := w.channel.Request(Request.MethodTRANSPORT_GET_STATS, nil, w.transportId)
	if err != nil {
		return nil, err
	}
	return []PlainTransportStat{}, nil
}

// Provide the PlainTransport remote parameters
// Override
func (w *plainTransportImpl) Connect(params *PlainTransportConnectParams) (_ error) {
	req := &plaintransport.ConnectRequestT{}
	if params.IP != nil {
		req.Ip = *params.IP
	}

	req.Port = params.Port
	req.RtcpPort = params.RTCPPort
	if params.SRTPParameters != nil {
		req.SrtpParameters = &srtpparameters.SrtpParametersT{
			CryptoSuite: srtpparameters.EnumValuesSrtpCryptoSuite[string(params.SRTPParameters.CryptoSuite)],
			KeyBase64:   params.SRTPParameters.KeyBase64,
		}
	}
	resp, err := w.channel.Request(Request.MethodPLAINTRANSPORT_CONNECT, &Request.BodyT{Type: Request.BodyPlainTransport_ConnectRequest, Value: req}, w.transportId)
	if err != nil {
		return err
	}
	resp2 := resp.Body.Value.(*plaintransport.ConnectResponseT)
	if resp2.Tuple != nil {
		w.data.tuple = TransportTuple{
			LocalIP:      resp2.Tuple.LocalAddress,
			LocalAddress: resp2.Tuple.LocalAddress,
			LocalPort:    (resp2.Tuple.LocalPort),
			RemoteIP:     &resp2.Tuple.RemoteIp,
			RemotePort:   &resp2.Tuple.RemotePort,
			Protocol:     TransportProtocol(strings.ToLower(resp2.Tuple.Protocol.String())),
		}
	}
	if resp2.RtcpTuple != nil {
		w.data.rtcpTuple = &TransportTuple{
			LocalIP:      resp2.RtcpTuple.LocalAddress,
			LocalAddress: resp2.RtcpTuple.LocalAddress,
			LocalPort:    (resp2.RtcpTuple.LocalPort),
			RemoteIP:     &resp2.RtcpTuple.RemoteIp,
			RemotePort:   &resp2.RtcpTuple.RemotePort,
			Protocol:     TransportProtocol(strings.ToLower(resp2.RtcpTuple.Protocol.String())),
		}
	}
	if resp2.SrtpParameters != nil {
		w.data.srtpParameters = &SrtpParameters{
			CryptoSuite: SrtpCryptoSuite(resp2.SrtpParameters.CryptoSuite.String()),
			KeyBase64:   resp2.SrtpParameters.KeyBase64,
		}
	}
	return err

}

func (w *plainTransportImpl) handleWorkerNotifications() {
	w.channel.On(events.EventName(w.transportId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventPLAINTRANSPORT_TUPLE:
			value := arg.Body.Value.(*plaintransport.TupleNotificationT)

			tuple := &TransportTuple{
				LocalIP:      value.Tuple.LocalAddress,
				LocalAddress: value.Tuple.LocalAddress,
				LocalPort:    value.Tuple.LocalPort,
				RemoteIP:     &value.Tuple.RemoteIp,
				RemotePort:   &value.Tuple.RemotePort,
				Protocol:     TransportProtocol(strings.ToLower(value.Tuple.Protocol.String())),
			}
			w.data.tuple = *tuple
			w.Emit("tuple", PlainTransportEvents{Tuple: events.NewEvent1(*tuple)})
			w.observer.Emit("tuple", PlainTransportObserverEvents{Tuple: events.NewEvent1(*tuple)})

		case Notification.EventPLAINTRANSPORT_RTCP_TUPLE:
			value := arg.Body.Value.(*plaintransport.RtcpTupleNotificationT)

			tuple := &TransportTuple{
				LocalIP:      value.Tuple.LocalAddress,
				LocalAddress: value.Tuple.LocalAddress,
				LocalPort:    value.Tuple.LocalPort,
				RemoteIP:     &value.Tuple.RemoteIp,
				RemotePort:   &value.Tuple.RemotePort,
				Protocol:     TransportProtocol(strings.ToLower(value.Tuple.Protocol.String())),
			}
			w.data.rtcpTuple = tuple
			w.Emit("rtcptuple", PlainTransportEvents{RTCPTuple: events.NewEvent1(*tuple)})
			w.observer.Emit("rtcptuple", PlainTransportObserverEvents{RTCPTuple: events.NewEvent1(*tuple)})
		case Notification.EventTRANSPORT_SCTP_STATE_CHANGE:
			value := arg.Body.Value.(*transport.SctpStateChangeNotificationT)
			s := SctpState(strings.ToLower(value.SctpState.String()))
			w.data.sctpState = s
			w.Emit("sctpstatechange", PlainTransportEvents{SctpStateChange: events.NewEvent1(s)})
			w.observer.Emit("sctpstatechange", PlainTransportObserverEvents{SctpStateChange: events.NewEvent1(s)})
		case Notification.EventTRANSPORT_TRACE:
			value := arg.Body.Value.(*transport.TraceNotificationT)

			trace := &TransportTraceEventData{
				Type:      TransportTraceEventType(strings.ToLower(value.Type.String())),
				Timestamp: value.Timestamp,
				Direction: value.Direction.String(),
				Info:      value.Info,
			}

			w.Emit("trace", PlainTransportEvents{TransportEvents: TransportEvents{Trace: events.NewEvent1(*trace)}})
			w.observer.Emit("trace", PlainTransportObserverEvents{TransportObserverEvents: TransportObserverEvents{Trace: events.NewEvent1(*trace)}})
		}
	})

}
func (w *plainTransportImpl) handleListenerError() {}
