package mediasoupgo

import (
	"log/slog"
	"strings"

	"mediasoupgo/FBS/Notification"
	"mediasoupgo/FBS/Request"
	transport "mediasoupgo/FBS/Transport"
	webrtctransport "mediasoupgo/FBS/WebRtcTransport"
	"mediasoupgo/events"
	"mediasoupgo/ptr"
)

var _ WebRtcTransport = &webRtcTransportImpl{}

type WebRtcTransportData struct {
	iceRole          IceRole
	iceParameters    IceParameters
	iceCandidates    []*IceCandidate
	iceState         IceState
	iceSelectedTuple *TransportTuple
	dtlsParameters   DtlsParameters
	dtlsState        DtlsState
	dtlsRemoteCert   *string
	sctpParameters   *SctpParameters
	sctpState        SctpState
}

type webRtcTransportImpl struct {
	*transportImpl
	data *WebRtcTransportData
	events.EventEmmiter[WebRtcTransportEvents]
	observer events.EventEmmiter[WebRtcTransportObserverEvents]
}

func NewWebRtcTransport(
	options *WebRtcTransportData,
	id TransportInternal,
	channel *Channel,
	appData TransportAppData,
	getRouterRtpCapabilities func() RtpCapabilities,
	getProducerById func(string) Producer,
	getDataProducerById func(string) DataProducer,
) WebRtcTransport {
	observer := events.New[WebRtcTransportObserverEvents]()
	eventEmmiter := events.New[WebRtcTransportEvents]()

	fn1 := func(en events.EventName, te TransportEvents) {
		eventEmmiter.Emit(en, WebRtcTransportEvents{TransportEvents: te})
	}
	fn2 := func(en events.EventName, toe TransportObserverEvents) {
		observer.Emit(en, WebRtcTransportObserverEvents{TransportObserverEvents: toe})
	}
	wi := &webRtcTransportImpl{
		data:         options,
		observer:     observer,
		EventEmmiter: eventEmmiter,
	}
	ti := NewTransport(
		id,
		channel,
		appData,
		getRouterRtpCapabilities,
		getProducerById,
		getDataProducerById,
		fn1,
		fn2,
		"webrtc")
	wi.transportImpl = ti
	wi.handleWorkerNotifications()
	wi.handleListenerError()
	return wi
}

// Transport type
// Override: always returns "webrtc"
func (w *webRtcTransportImpl) Type() string {
	return "webrtc"
}

// Observer
// Override: returns WebRtcTransportObserver
func (w *webRtcTransportImpl) Observer() WebRtcTransportObserver {
	return w.observer
}

// ICE role
// Always returns "controlled"
func (w *webRtcTransportImpl) IceRole() string {
	return string(w.data.iceRole)
}

// ICE parameters
func (w *webRtcTransportImpl) IceParameters() IceParameters {
	return w.data.iceParameters
}

// ICE candidates
func (w *webRtcTransportImpl) IceCandidates() []*IceCandidate {
	return w.data.iceCandidates
}

// ICE state
func (w *webRtcTransportImpl) IceState() IceState {
	return w.data.iceState
}

// ICE selected tuple
func (w *webRtcTransportImpl) IceSelectedTuple() *TransportTuple {
	return w.data.iceSelectedTuple
}

// DTLS parameters
func (w *webRtcTransportImpl) DtlsParameters() DtlsParameters {
	return w.data.dtlsParameters
}

// DTLS state
func (w *webRtcTransportImpl) DtlsState() DtlsState {
	return w.data.dtlsState
}

// Remote certificate in PEM format
func (w *webRtcTransportImpl) DtlsRemoteCert() *string {
	return w.data.dtlsRemoteCert
}

// SCTP parameters
func (w *webRtcTransportImpl) SctpParameters() *SctpParameters {
	return w.data.sctpParameters
}

// SCTP state
func (w *webRtcTransportImpl) SctpState() *SctpState {
	return &w.data.sctpState
}

