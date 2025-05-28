package mediasoupgo

import (
	"mediasoupgo/internal/events"
)

type WebRtcTransportOptions struct {
	EnableUdp                       *bool
	EnableTcp                       *bool
	PreferUdp                       *bool
	PreferTcp                       *bool
	IceConsentTimeout               *int
	InitialAvailableOutgoingBitrate *int
	EnableSctp                      *bool
	NumSctpStreams                  *NumSctpStreams
	MaxSctpMessageSize              *int
	SctpSendBufferSize              *int
	AppData                         AppData
	WebRtcTransportListen           *WebRtcTransportListen
}

// WebRtcTransportListen is an interface that can be one of the three types
type WebRtcTransportListen struct {
	ListenInfos  []TransportListenInfo `json:"listenInfos"`
	WebRtcServer WebRtcServer          `json:"webRtcServer"`
}

type IceParameters struct {
	UsernameFragment string `json:"usernameFragment"`
	Password         string `json:"password"`
	IceLite          *bool  `json:"iceLite"`
}

type IceCandidate struct {
	Foundation string               `json:"foundation"`
	Priority   uint32               `json:"priority"`
	IP         string               `json:"ip"`
	Address    string               `json:"address"`
	Protocol   TransportProtocol    `json:"protocol"`
	Port       uint16               `json:"port"`
	Type       IceCandidateType     `json:"type"`
	TCPType    *IceCandidateTcpType `json:"tcpType"`
}

type DtlsParameters struct {
	Role         *DtlsRole         `json:"role"`
	Fingerprints []DtlsFingerprint `json:"fingerprints"`
}

type FingerprintAlgorithm string

const (
	SHA1FingerprintAlgorithm   FingerprintAlgorithm = "sha-1"
	SHA224FingerprintAlgorithm FingerprintAlgorithm = "sha-224"
	SHA256FingerprintAlgorithm FingerprintAlgorithm = "sha-256"
	SHA384FingerprintAlgorithm FingerprintAlgorithm = "sha-384"
	SHA512FingerprintAlgorithm FingerprintAlgorithm = "sha-512"
)

type DtlsFingerprint struct {
	Algorithm FingerprintAlgorithm `json:"algorithm"`
	Value     string               `json:"value"`
}

type IceRole string

const (
	ControlledIceRole  IceRole = "controlled"
	ControllingIceRole IceRole = "controlling"
)

type IceState string

const (
	NewIceState          IceState = "new"
	ConnectedIceState    IceState = "connected"
	CompletedIceState    IceState = "completed"
	DisconnectedIceState IceState = "disconnected"
	ClosedIceState       IceState = "closed"
)

type IceCandidateType string

const (
	HostIceCandidateType IceCandidateType = "host"
)

type IceCandidateTcpType string

const (
	PassiveIceCandidateTcpType IceCandidateTcpType = "passive"
)

type DtlsRole string

const (
	AutoDtlsRole   DtlsRole = "auto"
	ClientDtlsRole DtlsRole = "client"
	ServerDtlsRole DtlsRole = "server"
)

type DtlsState string

const (
	NewDtlsState        DtlsState = "new"
	ConnectingDtlsState DtlsState = "connecting"
	ConnectedDtlsState  DtlsState = "connected"
	FailedDtlsState     DtlsState = "failed"
	ClosedDtlsState     DtlsState = "closed"
)

type WebRtcTransportDump struct {
	BaseTransportDump
	IceRole          IceRole
	IceParameters    IceParameters
	IceCandidates    []IceCandidate
	IceState         IceState
	IceSelectedTuple *TransportTuple
	DtlsParameters   DtlsParameters
	DtlsState        DtlsState
	DtlsRemoteCert   *string
}

type WebRtcTransportStat struct {
	BaseTransportStats
	Type             string
	IceRole          string
	IceState         IceState
	IceSelectedTuple *TransportTuple
	DtlsState        DtlsState
}

type WebRtcTransportEvents struct {
	TransportEvents
	IceStateChange         *events.Event1[IceState]
	IceSelectedTupleChange *events.Event1[TransportTuple]
	DtlsStateChange        *events.Event1[DtlsState]
	SctpStateChange        *events.Event1[SctpState]
}

type (
	WebRtcTransportObserver       events.EventEmmiter[WebRtcTransportObserverEvents]
	WebRtcTransportObserverEvents struct {
		TransportObserverEvents
		IceStateChange         *events.Event1[IceState]
		IceSelectedTupleChange *events.Event1[TransportTuple]
		DtlsStateChange        *events.Event1[DtlsState]
		SctpStateChange        *events.Event1[SctpState]
	}
)

type (
	WebRtcTransport interface {
		events.EventEmmiter[WebRtcTransportEvents]
		Transport
		// Transport type
		// Override: always returns "webrtc"
		Type() string

		// Observer
		// Override: returns WebRtcTransportObserver
		Observer() WebRtcTransportObserver

		// ICE role
		// Always returns "controlled"
		IceRole() string

		// ICE parameters
		IceParameters() IceParameters

		// ICE candidates
		IceCandidates() []*IceCandidate

		// ICE state
		IceState() IceState

		// ICE selected tuple
		IceSelectedTuple() *TransportTuple // Using pointer for undefined possibility

		// DTLS parameters
		DtlsParameters() DtlsParameters

		// DTLS state
		DtlsState() DtlsState

		// Remote certificate in PEM format
		DtlsRemoteCert() *string // Using pointer for undefined possibility

		// SCTP parameters
		SctpParameters() *SctpParameters // Using pointer for undefined possibility

		// SCTP state
		SctpState() *SctpState // Using pointer for undefined possibility

		// Dump WebRtcTransport
		// Override
		Dump() (WebRtcTransportDump, error)

		// Get WebRtcTransport stats
		// Override
		GetStats() ([]WebRtcTransportStat, error)

		// Provide the WebRtcTransport remote parameters
		// Override
		Connect(dtlsParameters DtlsParameters) error

		// Restart ICE
		RestartIce() (IceParameters, error)
	}
)
