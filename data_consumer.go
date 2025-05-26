package mediasoupgo

import (
	dataconsumer "mediasoupgo/FBS/DataConsumer"
	"mediasoupgo/FBS/Notification"
	"mediasoupgo/FBS/Request"
	transport "mediasoupgo/FBS/Transport"
	"mediasoupgo/events"
)

var _ DataConsumer = &dataConsumerImpl{}

type DataConsumerData struct {
	dataProducerId             string
	typ                        DataConsumerType
	sctpStreamParameters       *SctpStreamParameters
	lable                      string
	protocol                   string
	bufferedAmountLowThreshold uint32
}
type DataConsumerInternal struct {
	TransportInternal
	dataConsumerId string
}

type dataConsumerImpl struct {
	DataConsumerInternal
	channel            *Channel
	data               *DataConsumerData
	closed             bool
	paused             bool
	dataProducerPaused bool
	subchannels        []uint16
	appData            AppData
	observer           events.EventEmmiter[DataConsumerObserverEvents]
	events.EventEmmiter[DataConsumerEvents]
}

func NewDataComsumer(internal DataConsumerInternal, data *DataConsumerData, c *Channel, paused bool, dataProducerPaused bool, subchannels []uint16, appData AppData) DataConsumer {
	d := &dataConsumerImpl{
		DataConsumerInternal: internal,
		channel:              c,
		data:                 data,
		closed:               false,
		paused:               paused,
		dataProducerPaused:   dataProducerPaused,
		subchannels:          subchannels,
		appData:              appData,
		observer:             events.New[DataConsumerObserverEvents](),
		EventEmmiter:         events.New[DataConsumerEvents](),
	}
	d.handleListenerError()
	d.handleWorkerNotifications()
	return d
}

// DataConsumer id
func (d *dataConsumerImpl) ID() (_ string) {
	return d.dataConsumerId
}

// Associated DataProducer id
func (d *dataConsumerImpl) DataProducerID() (_ string) {
	return d.data.dataProducerId
}

// Whether the DataConsumer is closed
func (d *dataConsumerImpl) Closed() (_ bool) {
	return d.closed
}

// DataConsumer type
func (d *dataConsumerImpl) Type() (_ DataConsumerType) {
	return d.data.typ
}

// SCTP stream parameters
func (d *dataConsumerImpl) SCTPStreamParameters() (_ *SctpStreamParameters) {
	return d.data.sctpStreamParameters
}

// DataChannel label
func (d *dataConsumerImpl) Label() (_ string) {
	return d.data.lable
}

// DataChannel protocol
func (d *dataConsumerImpl) Protocol() (_ string) {
	return d.data.protocol
}

// Whether the DataConsumer is paused
func (d *dataConsumerImpl) Paused() (_ bool) {
	return d.paused
}

// Whether the associate DataProducer is paused
func (d *dataConsumerImpl) DataProducerPaused() (_ bool) {
	return d.dataProducerPaused
}

// Get current subchannels this data consumer is subscribed to
func (d *dataConsumerImpl) Subchannels() (_ []uint16) {
	return d.subchannels
}

// App custom data
func (d *dataConsumerImpl) AppData() (_ AppData) {
	return d.appData
}

func (d *dataConsumerImpl) SetAppData(appData AppData) {
	d.appData = appData
}

// Observer
func (d *dataConsumerImpl) Observer() (_ DataConsumerObserver) {
	return d.observer
}

// Close the DataConsumer
func (d *dataConsumerImpl) Close() {
	if d.closed {
		return
	}
	d.closed = true
	d.channel.RemoveAllListeners(events.EventName(d.dataConsumerId))
	d.channel.Request(Request.MethodTRANSPORT_CLOSE_DATACONSUMER,
		&Request.BodyT{
			Type:  Request.BodyTransport_CloseDataConsumerRequest,
			Value: &transport.CloseDataConsumerRequestT{DataConsumerId: d.dataConsumerId},
		}, d.transportId,
	)
	d.Emit("@close", DataConsumerEvents{AtClose: struct{}{}})
	d.observer.Emit("close", DataConsumerObserverEvents{Close: struct{}{}})
}