func (w *webRtcTransportImpl) Close() {
	if w.closed.Load() {
		return
	}
	w.data.iceState = "closed"
	w.data.iceSelectedTuple = nil
	w.data.dtlsState = "closed"
	w.data.sctpState = "closed"
	w.transportImpl.Close()
}

func (w *webRtcTransportImpl) RouterClosed() {
	if w.closed.Load() {
		return
	}
	w.data.iceState = "closed"
	w.data.iceSelectedTuple = nil
	w.data.dtlsState = "closed"
	w.data.sctpState = "closed"
	w.transportImpl.RouterClosed()
}

func (w *webRtcTransportImpl) ListenServerClosed() {
	if w.closed.Load() {
		return
	}
	w.data.iceState = "closed"
	w.data.iceSelectedTuple = nil
	w.data.dtlsState = "closed"
	w.data.sctpState = "closed"
	w.transportImpl.ListenServerClosed()
}

// Dump WebRtcTransport
// Override
func (w *webRtcTransportImpl) Dump() (WebRtcTransportDump, error) {
	_, err := w.channel.Request(
		Request.MethodTRANSPORT_DUMP,
		&Request.BodyT{Type: Request.BodyNONE},
		w.transportId,
	)
	if err != nil {
		return WebRtcTransportDump{}, err
	}
	return WebRtcTransportDump{}, nil
}

// Get WebRtcTransport stats
// Override
func (w *webRtcTransportImpl) GetStats() (_ []WebRtcTransportStat, _ error) {
	_, err := w.channel.Request(Request.MethodTRANSPORT_GET_STATS, nil, w.transportId)
	if err != nil {
		return nil, err
	}
	return []WebRtcTransportStat{}, nil
}

// Provide the WebRtcTransport remote parameters
// Override
func (w *webRtcTransportImpl) Connect(dtlsParameters DtlsParameters) (_ error) {
	var fgp []*webrtctransport.FingerprintT
	for _, v := range dtlsParameters.Fingerprints {
		fgp = append(fgp, ToFbsDtlsFingerprint(v))
	}
	var role webrtctransport.DtlsRole
	if dtlsParameters.Role != nil {
		role = webrtctransport.EnumValuesDtlsRole[strings.ToUpper(string(*dtlsParameters.Role))]
	}
	req := &webrtctransport.ConnectRequestT{
		DtlsParameters: &webrtctransport.DtlsParametersT{
			Fingerprints: fgp,
			Role:         role,
		},
	}
	response, err := w.channel.Request(
		Request.MethodWEBRTCTRANSPORT_CONNECT,
		&Request.BodyT{Type: Request.BodyWebRtcTransport_ConnectRequest, Value: req},
		w.transportId,
	)
	if err != nil {
		return err
	}
	resp2 := response.Body.Value.(*webrtctransport.ConnectResponseT)
	w.data.dtlsParameters.Role = ptr.To(
		DtlsRole(strings.ToLower(resp2.DtlsLocalRole.String())),
	)
	return nil
}

// Restart ICE
func (w *webRtcTransportImpl) RestartIce() (_ IceParameters, _ error) {
	resp, err := w.channel.Request(
		Request.MethodTRANSPORT_RESTART_ICE,
		&Request.BodyT{Type: Request.BodyNONE, Value: nil},
		w.transportId,
	)
	if err != nil {
		return IceParameters{}, err
	}
	resp2 := resp.Body.Value.(*transport.RestartIceResponseT)
	return IceParameters{
		UsernameFragment: resp2.UsernameFragment,
		Password:         resp2.Password,
		IceLite:          &resp2.IceLite,
	}, nil
}

