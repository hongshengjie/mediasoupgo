package mediasoupgo

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"

	consumer "mediasoupgo/FBS/Consumer"
	dataconsumer "mediasoupgo/FBS/DataConsumer"
	dataproducer "mediasoupgo/FBS/DataProducer"
	"mediasoupgo/FBS/Request"
	router "mediasoupgo/FBS/Router"
	rtpparameters "mediasoupgo/FBS/RtpParameters"
	sctpstreamparameters "mediasoupgo/FBS/SctpParameters"
	transport "mediasoupgo/FBS/Transport"
	"mediasoupgo/events"
	"mediasoupgo/ptr"
	"mediasoupgo/smap"
)

var _ Transport = &transportImpl{}

type TransportInternal struct {
	RouterInternal
	transportId string
}

type (
	transportImpl struct {
		TransportInternal
		channel                  *Channel
		closed                   atomic.Bool
		appData                  TransportAppData
		getRouterRtpCapabilities func() RtpCapabilities
		getProducerById          func(producerId string) Producer
		getDataProducerById      func(dataProducerId string) DataProducer
		EmitEvent                func(events.EventName, TransportEvents)
		EmitObserverEvent        func(events.EventName, TransportObserverEvents)
		producers                *smap.Map[string, Producer]
		consumers                *smap.Map[string, Consumer]
		dataProducers            *smap.Map[string, DataProducer]
		dataConsumers            *smap.Map[string, DataConsumer]
		cnameForProducers        string
		nextMidForConsumers      atomic.Uint32
		sctpStreamIds            bytes.Buffer
		nextSctpStreamId         int
		typ                      string
	}
)

func NewTransport(
	id TransportInternal,
	channel *Channel,
	appData TransportAppData,
	getRouterRtpCapabilities func() RtpCapabilities,
	getProducerById func(string) Producer,
	getDataProducerById func(string) DataProducer,
	emitEvent func(events.EventName, TransportEvents),
	emitObserverEvent func(events.EventName, TransportObserverEvents),
	typ string,
) *transportImpl {
	t := &transportImpl{
		TransportInternal:        id,
		channel:                  channel,
		closed:                   atomic.Bool{},
		appData:                  appData,
		getRouterRtpCapabilities: getRouterRtpCapabilities,
		getProducerById:          getProducerById,
		getDataProducerById:      getDataProducerById,
		EmitEvent:                emitEvent,
		EmitObserverEvent:        emitObserverEvent,
		producers:                smap.New[string, Producer](),
		consumers:                smap.New[string, Consumer](),
		dataProducers:            smap.New[string, DataProducer](),
		dataConsumers:            smap.New[string, DataConsumer](),
		typ:                      typ,
	}
	return t
}

// Transport id
func (t *transportImpl) ID() string {
	return t.transportId
}

// Whether the Transport is closed
func (t *transportImpl) Closed() bool {
	return t.closed.Load()
}

// App custom data
func (t *transportImpl) AppData() TransportAppData {
	return t.appData
}

func (t *transportImpl) SetAppData(appData TransportAppData) {
	t.appData = appData
}

// Close the Transport
func (t *transportImpl) Close() {
	if t.closed.Load() {
		return
	}
	t.closed.Store(true)
	t.channel.RemoveAllListeners(events.EventName(t.transportId))
	_, err := t.channel.Request(
		Request.MethodROUTER_CLOSE_TRANSPORT,
		&Request.BodyT{
			Type:  Request.BodyRouter_CloseTransportRequest,
			Value: router.CloseTransportRequestT{TransportId: t.transportId},
		}, t.routerId)
	if err != nil {
	}

	t.producers.Range(func(key string, value Producer) bool {
		value.TransportClosed()

		t.EmitEvent("@producerclose", TransportEvents{Producerclose: events.NewEvent1(value)})
		return true
	})
	t.consumers.Range(func(key string, value Consumer) bool {
		value.TransportClosed()
		return true
	})
	t.dataProducers.Range(func(key string, value DataProducer) bool {
		value.TransportClosed()
		t.EmitEvent("@dataproducerclose", TransportEvents{Dataproducerclose: events.NewEvent1(value)})
		return true
	})
	t.dataConsumers.Range(func(key string, value DataConsumer) bool {
		value.TransportClosed()
		return true
	})

	t.producers = smap.New[string, Producer]()
	t.consumers = smap.New[string, Consumer]()
	t.dataProducers = smap.New[string, DataProducer]()
	t.dataConsumers = smap.New[string, DataConsumer]()
	t.EmitEvent("@close", TransportEvents{Close: struct{}{}})
	t.EmitObserverEvent("close", TransportObserverEvents{Close: struct{}{}})
}