// Transport was closed
func (d *dataConsumerImpl) TransportClosed() {
	if d.closed {
		return
	}
	d.closed = true
	d.channel.RemoveAllListeners(events.EventName(d.dataConsumerId))
	d.Emit("transportclose", DataConsumerEvents{TransportClose: struct{}{}})
	d.observer.Emit("close", DataConsumerObserverEvents{Close: struct{}{}})
}

// Dump DataConsumer
func (d *dataConsumerImpl) Dump() (_ DataConsumerDump, _ error) {
	d.channel.Request(Request.MethodDATACONSUMER_DUMP, nil, d.dataConsumerId)
	return DataConsumerDump{}, nil
}

// Get DataConsumer stats
func (d *dataConsumerImpl) GetStats() (_ []DataConsumerStat, _ error) {
	d.channel.Request(Request.MethodDATACONSUMER_GET_STATS, nil, d.dataConsumerId)
	return nil, nil
}

// Pause the DataConsumer
func (d *dataConsumerImpl) Pause() (_ error) {
	d.channel.Request(Request.MethodDATACONSUMER_PAUSE, nil, d.dataConsumerId)
	wasPaused := d.paused
	d.paused = true
	if !wasPaused && !d.dataProducerPaused {
		d.observer.Emit("pause", DataConsumerObserverEvents{Pause: struct{}{}})
	}
	return nil
}

// Resume the DataConsumer
func (d *dataConsumerImpl) Resume() (_ error) {
	d.channel.Request(Request.MethodDATACONSUMER_RESUME, nil, d.dataConsumerId)

	wasPaused := d.paused
	d.paused = false
	if wasPaused && !d.dataProducerPaused {
		d.observer.Emit("resume", DataConsumerObserverEvents{Resume: struct{}{}})
	}
	return nil
}

// Set buffered amount low threshold
func (d *dataConsumerImpl) SetBufferedAmountLowThreshold(threshold int) (_ error) {
	d.channel.Request(Request.MethodDATACONSUMER_SET_BUFFERED_AMOUNT_LOW_THRESHOLD, &Request.BodyT{
		Type:  Request.BodyDataConsumer_SetBufferedAmountLowThresholdRequest,
		Value: &dataconsumer.SetBufferedAmountLowThresholdRequestT{Threshold: uint32(threshold)},
	}, d.dataConsumerId)
	return nil
}

// Get buffered amount size
func (d *dataConsumerImpl) GetBufferedAmount() (_ int, _ error) {
	resp, err := d.channel.Request(Request.MethodDATACONSUMER_GET_BUFFERED_AMOUNT, nil, d.dataConsumerId)
	if err != nil {
		return 0, err
	}
	resp2 := resp.Body.Value.(*dataconsumer.GetBufferedAmountResponseT)

	return int(resp2.BufferedAmount), nil
}

// Send a message
func (d *dataConsumerImpl) Send(message []byte, isString bool) (_ error) {

	var ppid int
	if isString {
		ppid = 51
		if len(message) == 0 {
			ppid = 56
			message = []byte(" ")
		}
	} else {
		ppid = 53
		if len(message) == 0 {
			ppid = 57
			message = []byte{1}
		}
	}
	d.channel.Request(Request.MethodDATACONSUMER_SEND, &Request.BodyT{
		Type: Request.BodyDataConsumer_SendRequest,
		Value: &dataconsumer.SendRequestT{
			Ppid: uint32(ppid),
			Data: message,
		},
	}, d.dataConsumerId)
	return nil
}

// Set subchannels
func (d *dataConsumerImpl) SetSubchannels(subchannels []uint16) (_ error) {
	resp, err := d.channel.Request(Request.MethodDATACONSUMER_SET_SUBCHANNELS, &Request.BodyT{
		Type: Request.BodyDataConsumer_SetSubchannelsRequest,
		Value: &dataconsumer.SetSubchannelsRequestT{
			Subchannels: subchannels,
		},
	}, d.dataConsumerId)
	if err != nil {
		return err
	}
	resp2 := resp.Body.Value.(*dataconsumer.SetSubchannelsResponseT)
	d.subchannels = resp2.Subchannels
	return nil
}

