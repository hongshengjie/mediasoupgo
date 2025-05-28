package mediasoupgo

import (
	"bytes"
	"log/slog"
	"strings"

	consumer "mediasoupgo/internal/FBS/Consumer"
	"mediasoupgo/internal/FBS/Notification"
	"mediasoupgo/internal/FBS/Request"
	transport "mediasoupgo/internal/FBS/Transport"
	"mediasoupgo/internal/events"
)

var _ Consumer = &consumerImpl{}

type ConsumerData struct {
	producerId    string
	kind          MediaKind
	rtpParameters RtpParameters
	typ           ConsumerType
}

type ConsumerInternal struct {
	TransportInternal
	consumerId string
}

type consumerImpl struct {
	ConsumerInternal
	data            *ConsumerData
	channel         *Channel
	closed          bool
	paused          bool
	producerPaused  bool
	priority        byte
	score           *ConsumerScore
	preferredLayers *ConsumerLayers
	currentLayers   *ConsumerLayers
	appData         ConsumerAppData
	observer        events.EventEmmiter[ConsumerObserverEvents]
	events.EventEmmiter[ConsumerEvents]
}

func NewConsumer(internal ConsumerInternal, data *ConsumerData, channel *Channel, appData ConsumerAppData, paused, producerPaused bool, score *ConsumerScore, preferredLayers *ConsumerLayers) Consumer {
	c := &consumerImpl{
		ConsumerInternal: internal,
		data:             data,
		channel:          channel,
		closed:           false,
		paused:           paused,
		producerPaused:   producerPaused,
		priority:         1,
		score:            score,
		preferredLayers:  preferredLayers,
		currentLayers:    &ConsumerLayers{},
		appData:          appData,
		observer:         events.New[ConsumerObserverEvents](),
		EventEmmiter:     events.New[ConsumerEvents](),
	}
	c.handleWorkerNotifications()
	c.handleListenerError()
	return c
}

// Consumer id
func (c *consumerImpl) ID() string {
	return c.consumerId
}

// Associated Producer id
func (c *consumerImpl) ProducerID() string {
	return c.data.producerId
}

// Whether the Consumer is closed
func (c *consumerImpl) Closed() bool {
	return c.closed
}

// Media kind
func (c *consumerImpl) Kind() MediaKind {
	return c.data.kind
}

// RTP parameters
func (c *consumerImpl) RTPParameters() RtpParameters {
	return c.data.rtpParameters
}

// Consumer type
func (c *consumerImpl) Type() ConsumerType {
	return c.data.typ
}

// Whether the Consumer is paused
func (c *consumerImpl) Paused() bool {
	return c.paused
}

// Whether the associate Producer is paused
func (c *consumerImpl) ProducerPaused() bool {
	return c.producerPaused
}

// Current priority
func (c *consumerImpl) Priority() byte {
	return c.priority
}

// Consumer score
func (c *consumerImpl) Score() ConsumerScore {
	return *c.score
}

// Preferred video layers
func (c *consumerImpl) PreferredLayers() *ConsumerLayers {
	return c.preferredLayers
}

// Current video layers
func (c *consumerImpl) CurrentLayers() *ConsumerLayers {
	return c.currentLayers
}

// App custom data
func (c *consumerImpl) AppData() AppData {
	return AppData(c.appData)
}

func (c *consumerImpl) SetAppData(appData AppData) {
	c.appData = ConsumerAppData(appData)
}

// Observer
func (c *consumerImpl) Observer() ConsumerObserver {
	return c.observer
}

// Close the Consumer
func (c *consumerImpl) Close() {
	if c.closed {
		return
	}
	c.closed = true
	c.channel.RemoveAllListeners(events.EventName(c.consumerId))

	c.channel.Request(Request.MethodTRANSPORT_CLOSE_CONSUMER,
		&Request.BodyT{Type: Request.BodyTransport_CloseConsumerRequest, Value: &transport.CloseConsumerRequestT{ConsumerId: c.consumerId}}, c.consumerId)
	c.Emit("@close", ConsumerEvents{AtClose: struct{}{}})
	c.observer.Emit("close", ConsumerObserverEvents{Close: struct{}{}})
}