// Router was closed
func (t *transportImpl) RouterClosed() {
	if t.closed.Load() {
		return
	}
	t.closed.Store(true)
	t.channel.RemoveAllListeners(events.EventName(t.transportId))

	t.producers.Range(func(key string, value Producer) bool {
		value.TransportClosed()
		return true
	})
	t.consumers.Range(func(key string, value Consumer) bool {
		value.TransportClosed()
		return true
	})
	t.dataProducers.Range(func(key string, value DataProducer) bool {
		value.TransportClosed()
		return true
	})
	t.dataConsumers.Range(func(key string, value DataConsumer) bool {
		value.TransportClosed()
		return true
	})

	t.producers = smap.New[string, Producer]()
	t.consumers = smap.New[string, Consumer]()
	t.dataProducers = smap.New[string, DataProducer]()
	t.dataConsumers = smap.New[string, DataConsumer]()
	t.EmitEvent("routerclose", TransportEvents{Close: struct{}{}})
	t.EmitObserverEvent("close", TransportObserverEvents{Close: struct{}{}})
}

// Listen server was closed
func (t *transportImpl) ListenServerClosed() {
	if t.closed.Load() {
		return
	}
	t.closed.Store(true)
	t.channel.RemoveAllListeners(events.EventName(t.transportId))

	t.producers.Range(func(key string, value Producer) bool {
		value.TransportClosed()
		return true
	})
	t.consumers.Range(func(key string, value Consumer) bool {
		value.TransportClosed()
		return true
	})
	t.dataProducers.Range(func(key string, value DataProducer) bool {
		value.TransportClosed()
		return true
	})
	t.dataConsumers.Range(func(key string, value DataConsumer) bool {
		value.TransportClosed()
		return true
	})

	t.producers = smap.New[string, Producer]()
	t.consumers = smap.New[string, Consumer]()
	t.dataProducers = smap.New[string, DataProducer]()
	t.dataConsumers = smap.New[string, DataConsumer]()
	t.EmitEvent("@listenserverclose", TransportEvents{Close: struct{}{}})
	t.EmitEvent("listenserverclose", TransportEvents{Close: struct{}{}})
	t.EmitObserverEvent("close", TransportObserverEvents{Close: struct{}{}})
}

// Set maximum incoming bitrate for receiving media
func (t *transportImpl) SetMaxIncomingBitrate(bitrate int) error {
	req := &transport.SetMaxIncomingBitrateRequestT{MaxIncomingBitrate: uint32(bitrate)}
	_, err := t.channel.Request(
		Request.MethodTRANSPORT_SET_MAX_INCOMING_BITRATE,
		&Request.BodyT{Type: Request.BodyTransport_SetMaxIncomingBitrateRequest, Value: req},
		t.transportId,
	)
	return err
}

// Set maximum outgoing bitrate for sending media
func (t *transportImpl) SetMaxOutgoingBitrate(bitrate int) error {
	req := &transport.SetMaxOutgoingBitrateRequestT{MaxOutgoingBitrate: uint32(bitrate)}
	_, err := t.channel.Request(
		Request.MethodTRANSPORT_SET_MAX_OUTGOING_BITRATE,
		&Request.BodyT{Type: Request.BodyTransport_SetMaxOutgoingBitrateRequest, Value: req},
		t.transportId,
	)
	return err
}

