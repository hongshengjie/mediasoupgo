package mediasoupgo

import (
	"mediasoupgo/FBS/Request"
	webrtcserver "mediasoupgo/FBS/WebRtcServer"
	worker "mediasoupgo/FBS/Worker"
	"mediasoupgo/events"
	"mediasoupgo/smap"
	"sync/atomic"
)

var _ WebRtcServer = &webRtcServerImpl{}

type webRtcServerImpl struct {
	webRtcServerId   string
	channel          *Channel
	closed           atomic.Bool
	appData          WebRtcServerAppData
	webRtcTransports *smap.Map[string, WebRtcTransport]
	observer         WebRtcServerObserver
	events.EventEmmiter[WebRtcServerEvents]
}

func NewWebRtcServer(id string, chnnel *Channel, appData WebRtcServerAppData) WebRtcServer {
	s := &webRtcServerImpl{
		webRtcServerId:   id,
		channel:          chnnel,
		appData:          appData,
		webRtcTransports: smap.New[string, WebRtcTransport](),
		observer:         events.New[WebRtcServerObserverEvents](),
		EventEmmiter:     events.New[WebRtcServerEvents](),
	}
	s.handleListernerError()
	return s
}

func (w *webRtcServerImpl) Id() string {
	return w.webRtcServerId
}

func (w *webRtcServerImpl) Closed() bool {
	return w.closed.Load()
}

func (w *webRtcServerImpl) AppData() WebRtcServerAppData {
	return w.appData
}

func (w *webRtcServerImpl) SetAppData(appData WebRtcServerAppData) {
	w.appData = appData
}

func (w *webRtcServerImpl) Observer() WebRtcServerObserver {
	return w.observer
}

func (w *webRtcServerImpl) Close() {
	if w.closed.Load() {
		return
	}
	w.closed.Store(true)

	w.channel.Request(
		Request.MethodWORKER_WEBRTCSERVER_CLOSE,
		&Request.BodyT{
			Type:  Request.BodyWorker_CloseWebRtcServerRequest,
			Value: &worker.CloseWebRtcServerRequestT{WebRtcServerId: w.Id()},
		},
		w.Id(),
	)
	w.webRtcTransports.Range(func(key string, value WebRtcTransport) bool {

		value.ListenServerClosed()
		w.observer.Emit(
			"webrtctransportunhandled",
			WebRtcServerObserverEvents{Webrtctransportunhandled: events.NewEvent1(value)},
		)
		return true
	})
	w.webRtcTransports = smap.New[string, WebRtcTransport]()
	w.Emit("@close", WebRtcServerEvents{AtClose: struct{}{}})
	w.observer.Emit("close", WebRtcServerObserverEvents{Close: struct{}{}})
}

func (w *webRtcServerImpl) WorkerClosed() {
	if w.closed.Load() {
		return
	}

	w.closed.Store(true)
	w.webRtcTransports = smap.New[string, WebRtcTransport]()
	w.Emit("workerclose", WebRtcServerEvents{WorkerClose: struct{}{}})
	w.observer.Emit("close", WebRtcServerObserverEvents{
		Close: struct{}{},
	})
}

func (w *webRtcServerImpl) HandleWebRtcTransport(webRtcTransport WebRtcTransport) {
	w.webRtcTransports.Set(webRtcTransport.ID(), webRtcTransport)
	w.observer.Emit("webrtctransporthandled", WebRtcServerObserverEvents{
		Webrtctransporthandled: events.NewEvent1(webRtcTransport),
	})
	webRtcTransport.On("@close", func(arg WebRtcTransportEvents) {
		w.webRtcTransports.Delete(webRtcTransport.ID())
		w.observer.Emit(
			"webrtctransportunhandled",
			WebRtcServerObserverEvents{Webrtctransportunhandled: events.NewEvent1(webRtcTransport)},
		)
	})
}

func (w *webRtcServerImpl) Dump() (*webrtcserver.DumpResponseT, error) {
	resp, err := w.channel.Request(Request.MethodWEBRTCSERVER_DUMP, nil, w.webRtcServerId)
	if err != nil {
		return nil, err
	}
	return resp.Body.Value.(*webrtcserver.DumpResponseT), nil
}

func (w *webRtcServerImpl) handleListernerError() {
	w.On("listenererror", func(arg WebRtcServerEvents) {
		// TODO
	})
}
