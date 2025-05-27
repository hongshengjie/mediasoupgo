package mediasoupgo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"

	fbsactivespeakerobserver "mediasoupgo/FBS/ActiveSpeakerObserver"
	fbsaudiolevelobserver "mediasoupgo/FBS/AudioLevelObserver"
	fbsdirecttransport "mediasoupgo/FBS/DirectTransport"
	fbspipetransport "mediasoupgo/FBS/PipeTransport"
	fbsplaintransport "mediasoupgo/FBS/PlainTransport"
	"mediasoupgo/FBS/Request"
	router "mediasoupgo/FBS/Router"
	sctpparameters "mediasoupgo/FBS/SctpParameters"
	transport "mediasoupgo/FBS/Transport"
	webrtctransport "mediasoupgo/FBS/WebRtcTransport"
	worker "mediasoupgo/FBS/Worker"
	"mediasoupgo/events"
	"mediasoupgo/ptr"
	"mediasoupgo/smap"
)

var _ Router = &routerImpl{}

type RouterInternal struct {
	routerId string
}

func (r RouterInternal) RouterId() string {
	return r.routerId
}

type routerImpl struct {
	RouterInternal
	rtpCapabilities *RtpCapabilities
	channel         *Channel
	closed          atomic.Bool
	appData         RouterAppData
	transports      *smap.Map[string, Transport]
	producers       *smap.Map[string, Producer]
	rtpObservers    *smap.Map[string, RtpObserver]
	dataProudcers   *smap.Map[string, DataProducer]
	// TODO
	mapRouterPairPipeTransportPairPromise *smap.Map[string, PipeTransportPair]
	observer                              RouterObserver
	events.EventEmmiter[RouterEvents]
}

func NewRouter(
	id string,
	rtpCaps *RtpCapabilities,
	channel *Channel,
	appData RouterAppData,
) (Router, error) {
	r := &routerImpl{
		RouterInternal:                        RouterInternal{routerId: id},
		rtpCapabilities:                       rtpCaps,
		channel:                               channel,
		closed:                                atomic.Bool{},
		appData:                               appData,
		transports:                            smap.New[string, Transport](),
		producers:                             smap.New[string, Producer](),
		rtpObservers:                          smap.New[string, RtpObserver](),
		dataProudcers:                         smap.New[string, DataProducer](),
		mapRouterPairPipeTransportPairPromise: smap.New[string, PipeTransportPair](),
		observer:                              events.New[RouterObserverEvents](),
		EventEmmiter:                          events.New[RouterEvents](),
	}
	r.handListenError()
	return r, nil
}

func (r *routerImpl) Id() string {
	return r.routerId
}

func (r *routerImpl) Closed() bool {
	return r.closed.Load()
}

func (r *routerImpl) RtpCapabilities() *RtpCapabilities {
	return r.rtpCapabilities
}

func (r *routerImpl) AppData() RouterAppData {
	return r.appData
}

func (r *routerImpl) SetAppData(appData RouterAppData) {
	r.appData = appData
}

func (r *routerImpl) Observer() RouterObserver {
	return r.observer
}

func (r *routerImpl) Close() {
	if r.closed.Load() {
		return
	}
	r.channel.Request(Request.MethodWORKER_CLOSE_ROUTER, &Request.BodyT{
		Type:  Request.BodyWorker_CloseRouterRequest,
		Value: worker.CloseRouterRequestT{RouterId: r.routerId},
	}, "")

	r.transports.Range(func(key string, value Transport) bool {
		value.RouterClosed()
		return true
	})
	r.transports = smap.New[string, Transport]()
	r.producers = smap.New[string, Producer]()
	r.dataProudcers = smap.New[string, DataProducer]()
	r.rtpObservers.Range(func(key string, value RtpObserver) bool {
		value.RouterClosed()
		return true
	})
	r.rtpObservers = smap.New[string, RtpObserver]()
	r.Emit("@close", RouterEvents{AtClose: struct{}{}})
	r.observer.Emit("close", RouterObserverEvents{Close: struct{}{}})
}