// Set minimum outgoing bitrate for sending media
func (t *transportImpl) SetMinOutgoingBitrate(bitrate int) error {
	req := &transport.SetMinOutgoingBitrateRequestT{MinOutgoingBitrate: uint32(bitrate)}
	_, err := t.channel.Request(
		Request.MethodTRANSPORT_SET_MIN_OUTGOING_BITRATE,
		&Request.BodyT{Type: Request.BodyTransport_SetMinOutgoingBitrateRequest, Value: req},
		t.transportId,
	)
	return err
}

// Create a Producer
func (t *transportImpl) Produce(options *ProducerOptions) (_ Producer, _ error) {
	var id string
	if options.ID != nil {
		id = *options.ID
	}
	if _, ok := t.producers.Get(id); ok {
		return nil, errors.New("a Producer with same id")
	}
	var paused bool
	if options.Paused != nil {
		paused = *options.Paused
	}
	var keyFrameRequestDelay uint32
	if options.KeyFrameRequestDelay != nil {
		keyFrameRequestDelay = *options.KeyFrameRequestDelay
	}
	rtpParameters := options.RTPParameters

	if err := ValidateRtpParameters(&rtpParameters); err != nil {
		return nil, err
	}
	if t.typ != "pipe" {
		if t.cnameForProducers == "" && rtpParameters.RTCP != nil && rtpParameters.RTCP.CNAME != nil {
			t.cnameForProducers = *rtpParameters.RTCP.CNAME
		} else {
			t.cnameForProducers = uuid.NewString()[0:8]
		}
		if rtpParameters.RTCP == nil {
			rtpParameters.RTCP = &RtcpParameters{
				CNAME: &t.cnameForProducers,
			}
		}
	}
	routerRtpCapabilities := t.getRouterRtpCapabilities()
	rtpMapping, err := GetProducerRtpParametersMapping(&rtpParameters, &routerRtpCapabilities)
	if err != nil {
		return nil, err
	}
	consumableRtpParameters := GetConsumableRtpParameters(
		string(options.Kind),
		rtpParameters, routerRtpCapabilities,
		*rtpMapping,
	)
	if id == "" {
		id = uuid.NewString()
	}
	req := &transport.ProduceRequestT{
		ProducerId:           id,
		Kind:                 rtpparameters.EnumValuesMediaKind[strings.ToUpper(string(options.Kind))],
		RtpParameters:        ToFbsRtpParameters(&options.RTPParameters),
		RtpMapping:           SerializeRtpMapping(rtpMapping),
		KeyFrameRequestDelay: keyFrameRequestDelay,
		Paused:               paused,
	}
	resp, err := t.channel.Request(
		Request.MethodTRANSPORT_PRODUCE,
		&Request.BodyT{
			Type:  Request.BodyTransport_ProduceRequest,
			Value: req,
		},
		t.transportId)
	if err != nil {
		return nil, err
	}
	resp2 := resp.Body.Value.(*transport.ProduceResponseT)
	p := NewProducer(
		ProducerInternal{
			TransportInternal: t.TransportInternal,
			producerId:        id,
		},
		&ProducerData{
			kind:                    options.Kind,
			typ:                     ProducerType(strings.ToLower(resp2.Type.String())),
			rtpParameters:           rtpParameters,
			consumableRtpParameters: consumableRtpParameters,
		},
		t.channel,
		ProducerAppData(t.appData),
		paused)
	t.producers.Set(p.ID(), p)
	p.On("@close", func(arg ProducerEvents) {
		t.producers.Delete(p.ID())
		t.EmitEvent("@producerclose", TransportEvents{Producerclose: events.NewEvent1(p)})
	})
	t.EmitEvent("@newproducer", TransportEvents{Newproducer: events.NewEvent1(p)})
	t.EmitObserverEvent("newproducer", TransportObserverEvents{Newproducer: events.NewEvent1(p)})
	return p, nil
}