// Add a subchannel
func (d *dataConsumerImpl) AddSubchannel(subchannel int) (_ error) {
	resp, err := d.channel.Request(Request.MethodDATACONSUMER_ADD_SUBCHANNEL, &Request.BodyT{
		Type: Request.BodyDataConsumer_AddSubchannelRequest,
		Value: &dataconsumer.AddSubchannelRequestT{
			Subchannel: uint16(subchannel),
		},
	}, d.dataConsumerId)
	if err != nil {
		return err
	}
	resp2 := resp.Body.Value.(*dataconsumer.AddSubchannelResponseT)
	d.subchannels = resp2.Subchannels
	return nil
}

// Remove a subchannel
func (d *dataConsumerImpl) RemoveSubchannel(subchannel int) (_ error) {
	resp, err := d.channel.Request(Request.MethodDATACONSUMER_REMOVE_SUBCHANNEL, &Request.BodyT{
		Type: Request.BodyDataConsumer_RemoveSubchannelRequest,
		Value: &dataconsumer.RemoveSubchannelRequestT{
			Subchannel: uint16(subchannel),
		},
	}, d.dataConsumerId)
	if err != nil {
		return err
	}
	resp2 := resp.Body.Value.(*dataconsumer.RemoveSubchannelResponseT)
	d.subchannels = resp2.Subchannels
	return nil
}

func (d *dataConsumerImpl) handleWorkerNotifications() {
	d.channel.On(events.EventName(d.dataConsumerId), func(arg *Notification.NotificationT) {
		switch arg.Event {
		case Notification.EventDATACONSUMER_BUFFERED_AMOUNT_LOW:
			value := arg.Body.Value.(*dataconsumer.BufferedAmountLowNotificationT)
			d.Emit("bufferedamountlow", DataConsumerEvents{BufferedAmountLow: value.BufferedAmount})
		case Notification.EventDATACONSUMER_SCTP_SENDBUFFER_FULL:
			d.Emit("sctpsendbufferfull", DataConsumerEvents{SCTPSendBufferFull: struct{}{}})
		case Notification.EventDATACONSUMER_DATAPRODUCER_PAUSE:
			if d.dataProducerPaused {
				break
			}
			d.dataProducerPaused = true
			d.Emit("dataproducerpause", DataConsumerEvents{DataProducerPause: struct{}{}})
			if !d.paused {
				d.observer.Emit("pause", DataConsumerObserverEvents{Pause: struct{}{}})
			}
		case Notification.EventDATACONSUMER_DATAPRODUCER_RESUME:

			if !d.dataProducerPaused {
				break
			}
			d.dataProducerPaused = false
			d.Emit("dataproducerresume", DataConsumerEvents{DataProducerResume: struct{}{}})
			if !d.paused {
				d.observer.Emit("pause", DataConsumerObserverEvents{Resume: struct{}{}})
			}
		case Notification.EventDATACONSUMER_DATAPRODUCER_CLOSE:
			if d.closed {
				break
			}
			d.closed = true
			d.channel.RemoveAllListeners(events.EventName(d.dataConsumerId))
			d.Emit("@dataproducerclose", DataConsumerEvents{AtDataProducerClose: struct{}{}})
			d.Emit("dataproducerclose", DataConsumerEvents{DataProducerClose: struct{}{}})
			d.observer.Emit("close", DataConsumerObserverEvents{Close: struct{}{}})
		case Notification.EventDATACONSUMER_MESSAGE:
			if d.closed {
				break
			}
			value := arg.Body.Value.(*dataconsumer.MessageNotificationT)
			d.Emit("message", DataConsumerEvents{Message: events.NewEvent1(&MessageItem{Ppid: int(value.Ppid), Data: value.Data})})
		}
	})
}
func (d *dataConsumerImpl) handleListenerError() {}
