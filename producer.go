package mediasoupgo

import (
	"log/slog"
	"mediasoupgo/internal/FBS/Notification"
	producer "mediasoupgo/internal/FBS/Producer"
	"mediasoupgo/internal/FBS/Request"
	transport "mediasoupgo/internal/FBS/Transport"
	"mediasoupgo/internal/events"
)

var _ Producer = &producerImpl{}

type ProducerInternal struct {
	TransportInternal
	producerId string
}
type ProducerData struct {
	kind                    MediaKind
	rtpParameters           RtpParameters
	typ                     ProducerType
	consumableRtpParameters RtpParameters
}
type producerImpl struct {
	ProducerInternal
	channel *Channel

	data *ProducerData

	closed   bool
	paused   bool
	appData  ProducerAppData
	score    []*ProducerScore
	observer ProducerObserver
	events.EventEmmiter[ProducerEvents]
}

func NewProducer(internal ProducerInternal, data *ProducerData, channel *Channel, appData ProducerAppData, paused bool) Producer {
	impl := &producerImpl{
		ProducerInternal: internal,
		channel:          channel,
		data:             data,
		closed:           false,
		paused:           paused,
		appData:          appData,
		score:            []*ProducerScore{},
		observer:         events.New[ProducerObserverEvents](),
		EventEmmiter:     events.New[ProducerEvents](),
	}
	impl.handleWorkerNotifications()
	impl.handleListenerError()
	return impl
}

// Producer id
func (p *producerImpl) ID() string {
	return p.producerId
}

// Whether the Producer is closed
func (p *producerImpl) Closed() bool {
	return p.closed
}

// Media kind
func (p *producerImpl) Kind() MediaKind {
	return p.data.kind
}

// RTP parameters
func (p *producerImpl) RTPParameters() RtpParameters {
	return p.data.rtpParameters
}

// Producer type
func (p *producerImpl) Type() ProducerType {
	return p.data.typ
}

// Consumable RTP parameters
func (p *producerImpl) ConsumableRTPParameters() RtpParameters {
	return p.data.consumableRtpParameters
}

// Whether the Producer is paused
func (p *producerImpl) Paused() bool {
	return p.paused
}

// Producer score list
func (p *producerImpl) Score() []*ProducerScore {
	return p.score
}

// App custom data
func (p *producerImpl) AppData() AppData {
	return AppData(p.appData)
}

func (p *producerImpl) SetAppData(appData AppData) {
	p.appData = ProducerAppData(appData)
}

// Observer
func (p *producerImpl) Observer() ProducerObserver {
	return p.observer
}

// Close the Producer
func (p *producerImpl) Close() {
	if p.closed {
		return
	}
	p.closed = true
	p.channel.RemoveAllListeners(events.EventName(p.producerId))
	p.channel.Request(Request.MethodTRANSPORT_CLOSE_PRODUCER,
		&Request.BodyT{
			Type:  Request.BodyTransport_CloseProducerRequest,
			Value: transport.CloseProducerRequestT{ProducerId: p.producerId},
		}, p.producerId,
	)
	p.Emit("@close", ProducerEvents{AtClose: struct{}{}})
	p.observer.Emit("close", ProducerObserverEvents{Close: struct{}{}})
}

// Transport was closed
func (p *producerImpl) TransportClosed() {

	if p.closed {
		return
	}
	p.closed = true
	p.channel.RemoveAllListeners(events.EventName(p.producerId))
	p.Emit("transportclose", ProducerEvents{TransportClose: struct{}{}})
	p.observer.Emit("close", ProducerObserverEvents{Close: struct{}{}})
}

// Dump Producer
func (p *producerImpl) Dump() (ProducerDump, error) {
	p.channel.Request(Request.MethodPRODUCER_DUMP,
		nil, p.producerId,
	)
	return ProducerDump{}, nil
}

// Get Producer stats
func (p *producerImpl) GetStats() ([]ProducerStat, error) {
	p.channel.Request(Request.MethodPRODUCER_GET_STATS,
		nil, p.producerId)
	return []ProducerStat{}, nil
}