// Create a Consumer
func (t *transportImpl) Consume(options *ConsumerOptions) (Consumer, error) {
	if options == nil {
		return nil, errors.New("nil options")
	}
	if options.ProducerID == "" {
		return nil, errors.New("empty produceId")
	}
	// if options.MID == nil || *options.MID == "" {
	// 	return nil, errors.New("empty mid")
	// }
	clonedRtpCapabilities := options.RTPCapabilities
	err := ValidateRtpCapabilities(&clonedRtpCapabilities)
	if err != nil {
		return nil, err
	}

	producer := t.getProducerById(options.ProducerID)
	if producer == nil {
		return nil, fmt.Errorf("no producer with id:%s ", options.ProducerID)
	}
	var enableRtx, ignoreDtx bool
	if options.EnableRtx != nil {
		enableRtx = *options.EnableRtx
	} else {
		enableRtx = producer.Kind() == "video"
	}
	if options.IgnoreDtx != nil {
		ignoreDtx = *options.IgnoreDtx
	}

	var paused, pipe bool
	if options.Paused != nil {
		paused = *options.Paused
	}
	if options.Pipe != nil {
		pipe = *options.Pipe
	}
	rtpParameters, err := GetConsumerRtpParameters(
		producer.ConsumableRTPParameters(),
		options.RTPCapabilities,
		pipe,
		enableRtx,
	)
	if err != nil {
		return nil, err
	}
	if !pipe {
		if options.MID != nil {
			rtpParameters.MID = options.MID
		} else {
			mid := t.nextMidForConsumers.Add(1)
			rtpParameters.MID = ptr.To(strconv.Itoa(int(mid)))
			if mid > 100000000 {
				t.nextMidForConsumers.Store(0)
			}
		}
	}
	var layer *consumer.ConsumerLayersT
	if options.PreferredLayers != nil {
		layer = &consumer.ConsumerLayersT{
			SpatialLayer:  options.PreferredLayers.SpatialLayer,
			TemporalLayer: options.PreferredLayers.TemporalLayer,
		}
	}
	var typ rtpparameters.Type
	var ctyp ConsumerType
	if t.typ == "pipe" {
		typ = rtpparameters.TypePIPE
		ctyp = PipeConsumerType
	} else {
		switch producer.Type() {
		case SimpleProducerType:
			typ = rtpparameters.TypeSIMPLE
			ctyp = SimpleConsumerType
		case SvcProducerType:
			typ = rtpparameters.TypeSVC
			ctyp = SvcConsumerType
		case SimulcastProducerType:
			typ = rtpparameters.TypeSIMULCAST
			ctyp = SimulcastConsumerType
		}
	}
	consumerId := uuid.NewString()
	req := &transport.ConsumeRequestT{
		ConsumerId:             consumerId,
		ProducerId:             options.ProducerID,
		Kind:                   rtpparameters.EnumValuesMediaKind[strings.ToUpper(string(producer.Kind()))],
		RtpParameters:          ToFbsRtpParameters(rtpParameters),
		Type:                   typ,
		ConsumableRtpEncodings: ToFBSRtpEncodingParameters(producer.ConsumableRTPParameters().Encodings),
		Paused:                 paused,
		PreferredLayers:        layer,
		IgnoreDtx:              ignoreDtx,
	}
	resp, err := t.channel.Request(
		Request.MethodTRANSPORT_CONSUME,
		&Request.BodyT{Type: Request.BodyTransport_ConsumeRequest, Value: req},
		t.transportId,
	)
	if err != nil {
		return nil, err
	}

	resp2 := resp.Body.Value.(*transport.ConsumeResponseT)
	var cscore *ConsumerScore
	var clayers *ConsumerLayers
	if resp2.Score != nil {
		cscore = &ConsumerScore{
			Score:          resp2.Score.Score,
			ProducerScore:  resp2.Score.ProducerScore,
			ProducerScores: resp2.Score.ProducerScores,
		}
	}
	if resp2.PreferredLayers != nil {
		clayers = &ConsumerLayers{
			SpatialLayer:  resp2.PreferredLayers.SpatialLayer,
			TemporalLayer: resp2.PreferredLayers.TemporalLayer,
		}
	}
	c := NewConsumer(
		ConsumerInternal{TransportInternal: t.TransportInternal, consumerId: consumerId},
		&ConsumerData{
			producerId:    options.ProducerID,
			kind:          producer.Kind(),
			rtpParameters: *rtpParameters,
			typ:           ctyp,
		},
		t.channel,
		ConsumerAppData(t.appData),
		resp2.Paused,
		resp2.ProducerPaused,
		cscore,
		clayers,
	)
	t.consumers.Set(consumerId, c)
	c.On("@close", func(arg ConsumerEvents) { t.consumers.Delete(consumerId) })
	c.On("@prodcuerclose", func(arg ConsumerEvents) { t.consumers.Delete(consumerId) })
	t.EmitObserverEvent("newconsumer", TransportObserverEvents{Newconsumer: events.NewEvent1(c)})
	return c, nil
}

