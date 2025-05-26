package mediasoupgo

import (
	"mediasoupgo/FBS/Notification"
	pipetransport "mediasoupgo/FBS/PipeTransport"
	"mediasoupgo/FBS/Request"
	srtpparameters "mediasoupgo/FBS/SrtpParameters"
	transport "mediasoupgo/FBS/Transport"
	"mediasoupgo/events"
	"strings"
)

var _ PipeTransport = &pipeTransportImpl{}

type PipeTransportData struct {
	tuple          TransportTuple
	sctpParameters *SctpParameters
	sctpState      SctpState
	rtx            bool
	srtpParameters *SrtpParameters
}

type pipeTransportImpl struct {
	*transportImpl
	data     *PipeTransportData
	observer events.EventEmmiter[PipeTransportObserverEvents]
	events.EventEmmiter[PipeTransportEvents]
}

func NewPipeTransport(
	options *PipeTransportData,
	id TransportInternal,
	channel *Channel,
	appData TransportAppData,
	getRouterRtpCapabilities func() RtpCapabilities,
	getProducerById func(string) Producer,
	getDataProducerById func(string) DataProducer,
) PipeTransport {
	p := &pipeTransportImpl{
		data:         options,
		observer:     events.New[PipeTransportObserverEvents](),
		EventEmmiter: events.New[PipeTransportEvents](),
	}

	ti := NewTransport(id, channel, appData,
		getRouterRtpCapabilities, getProducerById, getDataProducerById,
		func(en events.EventName, te TransportEvents) {
			p.EmitEvent(en, te)
		}, func(en events.EventName, toe TransportObserverEvents) {
			p.observer.Emit(en, PipeTransportObserverEvents{TransportObserverEvents: toe})
		}, "pipe")
	p.transportImpl = ti
	p.handleWorkerNotifications()
	p.handleListenerError()
	return p
}

// PipeTransport tuple
func (p *pipeTransportImpl) Tuple() (_ TransportTuple) {
	return p.data.tuple
}

// SCTP parameters
func (p *pipeTransportImpl) SCTPParameters() (_ *SctpParameters) {
	return p.data.sctpParameters
}

// SCTP state
func (p *pipeTransportImpl) SCTPState() (_ *SctpState) {
	return &p.data.sctpState
}

// SRTP parameters
func (p *pipeTransportImpl) SRTPParameters() (_ *SrtpParameters) {
	return p.data.srtpParameters
}

func (p *pipeTransportImpl) Type() (_ string) {
	return "pipe"
}

// Observer
// Override: returns PipeTransportObserver
func (p *pipeTransportImpl) Observer() (_ PipeTransportObserver) {
	return p.observer
}
func (p *pipeTransportImpl) Close() {
	if p.closed.Load() {
		return
	}
	p.data.sctpState = ClosedSctpState
	p.transportImpl.Close()
}

func (p *pipeTransportImpl) RouterClosed() {
	if p.closed.Load() {
		return
	}
	p.data.sctpState = ClosedSctpState
	p.transportImpl.RouterClosed()
}

// Dump PipeTransport
// Override
func (w *pipeTransportImpl) Dump() (_ PipeTransportDump, _ error) {
	_, err := w.channel.Request(Request.MethodTRANSPORT_DUMP, &Request.BodyT{Type: Request.BodyNONE}, w.transportId)
	if err != nil {
		return PipeTransportDump{}, err
	}
	return PipeTransportDump{}, nil
}

// Get PipeTransport stats
// Override
func (w *pipeTransportImpl) GetStats() (_ []PipeTransportStat, _ error) {
	_, err := w.channel.Request(Request.MethodTRANSPORT_GET_STATS, nil, w.transportId)
	if err != nil {
		return nil, err
	}
	return []PipeTransportStat{}, nil
}

// Provide the PipeTransport remote parameters
// Override
func (w *pipeTransportImpl) Connect(params PipeTransportConnectParams) (_ error) {
	req := &pipetransport.ConnectRequestT{
		Ip:   params.IP,
		Port: &params.Port,
	}
	if params.SRTPParameters != nil {
		req.SrtpParameters = &srtpparameters.SrtpParametersT{
			CryptoSuite: srtpparameters.EnumValuesSrtpCryptoSuite[string(params.SRTPParameters.CryptoSuite)],
			KeyBase64:   params.SRTPParameters.KeyBase64,
		}
	}
	resp, err := w.channel.Request(Request.MethodPIPETRANSPORT_CONNECT, &Request.BodyT{Type: Request.BodyPipeTransport_ConnectRequest, Value: req}, w.transportId)
	if err != nil {
		return err
	}
	resp2 := resp.Body.Value.(*pipetransport.ConnectResponseT)
	if resp2.Tuple != nil {
		w.data.tuple = TransportTuple{
			LocalIP:      resp2.Tuple.LocalAddress,
			LocalAddress: resp2.Tuple.LocalAddress,
			LocalPort:    resp2.Tuple.LocalPort,
			RemoteIP:     &resp2.Tuple.RemoteIp,
			RemotePort:   &resp2.Tuple.RemotePort,
			Protocol:     TransportProtocol(strings.ToLower(resp2.Tuple.Protocol.String())),
		}
	}
	return nil
}

// Create a pipe Consumer
// Override
func (t *pipeTransportImpl) Consume(options *ConsumerOptions) (_ Consumer, _ error) {
	req := &transport.ConsumeRequestT{}
	_, err := t.channel.Request(Request.MethodTRANSPORT_CONSUME, &Request.BodyT{Type: Request.BodyTransport_ConsumeRequest, Value: req}, t.transportId)
	return &consumerImpl{}, err
}

func (p *pipeTransportImpl) handleListenerError() {}

func (p *pipeTransportImpl) handleWorkerNotifications() {
	p.channel.On(events.EventName(p.transportId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventTRANSPORT_SCTP_STATE_CHANGE:

			value := arg.Body.Value.(*transport.SctpStateChangeNotificationT)
			s := SctpState(strings.ToLower(value.SctpState.String()))
			p.data.sctpState = s
			p.Emit("sctpstatechange", PipeTransportEvents{SctpStateChange: events.NewEvent1(s)})
			p.observer.Emit("sctpstatechange", PipeTransportObserverEvents{SctpStateChange: events.NewEvent1(s)})
		case Notification.EventTRANSPORT_TRACE:

			value := arg.Body.Value.(*transport.TraceNotificationT)

			trace := &TransportTraceEventData{
				Type:      TransportTraceEventType(strings.ToLower(value.Type.String())),
				Timestamp: value.Timestamp,
				Direction: value.Direction.String(),
				Info:      value.Info,
			}

			p.Emit("trace", PipeTransportEvents{TransportEvents: TransportEvents{Trace: events.NewEvent1(*trace)}})
			p.observer.Emit("trace", PipeTransportObserverEvents{TransportObserverEvents: TransportObserverEvents{Trace: events.NewEvent1(*trace)}})
		}

	})
}