// Transport was closed
func (c *consumerImpl) TransportClosed() {
	if c.closed {
		return
	}
	c.closed = true
	c.channel.RemoveAllListeners(events.EventName(c.consumerId))

	c.Emit("transportclose", ConsumerEvents{TransportClose: struct{}{}})
	c.observer.Emit("close", ConsumerObserverEvents{Close: struct{}{}})
}

// Dump Consumer
func (c *consumerImpl) Dump() (ConsumerDump, error) {
	c.channel.Request(Request.MethodCONSUMER_DUMP, nil, c.consumerId)
	return SimpleConsumerDump{}, nil
}

// Get Consumer stats
func (c *consumerImpl) GetStats() ([]interface{}, error) {
	c.channel.Request(Request.MethodCONSUMER_GET_STATS, nil, c.consumerId)
	return nil, nil
}

// Pause the Consumer
func (c *consumerImpl) Pause() error {
	_, err := c.channel.Request(Request.MethodCONSUMER_PAUSE, nil, c.consumerId)
	if err != nil {
		return err
	}
	wasPaused := c.paused
	c.paused = true
	if !wasPaused && !c.producerPaused {
		c.observer.Emit("pause", ConsumerObserverEvents{Pause: struct{}{}})
	}
	return nil
}

// Resume the Consumer
func (c *consumerImpl) Resume() error {
	_, err := c.channel.Request(Request.MethodCONSUMER_RESUME, nil, c.consumerId)
	if err != nil {
		return err
	}
	wasPaused := c.paused
	c.paused = false
	if wasPaused && !c.producerPaused {
		c.observer.Emit("resume", ConsumerObserverEvents{Resume: struct{}{}})
	}
	return nil
}

// Set preferred video layers
func (c *consumerImpl) SetPreferredLayers(layers ConsumerLayers) error {
	var temporallayer *byte
	if layers.TemporalLayer != nil {
		temporallayer = layers.TemporalLayer
	}
	resp, err := c.channel.Request(Request.MethodCONSUMER_SET_PREFERRED_LAYERS, &Request.BodyT{
		Type: Request.BodyConsumer_SetPreferredLayersRequest,

		Value: &consumer.ConsumerLayersT{
			SpatialLayer:  byte(layers.SpatialLayer),
			TemporalLayer: temporallayer,
		},
	}, c.consumerId)
	if err != nil {
		return err
	}
	resp2 := resp.Body.Value.(*consumer.SetPreferredLayersResponseT)
	t2 := resp2.PreferredLayers.TemporalLayer
	preferredLayers := &ConsumerLayers{
		TemporalLayer: t2,
		SpatialLayer:  resp2.PreferredLayers.SpatialLayer,
	}
	c.preferredLayers = preferredLayers
	return nil
}

// Set priority
func (c *consumerImpl) SetPriority(priority byte) error {
	resp, err := c.channel.Request(Request.MethodCONSUMER_SET_PRIORITY, &Request.BodyT{
		Type: Request.BodyConsumer_SetPriorityRequest,
		Value: &consumer.SetPriorityRequestT{
			Priority: byte(priority),
		},
	}, c.consumerId)
	if err != nil {
		return err
	}
	resp2 := resp.Body.Value.(*consumer.SetPriorityResponseT)
	c.priority = resp2.Priority
	return nil
}

// Unset priority
func (c *consumerImpl) UnsetPriority() error {
	return c.SetPriority(1)
}

// Request a key frame to the Producer
func (c *consumerImpl) RequestKeyFrame() error {
	_, err := c.channel.Request(Request.MethodCONSUMER_REQUEST_KEY_FRAME, nil, c.consumerId)
	return err
}