// Create a DataProducer
func (t *transportImpl) ProduceData(options *DataProducerOptions) (DataProducer, error) {
	var id string
	if options.ID != nil {
		id = *options.ID
	}

	if id == "" {
		id = uuid.NewString()
	}
	var tye dataproducer.Type
	var dtye DataProducerType
	clonedSctpStreamParameters := *options.SCTPStreamParameters

	req := &transport.ProduceDataRequestT{
		DataProducerId: id,
	}
	if t.typ != "direct" {
		dtye = SCTPDataProducerType
		tye = dataproducer.TypeSCTP
		if err := ValidateSctpStreamParameters(&clonedSctpStreamParameters); err != nil {
			return nil, err
		}
		req.Type = tye
	} else {
		dtye = DirectDataProducerType
		tye = dataproducer.TypeDIRECT
		clonedSctpStreamParameters = SctpStreamParameters{}
		req.Type = tye
	}
	if options.Label != nil {
		req.Label = *options.Label
	}
	if options.Protocol != nil {
		req.Protocol = *options.Protocol
	}
	if options.Paused != nil {
		req.Paused = *options.Paused
	}
	resp, err := t.channel.Request(
		Request.MethodTRANSPORT_PRODUCE_DATA,
		&Request.BodyT{Type: Request.BodyTransport_ProduceDataRequest, Value: req},
		t.transportId,
	)
	resp2 := resp.Body.Value.(*dataproducer.DumpResponseT)
	p := NewDataProducer(
		DataProducerInternal{TransportInternal: t.TransportInternal, dataProducerId: id},
		&DataProducerData{
			typ:                  dtye,
			sctpStreamParameters: &clonedSctpStreamParameters,
			lable:                resp2.Label,
			protocol:             resp2.Protocol,
		},
		t.channel,
		resp2.Paused,
		AppData(t.appData),
	)
	t.dataProducers.Set(id, p)
	p.On("@close", func(arg DataProducerEvents) {
		t.dataProducers.Delete(id)
		t.EmitEvent("@dataproducerclose", TransportEvents{Dataproducerclose: events.NewEvent1(p)})
	})
	t.EmitEvent("@newdataproducer", TransportEvents{Newdataproducer: events.NewEvent1(p)})
	t.EmitObserverEvent(
		"newdataproducer",
		TransportObserverEvents{Newdataproducer: events.NewEvent1(p)},
	)
	return p, err
}

