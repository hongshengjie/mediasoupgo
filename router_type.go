package mediasoupgo

import (
	"context"

	"mediasoupgo/internal/events"
)

type RouterOptions struct {
	MediaCodecs []*RtpCodecCapability
	AppData     AppData
}

type (
	PipeToRouterOptions struct {
		ProducerID     *string
		DataProducerID *string
		Router         Router
		EnableSctp     *bool
		NumSctpStreams *NumSctpStreams
		EnableRtx      *bool
		EnableSrtp     *bool
		// Either<PipeToRouterListenInfo, PipeToRouterListenIp>
		ListenInfo *TransportListenInfo // Using pointer to represent optional union
		ListenIp   *interface{}         // Using interface{} to represent TransportListenIp or string
	}
)

type PipeToRouterListenInfo struct {
	ListenInfo TransportListenInfo
}

type PipeToRouterListenIp struct {
	ListenIp interface{} // Can be TransportListenIp or string
}

type PipeToRouterResult struct {
	PipeConsumer     *Consumer
	PipeProducer     *Producer
	PipeDataConsumer *DataConsumer
	PipeDataProducer *DataProducer
}

type PipeTransportPair map[string]PipeTransport

type RouterDump struct {
	ID                       string
	TransportIDs             []string
	RTPObserverIDs           []string
	MapProducerIDConsumerIDs []struct {
		Key    string
		Values []string
	}
	MapConsumerIDProducerID []struct {
		Key   string
		Value string
	}
	MapProducerIDObserverIDs []struct {
		Key    string
		Values []string
	}
	MapDataProducerIDDataConsumerIDs []struct {
		Key    string
		Values []string
	}
	MapDataConsumerIDDataProducerID []struct {
		Key   string
		Value string
	}
}

type (
	RouterAppData  AppData
	RouterObserver events.EventEmmiter[RouterObserverEvents]
	RouterOption   struct {
		MediaCodecs []*RtpCodecCapability
		AppData     RouterAppData
	}
)

type RouterObserverEvents struct {
	Close          struct{}
	Newtransport   *events.Event1[Transport]
	Newrtpobserver *events.Event1[RtpObserver]
}

type RouterEvents struct {
	WorkerClose struct{}
	AtClose     struct{}
}

type Router interface {
	events.EventEmmiter[RouterEvents]
	Id() string
	Closed() bool
	RtpCapabilities() *RtpCapabilities
	AppData() RouterAppData
	SetAppData(appData RouterAppData)
	Observer() RouterObserver
	Close()
	WorkerClosed()
	Dump(ctx context.Context) (*RouterDump, error)
	CreateWebRtcTransport(options *WebRtcTransportOptions) (WebRtcTransport, error)
	CreatePlainTransport(options *PlainTransportOptions) (PlainTransport, error)
	CreatePipeTransport(options *PipeTransportOptions) (PipeTransport, error)
	CreateDirectTransport(options *DirectTransportOptions) (DirectTransport, error)
	PipeToRouter(options *PipeToRouterOptions) (PipeToRouterResult, error)
	AddPipeTransportPair(pipeTransportPairKey string, pipeTransportPairPromise PipeTransportPair)
	CreateActiveSpeakerObserver(options *ActiveSpeakerObserverOptions) (ActiveSpeakerObserver, error)
	CreateAudioLevelObserver(options *AudioLevelObserverOptions) (AudioLevelObserver, error)
	CanConsume(producerId string, rtpCapabilities *RtpCapabilities) bool
}