func (r *routerImpl) WorkerClosed() {
	if r.closed.Load() {
		return
	}
	r.transports.Range(func(key string, value Transport) bool {
		value.RouterClosed()
		return true
	})
	r.transports = smap.New[string, Transport]()
	r.producers = smap.New[string, Producer]()
	r.dataProudcers = smap.New[string, DataProducer]()
	r.rtpObservers.Range(func(key string, value RtpObserver) bool {
		value.RouterClosed()
		return true
	})
	r.rtpObservers = smap.New[string, RtpObserver]()
	r.Emit("workerclose", RouterEvents{WorkerClose: struct{}{}})
	r.observer.Emit("close", RouterObserverEvents{Close: struct{}{}})
}

func (r *routerImpl) Dump(ctx context.Context) (*RouterDump, error) {
	resp, err := r.channel.Request(Request.MethodROUTER_DUMP, nil, r.Id())
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)
	ret := &RouterDump{}
	return ret, nil
}

func (r *routerImpl) CreateWebRtcTransport(
	options *WebRtcTransportOptions,
) (WebRtcTransport, error) {
	if options.WebRtcTransportListen == nil {
		return nil, errors.New("miss listren")
	}
	if options.WebRtcTransportListen.WebRtcServer == nil &&
		len(options.WebRtcTransportListen.ListenInfos) == 0 {
		return nil, errors.New("empty WebRtcTransportListen")
	}
	if options.WebRtcTransportListen.WebRtcServer != nil &&
		len(options.WebRtcTransportListen.ListenInfos) > 0 {
		return nil, errors.New("just choose one type ")
	}
	preferUdp := false
	if options.PreferUdp != nil {
		preferUdp = *options.EnableUdp
	}
	preferTcp := false
	if options.PreferTcp != nil {
		preferTcp = *options.PreferTcp
	}
	initialavailableoutgoingbitrate := uint32(600000)
	if options.InitialAvailableOutgoingBitrate != nil {
		initialavailableoutgoingbitrate = uint32(*options.InitialAvailableOutgoingBitrate)
	}
	enableSctp := false
	if options.EnableSctp != nil {
		enableSctp = *options.EnableSctp
	}

	numSctpStreams := options.NumSctpStreams
	if options.NumSctpStreams == nil {
		numSctpStreams = &NumSctpStreams{OS: 1024, MIS: 1024}
	}
	maxSctpMessageSzie := 2621244
	sctpSendBufferSize := 2621244
	if options.MaxSctpMessageSize != nil {
		maxSctpMessageSzie = *options.MaxSctpMessageSize
	}
	if options.SctpSendBufferSize != nil {
		sctpSendBufferSize = *options.SctpSendBufferSize
	}
	iceConsentTimeout := 30
	if options.IceConsentTimeout != nil {
		iceConsentTimeout = *options.IceConsentTimeout
	}
	enableUdp := true
	enableTcp := true
	if options.EnableUdp != nil {
		enableUdp = *options.EnableUdp
	}
	if options.EnableTcp != nil {
		enableTcp = *options.EnableTcp
	}
	transportId := uuid.NewString()
	req := &router.CreateWebRtcTransportRequestT{
		TransportId: transportId,
	}

	req.Options = &webrtctransport.WebRtcTransportOptionsT{
		Base: &transport.OptionsT{
			Direct:                          false,
			MaxMessageSize:                  nil,
			InitialAvailableOutgoingBitrate: &initialavailableoutgoingbitrate,
			EnableSctp:                      enableSctp,
			NumSctpStreams: &sctpparameters.NumSctpStreamsT{
				Os:  uint16(numSctpStreams.OS),
				Mis: uint16(numSctpStreams.MIS),
			},
			MaxSctpMessageSize: uint32(maxSctpMessageSzie),
			SctpSendBufferSize: uint32(sctpSendBufferSize),
			IsDataChannel:      true,
		},
		EnableUdp:         enableUdp,
		EnableTcp:         enableTcp,
		PreferUdp:         preferUdp,
		PreferTcp:         preferTcp,
		IceConsentTimeout: byte(iceConsentTimeout),
	}
	if options.WebRtcTransportListen.WebRtcServer != nil {
		req.Options.Listen = &webrtctransport.ListenT{
			Type: webrtctransport.ListenListenServer,
			Value: &webrtctransport.ListenServerT{
				WebRtcServerId: options.WebRtcTransportListen.WebRtcServer.Id(),
			},
		}
	}
	if len(options.WebRtcTransportListen.ListenInfos) > 0 {
		list := options.WebRtcTransportListen
		linstens := &webrtctransport.ListenIndividualT{}
		for _, v := range list.ListenInfos {
			info := &transport.ListenInfoT{
				Protocol:  transport.EnumValuesProtocol[strings.ToUpper(string(v.Protocol))],
				Ip:        v.IP,
				PortRange: &transport.PortRangeT{},
				Flags:     &transport.SocketFlagsT{},
			}
			if v.AnnouncedAddress != nil {
				info.AnnouncedAddress = *v.AnnouncedAddress
			}
			if v.Port != nil {
				info.Port = *v.Port
			}
			if v.SendBufferSize != nil {
				info.SendBufferSize = *v.SendBufferSize
			}
			if v.RecvBufferSize != nil {
				info.RecvBufferSize = *v.RecvBufferSize
			}
			if v.PortRange != nil {
				info.PortRange = &transport.PortRangeT{Max: v.PortRange.Max, Min: v.PortRange.Min}
			}
			if v.Flags != nil {
				info.Flags = &transport.SocketFlagsT{
					Ipv6Only:     v.Flags.IPV6Only,
					UdpReusePort: v.Flags.UDPReusePort,
				}
			}
			linstens.ListenInfos = append(linstens.ListenInfos, info)
		}
		req.Options.Listen = &webrtctransport.ListenT{
			Type:  webrtctransport.ListenListenIndividual,
			Value: linstens,
		}
	}

	resp, err := r.channel.Request(
		Request.MethodROUTER_CREATE_WEBRTCTRANSPORT_WITH_SERVER,
		&Request.BodyT{
			Type:  Request.BodyRouter_CreateWebRtcTransportRequest,
			Value: req,
		},
		r.routerId)
	if err != nil {
		return nil, err
	}
	resp2 := resp.Body.Value.(*webrtctransport.DumpResponseT)
	data := ToWebRtcTransportData(resp2)

	p := NewWebRtcTransport(
		data,
		TransportInternal{
			RouterInternal: RouterInternal{routerId: r.routerId},
			transportId:    transportId,
		},
		r.channel,
		TransportAppData(r.appData),
		func() RtpCapabilities {
			return *r.rtpCapabilities
		}, func(s string) Producer {
			p, _ := r.producers.Get(s)
			return p
		}, func(s string) DataProducer {
			dp, _ := r.dataProudcers.Get(s)
			return dp
		},
	)
	id := transportId
	r.transports.Set(id, p)
	p.On("@close", func(arg WebRtcTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@listenserverclose", func(arg WebRtcTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@newproducer", func(arg WebRtcTransportEvents) {
		r.producers.Set(arg.Newproducer.Arg1.ID(), arg.Newproducer.Arg1)
	})
	p.On("@producerclose", func(arg WebRtcTransportEvents) {
		r.producers.Delete(arg.Producerclose.Arg1.ID())
	})
	p.On("@newdataproducer", func(arg WebRtcTransportEvents) {
		r.dataProudcers.Set(arg.Newdataproducer.Arg1.ID(), arg.Newdataproducer.Arg1)
	})
	p.On("@dataproducerclose", func(arg WebRtcTransportEvents) {
		r.dataProudcers.Delete(arg.Dataproducerclose.Arg1.ID())
	})
	r.observer.Emit(
		"newtransport",
		RouterObserverEvents{Newtransport: events.NewEvent1(Transport(p))},
	)
	if options.WebRtcTransportListen != nil && options.WebRtcTransportListen.WebRtcServer != nil {
		options.WebRtcTransportListen.WebRtcServer.HandleWebRtcTransport(p)
	}
	return p, nil
}

func ToWebRtcTransportData(resp2 *webrtctransport.DumpResponseT) *WebRtcTransportData {
	data := &WebRtcTransportData{
		iceState:       IceState(strings.ToLower(resp2.IceState.String())),
		dtlsState:      DtlsState(strings.ToLower(resp2.DtlsState.String())),
		dtlsRemoteCert: nil,
		sctpState:      SctpState(strings.ToLower(resp2.Base.SctpState.String())),
	}
	if resp2.Base != nil {
		data.sctpParameters = &SctpParameters{
			Port:           resp2.Base.SctpParameters.Port,
			OS:             resp2.Base.SctpParameters.Os,
			MIS:            resp2.Base.SctpParameters.Mis,
			MaxMessageSize: resp2.Base.MaxMessageSize,
		}
	}
	if resp2.IceParameters != nil {
		data.iceParameters = IceParameters{
			UsernameFragment: resp2.IceParameters.UsernameFragment,
			Password:         resp2.IceParameters.Password,
			IceLite:          &resp2.IceParameters.IceLite,
		}
	}
	if resp2.IceSelectedTuple != nil {
		data.iceSelectedTuple = &TransportTuple{
			LocalIP:      resp2.IceSelectedTuple.LocalAddress,
			LocalAddress: resp2.IceSelectedTuple.LocalAddress,
			LocalPort:    resp2.IceSelectedTuple.LocalPort,
			RemoteIP:     &resp2.IceSelectedTuple.RemoteIp,
			RemotePort:   &resp2.IceSelectedTuple.RemotePort,
			Protocol: TransportProtocol(
				strings.ToLower(resp2.IceSelectedTuple.Protocol.String()),
			),
		}
	}
	if len(resp2.IceCandidates) > 0 {
		var iceCandidates []*IceCandidate
		for _, v := range resp2.IceCandidates {
			item := &IceCandidate{
				Foundation: v.Foundation,
				Priority:   v.Priority,
				IP:         v.Address,
				Address:    v.Address,
				Protocol:   TransportProtocol(strings.ToLower(v.Protocol.String())),
				Port:       v.Port,
				Type:       IceCandidateType(strings.ToLower(v.Type.String())),
			}
			if v.TcpType != nil {
				item.TCPType = ptr.To(IceCandidateTcpType(strings.ToLower(v.TcpType.String())))
			}
			iceCandidates = append(iceCandidates, item)
		}
		data.iceCandidates = iceCandidates
	}
	if resp2.DtlsParameters != nil {
		data.dtlsParameters.Role = ptr.To(
			DtlsRole(strings.ToLower(resp2.DtlsParameters.Role.String())),
		)
		var figers []DtlsFingerprint
		for _, v := range resp2.DtlsParameters.Fingerprints {
			figers = append(figers, ToDtlsFingerprint(v))
		}
		data.dtlsParameters.Fingerprints = figers
	}
	return data
}

func (r *routerImpl) CreatePlainTransport(options *PlainTransportOptions) (PlainTransport, error) {
	// Implementation for creating PlainTransport
	id := uuid.NewString()
	req := &router.CreatePlainTransportRequestT{
		TransportId: "",
		Options: &fbsplaintransport.PlainTransportOptionsT{
			Base: &transport.OptionsT{
				Direct:                          false,
				MaxMessageSize:                  new(uint32),
				InitialAvailableOutgoingBitrate: new(uint32),
				EnableSctp:                      false,
				NumSctpStreams:                  &sctpparameters.NumSctpStreamsT{},
				MaxSctpMessageSize:              0,
				SctpSendBufferSize:              0,
				IsDataChannel:                   false,
			},
			ListenInfo: &transport.ListenInfoT{
				Protocol:         0,
				Ip:               "",
				AnnouncedAddress: "",
				Port:             0,
				PortRange:        &transport.PortRangeT{},
				Flags:            &transport.SocketFlagsT{},
				SendBufferSize:   0,
				RecvBufferSize:   0,
			},
			RtcpListenInfo: &transport.ListenInfoT{
				Protocol:         0,
				Ip:               "",
				AnnouncedAddress: "",
				Port:             0,
				PortRange:        &transport.PortRangeT{},
				Flags:            &transport.SocketFlagsT{},
				SendBufferSize:   0,
				RecvBufferSize:   0,
			},
			RtcpMux:         false,
			Comedia:         false,
			EnableSrtp:      false,
			SrtpCryptoSuite: nil,
		},
	}
	resp, err := r.channel.Request(
		Request.MethodROUTER_CREATE_PLAINTRANSPORT,
		&Request.BodyT{Type: Request.BodyRouter_CreatePlainTransportRequest, Value: req},
		r.routerId,
	)
	if err != nil {
		return nil, err
	}
	resp2 := resp.Body.Value.(*fbsplaintransport.DumpResponseT)

	p := NewPlainTransport(&PlainTransportData{
		rtcpMux:        resp2.RtcpMux,
		comedia:        false,
		tuple:          TransportTuple{},
		rtcpTuple:      &TransportTuple{},
		sctpParameters: &SctpParameters{},
		sctpState:      "",
		srtpParameters: &SrtpParameters{},
	}, TransportInternal{RouterInternal: r.RouterInternal, transportId: id}, r.channel, TransportAppData(r.appData),
		func() RtpCapabilities {
			return *r.rtpCapabilities
		}, func(s string) Producer {
			x, _ := r.producers.Get(s)
			return x
		}, func(s string) DataProducer {
			x, _ := r.dataProudcers.Get(s)
			return x
		},
	)
	r.transports.Set(id, p)
	p.On("@close", func(arg PlainTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@listenserverclose", func(arg PlainTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@newproducer", func(arg PlainTransportEvents) {
		r.producers.Set(arg.Newdataproducer.Arg1.ID(), arg.Newproducer.Arg1)
	})
	p.On("@producerclose", func(arg PlainTransportEvents) {
		r.producers.Delete(arg.Producerclose.Arg1.ID())
	})
	p.On("@newdataproducer", func(arg PlainTransportEvents) {
		r.dataProudcers.Set(arg.Newdataproducer.Arg1.ID(), arg.Newdataproducer.Arg1)
	})
	p.On("@dataproducerclose", func(arg PlainTransportEvents) {
		r.dataProudcers.Delete(arg.Dataproducerclose.Arg1.ID())
	})
	r.observer.Emit(
		"newtransport",
		RouterObserverEvents{Newtransport: events.NewEvent1(Transport(p))},
	)
	return p, nil
}

func (r *routerImpl) CreatePipeTransport(options *PipeTransportOptions) (PipeTransport, error) {
	id := uuid.NewString()
	req := &router.CreatePipeTransportRequestT{
		TransportId: "",
		Options: &fbspipetransport.PipeTransportOptionsT{
			Base:       &transport.OptionsT{},
			ListenInfo: &transport.ListenInfoT{},
			EnableRtx:  false,
			EnableSrtp: false,
		},
	}
	resp, err := r.channel.Request(
		Request.MethodROUTER_CREATE_PIPETRANSPORT,
		&Request.BodyT{Type: Request.BodyRouter_CreatePipeTransportRequest, Value: req},
		r.routerId,
	)
	if err != nil {
		return nil, err
	}
	resp2 := resp.Body.Value.(*fbspipetransport.DumpResponseT)
	p := NewPipeTransport(&PipeTransportData{
		tuple:          TransportTuple{},
		sctpParameters: &SctpParameters{},
		sctpState:      "",
		rtx:            resp2.Rtx,
		srtpParameters: &SrtpParameters{},
	}, TransportInternal{RouterInternal: r.RouterInternal, transportId: id}, r.channel, TransportAppData(r.appData),
		func() RtpCapabilities {
			return *r.rtpCapabilities
		}, func(s string) Producer {
			x, _ := r.producers.Get(s)
			return x
		}, func(s string) DataProducer {
			x, _ := r.dataProudcers.Get(s)
			return x
		},
	)

	r.transports.Set(id, p)
	p.On("@close", func(arg PipeTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@listenserverclose", func(arg PipeTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@newproducer", func(arg PipeTransportEvents) {
		r.producers.Set(arg.Newdataproducer.Arg1.ID(), arg.Newproducer.Arg1)
	})
	p.On("@producerclose", func(arg PipeTransportEvents) {
		r.producers.Delete(arg.Producerclose.Arg1.ID())
	})
	p.On("@newdataproducer", func(arg PipeTransportEvents) {
		r.dataProudcers.Set(arg.Newdataproducer.Arg1.ID(), arg.Newdataproducer.Arg1)
	})
	p.On("@dataproducerclose", func(arg PipeTransportEvents) {
		r.dataProudcers.Delete(arg.Dataproducerclose.Arg1.ID())
	})
	r.observer.Emit(
		"newtransport",
		RouterObserverEvents{Newtransport: events.NewEvent1(Transport(p))},
	)
	return p, nil
}

func (r *routerImpl) CreateDirectTransport(
	options *DirectTransportOptions,
) (DirectTransport, error) {
	id := uuid.NewString()
	req := &router.CreateDirectTransportRequestT{
		TransportId: id,
		Options: &fbsdirecttransport.DirectTransportOptionsT{
			Base: &transport.OptionsT{
				Direct:                          false,
				MaxMessageSize:                  new(uint32),
				InitialAvailableOutgoingBitrate: new(uint32),
				EnableSctp:                      false,
				NumSctpStreams:                  &sctpparameters.NumSctpStreamsT{},
				MaxSctpMessageSize:              0,
				SctpSendBufferSize:              0,
				IsDataChannel:                   false,
			},
		},
	}
	resp, err := r.channel.Request(
		Request.MethodROUTER_CREATE_DIRECTTRANSPORT,
		&Request.BodyT{Type: Request.BodyRouter_CreateDirectTransportRequest, Value: req},
		r.routerId,
	)
	if err != nil {
		return nil, err
	}
	resp2 := resp.Body.Value.(*fbsdirecttransport.DumpResponseT)
	fmt.Println(resp2)
	p := NewDirectTransport(&DirectTransportData{
		sctpParameters: &SctpParameters{},
	}, TransportInternal{RouterInternal: r.RouterInternal, transportId: id}, r.channel, TransportAppData(r.appData),
		func() RtpCapabilities {
			return *r.rtpCapabilities
		}, func(s string) Producer {
			x, _ := r.producers.Get(s)
			return x
		}, func(s string) DataProducer {
			x, _ := r.dataProudcers.Get(s)
			return x
		},
	)

	r.transports.Set(id, p)
	p.On("@close", func(arg DirectTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@listenserverclose", func(arg DirectTransportEvents) {
		r.transports.Delete(id)
	})
	p.On("@newproducer", func(arg DirectTransportEvents) {
		r.producers.Set(arg.Newdataproducer.Arg1.ID(), arg.Newproducer.Arg1)
	})
	p.On("@producerclose", func(arg DirectTransportEvents) {
		r.producers.Delete(arg.Producerclose.Arg1.ID())
	})
	p.On("@newdataproducer", func(arg DirectTransportEvents) {
		r.dataProudcers.Set(arg.Newdataproducer.Arg1.ID(), arg.Newdataproducer.Arg1)
	})
	p.On("@dataproducerclose", func(arg DirectTransportEvents) {
		r.dataProudcers.Delete(arg.Dataproducerclose.Arg1.ID())
	})
	r.observer.Emit(
		"newtransport",
		RouterObserverEvents{Newtransport: events.NewEvent1(Transport(p))},
	)
	return &directTransportImpl{}, nil
}

func (r *routerImpl) PipeToRouter(options *PipeToRouterOptions) (PipeToRouterResult, error) {
	return PipeToRouterResult{}, nil
}

func (r *routerImpl) AddPipeTransportPair(
	pipeTransportPairKey string,
	pipeTransportPairPromise PipeTransportPair,
) {
}

func (r *routerImpl) CreateActiveSpeakerObserver(
	options *ActiveSpeakerObserverOptions,
) (ActiveSpeakerObserver, error) {
	id := uuid.NewString()
	req := &router.CreateActiveSpeakerObserverRequestT{
		RtpObserverId: id,
		Options: &fbsactivespeakerobserver.ActiveSpeakerObserverOptionsT{
			Interval: uint16(options.Interval),
		},
	}
	r.channel.Request(
		Request.MethodROUTER_CREATE_ACTIVESPEAKEROBSERVER,
		&Request.BodyT{Type: Request.BodyRouter_CreateActiveSpeakerObserverRequest, Value: req},
		r.routerId,
	)
	aso := NewActiveSpeakerObserver(
		RtpObserverInternal{RouterInternal: r.RouterInternal, rtpObserverId: id},
		r.channel,
		RtpObserverAppData(r.appData),
		func(producerId string) Producer {
			x, _ := r.producers.Get(producerId)
			return x
		},
	)
	aso.On("@close", func(arg ActiveSpeakerObserverEvents) {
		r.rtpObservers.Delete(aso.ID())
	})
	r.observer.Emit(
		"newrtpobserver",
		RouterObserverEvents{Newrtpobserver: events.NewEvent1(RtpObserver(aso))},
	)
	return aso, nil
}

func (r *routerImpl) CreateAudioLevelObserver(
	options *AudioLevelObserverOptions,
) (AudioLevelObserver, error) {
	maxEntries := options.MaxEntries
	id := uuid.NewString()
	req := &router.CreateAudioLevelObserverRequestT{
		RtpObserverId: id,
		Options: &fbsaudiolevelobserver.AudioLevelObserverOptionsT{
			MaxEntries: uint16(maxEntries),
			Threshold:  int8(options.Threshold),
			Interval:   uint16(options.Interval),
		},
	}
	r.channel.Request(
		Request.MethodROUTER_CREATE_AUDIOLEVELOBSERVER,
		&Request.BodyT{Type: Request.BodyRouter_CreateAudioLevelObserverRequest, Value: req},
		r.routerId,
	)
	alo := NewAudioLevelObserver(
		RtpObserverInternal{RouterInternal: r.RouterInternal, rtpObserverId: id},
		r.channel,
		RtpObserverAppData(r.appData),
		func(producerId string) Producer {
			x, _ := r.producers.Get(producerId)
			return x
		},
	)
	alo.On("@close", func(arg AudioLevelObserverEvents) {
		r.rtpObservers.Delete(alo.ID())
	})

	r.rtpObservers.Set(id, alo)
	r.observer.Emit(
		"newrtpobserver",
		RouterObserverEvents{Newrtpobserver: events.NewEvent1(RtpObserver(alo))},
	)
	return alo, nil
}

func (r *routerImpl) CanConsume(producerId string, rtpCapabilities *RtpCapabilities) bool {
	producer, ok := r.producers.Get(producerId)
	if !ok {
		return false
	}
	can, err := CanConsume(ptr.To(producer.ConsumableRTPParameters()), rtpCapabilities)
	if err != nil {
		return false
	}
	return can
}

func (r *routerImpl) handListenError() {
	r.On("listenererror", func(arg RouterEvents) {
		// TODO
	})
}