// Pause the Producer
func (p *producerImpl) Pause() error {
	_, err := p.channel.Request(Request.MethodPRODUCER_PAUSE, nil, p.producerId)
	if err != nil {
		return err
	}
	wasPaused := p.paused
	p.paused = true
	if !wasPaused {
		p.observer.Emit("pause", ProducerObserverEvents{Pause: struct{}{}})
	}
	return nil
}

// Resume the Producer
func (p *producerImpl) Resume() error {
	_, err := p.channel.Request(Request.MethodPRODUCER_RESUME, nil, p.producerId)
	if err != nil {
		return err
	}
	wasPaused := p.paused
	p.paused = false
	if wasPaused {
		p.observer.Emit("resume", ProducerObserverEvents{Resume: struct{}{}})
	}
	return nil
}

// Enable 'trace' event
func (p *producerImpl) EnableTraceEvent(types []ProducerTraceEventType) error {

	var events []producer.TraceEventType

	for _, v := range types {
		events = append(events, producer.EnumValuesTraceEventType[string(v)])
	}
	p.channel.Request(Request.MethodPRODUCER_ENABLE_TRACE_EVENT,
		&Request.BodyT{
			Type:  Request.BodyProducer_EnableTraceEventRequest,
			Value: &producer.EnableTraceEventRequestT{Events: events},
		}, p.producerId,
	)
	return nil
}

// Send RTP packet (just valid for Producers created on a DirectTransport)
func (p *producerImpl) Send(rtpPacket []byte) {
	p.channel.Notify(Notification.EventPRODUCER_SEND,
		&Notification.BodyT{Type: Notification.BodyDataProducer_SendNotification, Value: &producer.SendNotificationT{Data: rtpPacket}}, p.producerId)
}

func (p *producerImpl) handleWorkerNotifications() {
	p.channel.On(events.EventName(p.producerId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventPRODUCER_SCORE:
			value := arg.Body.Value.(*producer.ScoreNotificationT)
			var scores []*ProducerScore
			for _, v := range value.Scores {
				s := &ProducerScore{
					EncodingIdx: (v.EncodingIdx),
					SSRC:        (v.Ssrc),
					RID:         (&v.Rid),
					Score:       (v.Score),
				}
				scores = append(scores, s)
			}
			p.score = scores
			p.Emit("score", ProducerEvents{Score: events.NewEvent1(scores)})
			p.observer.Emit("score", ProducerObserverEvents{Score: events.NewEvent1(scores)})
			break
		case Notification.EventPRODUCER_VIDEO_ORIENTATION_CHANGE:
			value := arg.Body.Value.(*producer.VideoOrientationChangeNotificationT)
			n := ProducerVideoOrientation{
				Camera:   value.Camera,
				Flip:     value.Flip,
				Rotation: int(value.Rotation),
			}
			p.Emit("videoorientationchange", ProducerEvents{VideoOrientationChange: events.NewEvent1(n)})

			p.observer.Emit("videoorientationchange", ProducerObserverEvents{VideoOrientationChange: events.NewEvent1(n)})
		case Notification.EventPRODUCER_TRACE:
			value := arg.Body.Value.(*producer.TraceNotificationT)
			e := ProducerTraceEventData{
				Type:      ProducerTraceEventType(producer.EnumNamesTraceEventType[value.Type]),
				Timestamp: int(value.Timestamp),
				Direction: value.Direction.String(),
				Info:      value.Info,
			}
			p.Emit("trace", ProducerEvents{Trace: events.NewEvent1(e)})
			p.observer.Emit("trace", ProducerObserverEvents{Trace: events.NewEvent1(e)})
		default:
			slog.Error("ignoring unknown event ", slog.Any("event", arg))
		}
	})
}
func (p *producerImpl) handleListenerError() {
	p.On("listrenererror", func(arg ProducerEvents) {
		slog.Error("event listrener threw an error ")
	})
}