// Create a DataConsumer
func (t *transportImpl) ConsumeData(options *DataConsumerOptions) (DataConsumer, error) {
	if options.DataProducerID == "" {
		return nil, errors.New("empty dataProducerId")
	}
	dataProducer := t.getDataProducerById(options.DataProducerID)
	if dataProducer == nil {
		return nil, fmt.Errorf("can not find dataProducer id:%s", options.DataProducerID)
	}
	var tye dataproducer.Type
	var sctpStreamParameters *SctpStreamParameters
	var sctpStreamId uint16
	if t.typ == "direct" {
		tye = dataproducer.TypeDIRECT
	} else {
		tye = dataproducer.TypeSCTP
		sctpStreamParameters = &SctpStreamParameters{
			StreamID:          sctpStreamId,
			Ordered:           new(bool),
			MaxPacketLifeTime: new(uint16),
			MaxRetransmits:    new(uint16),
		}
	}

	cid := uuid.NewString()
	lable := dataProducer.Label()
	protocol := dataProducer.Protocol()
	var paused bool
	if options.Paused != nil {
		paused = *options.Paused
	}
	req := &transport.ConsumeDataRequestT{
		DataConsumerId: cid,
		DataProducerId: options.DataProducerID,
		Type:           tye,
		SctpStreamParameters: &sctpstreamparameters.SctpStreamParametersT{
			StreamId:          sctpStreamId,
			Ordered:           sctpStreamParameters.Ordered,
			MaxPacketLifeTime: sctpStreamParameters.MaxPacketLifeTime,
			MaxRetransmits:    sctpStreamParameters.MaxRetransmits,
		},
		Label:       lable,
		Protocol:    protocol,
		Paused:      paused,
		Subchannels: options.Subchannels,
	}
	resp, err := t.channel.Request(
		Request.MethodTRANSPORT_CONSUME_DATA,
		&Request.BodyT{Type: Request.BodyTransport_ConsumeDataRequest, Value: req},
		t.transportId,
	)
	resp2 := resp.Body.Value.(*dataconsumer.DumpResponseT)
	c := NewDataComsumer(
		DataConsumerInternal{TransportInternal: t.TransportInternal, dataConsumerId: cid},
		&DataConsumerData{
			dataProducerId: resp2.DataProducerId,
			typ:            DataConsumerType(strings.ToLower(resp2.Type.String())),
			sctpStreamParameters: &SctpStreamParameters{
				StreamID:          resp2.SctpStreamParameters.StreamId,
				Ordered:           resp2.SctpStreamParameters.Ordered,
				MaxPacketLifeTime: resp2.SctpStreamParameters.MaxPacketLifeTime,
				MaxRetransmits:    resp2.SctpStreamParameters.MaxRetransmits,
			},
			lable:                      resp2.Label,
			protocol:                   resp2.Protocol,
			bufferedAmountLowThreshold: resp2.BufferedAmountLowThreshold,
		},
		t.channel,
		resp2.Paused,
		resp2.DataProducerPaused,
		resp2.Subchannels,
		options.AppData,
	)
	t.dataConsumers.Set(cid, c)
	c.On("@close", func(arg DataConsumerEvents) {
		t.dataConsumers.Delete(cid)
	})
	c.On("@dataproducerclose", func(arg DataConsumerEvents) {
		t.dataConsumers.Delete(cid)
	})
	t.EmitObserverEvent(
		"newdataconsumer",
		TransportObserverEvents{Newdataconsumer: events.NewEvent1(c)},
	)
	return c, err
}

// Enable 'trace' event
func (t *transportImpl) EnableTraceEvent(types []TransportTraceEventType) (_ error) {
	var traces []transport.TraceEventType
	for _, v := range types {
		traces = append(traces, transport.EnumValuesTraceEventType[strings.ToUpper(string(v))])
	}
	_, err := t.channel.Request(Request.MethodTRANSPORT_ENABLE_TRACE_EVENT, &Request.BodyT{
		Type: Request.BodyTransport_EnableTraceEventRequest,
		Value: &transport.EnableTraceEventRequestT{
			Events: traces,
		},
	}, t.transportId)
	return err
}