// Enable 'trace' event
func (c *consumerImpl) EnableTraceEvent(types []ConsumerTraceEventType) error {
	var events []consumer.TraceEventType
	for _, v := range types {
		events = append(events, consumer.EnumValuesTraceEventType[string(v)])
	}
	_, err := c.channel.Request(Request.MethodCONSUMER_ENABLE_TRACE_EVENT,
		&Request.BodyT{
			Type:  Request.BodyConsumer_EnableTraceEventRequest,
			Value: consumer.EnableTraceEventRequestT{Events: events},
		}, c.consumerId,
	)
	return err
}

func (p *consumerImpl) handleWorkerNotifications() {
	p.channel.On(events.EventName(p.consumerId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventCONSUMER_PRODUCER_CLOSE:
			if p.closed {
				break
			}
			p.closed = true
			p.channel.RemoveAllListeners(events.EventName(p.consumerId))
			p.Emit("@producerclose", ConsumerEvents{AtProducerClose: struct{}{}})
			p.Emit("producerclose", ConsumerEvents{ProducerClose: struct{}{}})
			p.observer.Emit("close", ConsumerObserverEvents{Close: struct{}{}})
		case Notification.EventCONSUMER_PRODUCER_PAUSE:
			if p.producerPaused {
				break
			}
			p.producerPaused = true
			p.Emit("producerpause", ConsumerEvents{ProducerPause: struct{}{}})
			if p.paused {
				p.observer.Emit("pause", ConsumerObserverEvents{Pause: struct{}{}})
			}
		case Notification.EventCONSUMER_PRODUCER_RESUME:

			if !p.producerPaused {
				break
			}
			p.producerPaused = false
			p.Emit("producerresume", ConsumerEvents{ProducerResume: struct{}{}})
			if !p.paused {
				p.observer.Emit("resume", ConsumerObserverEvents{Resume: struct{}{}})
			}
		case Notification.EventCONSUMER_SCORE:
			value := arg.Body.Value.(*consumer.ScoreNotificationT)
			s := &ConsumerScore{
				Score:          value.Score.Score,
				ProducerScore:  value.Score.ProducerScore,
				ProducerScores: value.Score.ProducerScores,
			}
			p.score = s
			p.Emit("score", ConsumerEvents{Score: events.NewEvent1(*s)})
			p.observer.Emit("score", ConsumerObserverEvents{Score: events.NewEvent1(*s)})
		case Notification.EventCONSUMER_LAYERS_CHANGE:
			value := arg.Body.Value.(*consumer.LayersChangeNotificationT)
			tl := value.Layers.TemporalLayer
			layer := &ConsumerLayers{
				SpatialLayer:  value.Layers.SpatialLayer,
				TemporalLayer: tl,
			}
			p.currentLayers = layer
			p.Emit("layerschange", ConsumerEvents{LayersChange: events.NewEvent1(*layer)})
			p.observer.Emit("layerschange", ConsumerObserverEvents{LayersChange: events.NewEvent1(*layer)})

		case Notification.EventCONSUMER_TRACE:
			value := arg.Body.Value.(*consumer.TraceNotificationT)
			trace := &ConsumerTraceEventData{
				Type:      ConsumerTraceEventType(strings.ToLower(value.Type.String())),
				Timestamp: int(value.Timestamp),
				Direction: value.Direction.String(),
				Info:      value.Info,
			}
			p.Emit("trace", ConsumerEvents{Trace: events.NewEvent1(*trace)})
			p.observer.Emit("trace", ConsumerObserverEvents{Trace: events.NewEvent1(*trace)})

		case Notification.EventCONSUMER_RTP:
			if p.closed {
				break
			}
			value := arg.Body.Value.(*consumer.RtpNotificationT)

			p.Emit("rtp", ConsumerEvents{RTP: *bytes.NewBuffer(value.Data)})

		default:
			slog.Error("ignoring unknown event ", slog.Any("event", arg))
		}
	})
}

func (p *consumerImpl) handleListenerError() {
	p.On("listenererror", func(arg ConsumerEvents) {
		slog.Error("error")
	})
}
