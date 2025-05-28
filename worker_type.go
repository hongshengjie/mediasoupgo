package mediasoupgo

import (
	"mediasoupgo/internal/events"
)

type WorkerLogLevel string

const (
	Debug WorkerLogLevel = "debug"
	Warn  WorkerLogLevel = "warn"
	Error WorkerLogLevel = "error"
	None  WorkerLogLevel = "none"
)

type WorkerLogTag string

const (
	Info      WorkerLogTag = "info"
	Ice       WorkerLogTag = "ice"
	Dtls      WorkerLogTag = "dtls"
	Rtp       WorkerLogTag = "rtp"
	Srtp      WorkerLogTag = "srtp"
	Rtcp      WorkerLogTag = "rtcp"
	Rtx       WorkerLogTag = "rtx"
	Bwe       WorkerLogTag = "bwe"
	Score     WorkerLogTag = "score"
	Simulcast WorkerLogTag = "simulcast"
	Svc       WorkerLogTag = "svc"
	Sctp      WorkerLogTag = "sctp"
	Message   WorkerLogTag = "message"
)

func DefaultLogTags() []WorkerLogTag {
	return []WorkerLogTag{
		Info,
		Ice,
		Dtls,
		Rtp,
		Srtp,
		Rtcp,
		Rtx,
		Bwe,
		Score,
		Simulcast,
		Svc,
		Sctp,
		Message,
	}
}

type WorkerSettings struct {
	LogLevel             *WorkerLogLevel
	LogTags              []WorkerLogTag
	RTCMinPort           *int
	RTCMaxPort           *int
	DTLSCertificateFile  *string
	DTLSPrivateKeyFile   *string
	LibwebrtcFieldTrials *string
	DisableLiburing      *bool
	AppData              WorkerAppData
}

type WorkerResourceUsage struct {
	RUUtime    uint64
	RUStime    uint64
	RUMaxrss   uint64
	RUIxrss    uint64
	RUIdrss    uint64
	RUIsrss    uint64
	RUMinflt   uint64
	RUMajflt   uint64
	RUNswap    uint64
	RUInblock  uint64
	RUOublock  uint64
	RUMsgsnd   uint64
	RUMsgrcv   uint64
	RUNsignals uint64
	RUNvcsw    uint64
	RUNivcsw   uint64
}

type WorkerDump struct {
	PID                    uint32
	WebRTCServerIDs        []string
	RouterIDs              []string
	ChannelMessageHandlers struct {
		ChannelRequestHandlers      []string
		ChannelNotificationHandlers []string
	}
	Liburing *struct {
		SQEProcessCount   uint64
		SQEMissCount      uint64
		UserDataMissCount uint64
	}
}

type (
	WorkerObserver = events.EventEmmiter[WorkerObserverEvents]
	WorkerAppData  = AppData
)

type WorkerEvents struct {
	Died            *events.Event1[error]
	Subprocessclose struct{}
	AtSuccess       struct{}
	AtFailure       *events.Event1[error]
}

type WorkerObserverEvents struct {
	Close           struct{}
	Newwebrtcserver *events.Event1[WebRtcServer]
	Newrouter       *events.Event1[Router]
}

type WorkerUpdateableSettings struct {
	LogLevel string
	LogTags  []string
}

type Worker interface {
	Pid() int
	Closed() bool

	Died() bool

	SubprocessClosed() bool

	AppData() WorkerAppData

	SetAppData(appData WorkerAppData)

	Observer() WorkerObserver

	Close()

	Dump() (*WorkerDump, error)

	GetResourceUsage() (*WorkerResourceUsage, error)

	UpdateSettings(setting *WorkerUpdateableSettings) error

	CreateWebRtcServer(options *WebRtcServerOption) (WebRtcServer, error)

	CreateRouter(options *RouterOption) (Router, error)
}