func (w *webRtcTransportImpl) handleWorkerNotifications() {
	w.channel.On(events.EventName(w.transportId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventWEBRTCTRANSPORT_ICE_STATE_CHANGE:
			value := arg.Body.Value.(*webrtctransport.IceStateChangeNotificationT)
			iceState := IceState(strings.ToLower(value.IceState.String()))

			w.data.iceState = IceState(iceState)
			w.Emit(
				"icestatechange",
				WebRtcTransportEvents{IceStateChange: events.NewEvent1(iceState)},
			)
			w.observer.Emit(
				"icestatechange",
				WebRtcTransportObserverEvents{IceStateChange: events.NewEvent1(iceState)},
			)
		case Notification.EventWEBRTCTRANSPORT_ICE_SELECTED_TUPLE_CHANGE:
			value := arg.Body.Value.(*webrtctransport.IceSelectedTupleChangeNotificationT)

			iceSelectedTuple := &TransportTuple{
				LocalIP:      value.Tuple.LocalAddress,
				LocalAddress: value.Tuple.LocalAddress,
				LocalPort:    value.Tuple.LocalPort,
				RemoteIP:     &value.Tuple.RemoteIp,
				RemotePort:   &value.Tuple.RemotePort,
				Protocol:     TransportProtocol(strings.ToLower(value.Tuple.Protocol.String())),
			}
			w.data.iceSelectedTuple = iceSelectedTuple
			w.Emit(
				"iceselectedtuplechange",
				WebRtcTransportEvents{IceSelectedTupleChange: events.NewEvent1(*iceSelectedTuple)},
			)
			w.observer.Emit(
				"iceselectedtuplechange",
				WebRtcTransportObserverEvents{
					IceSelectedTupleChange: events.NewEvent1(*iceSelectedTuple),
				},
			)

		case Notification.EventWEBRTCTRANSPORT_DTLS_STATE_CHANGE:
			value := arg.Body.Value.(*webrtctransport.DtlsStateChangeNotificationT)
			w.data.dtlsState = DtlsState(value.DtlsState.String())
			if w.data.dtlsState == ConnectedDtlsState {
				w.data.dtlsRemoteCert = &value.RemoteCert
			}
			w.Emit(
				"dtlsstatechange",
				WebRtcTransportEvents{DtlsStateChange: events.NewEvent1(w.data.dtlsState)},
			)
			w.observer.Emit(
				"dtlsstatechange",
				WebRtcTransportObserverEvents{DtlsStateChange: events.NewEvent1(w.data.dtlsState)},
			)
		case Notification.EventTRANSPORT_SCTP_STATE_CHANGE:
			value := arg.Body.Value.(*transport.SctpStateChangeNotificationT)
			w.data.sctpState = SctpState(strings.ToLower(value.SctpState.String()))
			w.Emit(
				"sctpstatechange",
				WebRtcTransportEvents{SctpStateChange: events.NewEvent1(w.data.sctpState)},
			)
			w.observer.Emit(
				"sctpstatechange",
				WebRtcTransportObserverEvents{SctpStateChange: events.NewEvent1(w.data.sctpState)},
			)
		case Notification.EventTRANSPORT_TRACE:

			value := arg.Body.Value.(*transport.TraceNotificationT)
			trace := &TransportTraceEventData{
				Type:      TransportTraceEventType(strings.ToLower(value.Type.String())),
				Timestamp: value.Timestamp,
				Direction: value.Direction.String(),
				Info:      value.Info,
			}

			w.Emit(
				"trace",
				WebRtcTransportEvents{
					TransportEvents: TransportEvents{Trace: events.NewEvent1(*trace)},
				},
			)
			w.observer.Emit(
				"trace",
				WebRtcTransportObserverEvents{
					TransportObserverEvents: TransportObserverEvents{
						Trace: events.NewEvent1(*trace),
					},
				},
			)
		default:
			slog.Error("ignoring unknown event ", slog.Any("event", arg))
		}
	})
}

func (w *webRtcTransportImpl) handleListenerError() {
	w.On("listenererror", func(arg WebRtcTransportEvents) {
		// TODO
	})
}
