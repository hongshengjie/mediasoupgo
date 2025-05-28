package mediasoupgo

import (
	webrtcserver "mediasoupgo/internal/FBS/WebRtcServer"
	"mediasoupgo/internal/events"
)

type WebRtcServerOption struct {
	ListenInfos []*TransportListenInfo
	AppData     WebRtcServerAppData
}

type IpPort struct {
	IP   string
	Port int
}

type IceUserNameFragment struct {
	LocalIceUsernameFragment string
	WebRtcTransportID        string
}

type TupleHash struct {
	TupleHash         int
	WebRtcTransportID string
}

type WebRtcServerDump struct {
	ID                        string
	UDPSockets                []IpPort
	TCPServers                []IpPort
	WebRtcTransportIDs        []string
	LocalIceUsernameFragments []IceUserNameFragment
	TupleHashes               []TupleHash
}

type WebRtcServerEvents struct {
	WorkerClose struct{}
	AtClose     struct{}
}

type WebRtcServerObserverEvents struct {
	Close                    struct{}
	Webrtctransporthandled   *events.Event1[WebRtcTransport]
	Webrtctransportunhandled *events.Event1[WebRtcTransport]
}

type (
	WebRtcServerObserver events.EventEmmiter[WebRtcServerObserverEvents]
	WebRtcServerAppData  = AppData
	WebRtcServer         interface {
		events.EventEmmiter[WebRtcServerEvents]
		Id() string

		Closed() bool

		AppData() WebRtcServerAppData

		SetAppData(appData WebRtcServerAppData)

		Observer() WebRtcServerObserver

		Close()

		WorkerClosed()

		Dump() (*webrtcserver.DumpResponseT, error)

		HandleWebRtcTransport(webRtcTransport WebRtcTransport)
	}
)
