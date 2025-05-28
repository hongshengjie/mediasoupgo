package mediasoupgo

import (
	dataproducer "mediasoupgo/internal/FBS/DataProducer"
	"mediasoupgo/internal/FBS/Notification"
	"mediasoupgo/internal/FBS/Request"
	transport "mediasoupgo/internal/FBS/Transport"
	"mediasoupgo/internal/events"
)

var _ DataProducer = &dataProducerImpl{}

type DataProducerInternal struct {
	TransportInternal
	dataProducerId string
}
type DataProducerData struct {
	typ                  DataProducerType
	sctpStreamParameters *SctpStreamParameters
	lable                string
	protocol             string
}
type dataProducerImpl struct {
	DataProducerInternal
	data     *DataProducerData
	channel  *Channel
	closed   bool
	paused   bool
	appData  AppData
	observer events.EventEmmiter[DataProducerObserverEvents]
	events.EventEmmiter[DataProducerEvents]
}

func NewDataProducer(internal DataProducerInternal, data *DataProducerData, c *Channel, paused bool, appData AppData) DataProducer {
	d := &dataProducerImpl{
		DataProducerInternal: internal,
		data:                 data,
		channel:              c,
		closed:               false,
		paused:               paused,
		appData:              appData,
		observer:             events.New[DataProducerObserverEvents](),
		EventEmmiter:         events.New[DataProducerEvents](),
	}
	d.handleWorkerNotifications()
	d.handleListenerError()
	return d
}

// DataProducer id
func (d *dataProducerImpl) ID() (_ string) {
	return d.dataProducerId
}

// Whether the DataProducer is closed
func (d *dataProducerImpl) Closed() (_ bool) {
	return d.closed
}

// DataProducer type
func (d *dataProducerImpl) Type() (_ DataProducerType) {
	return d.data.typ
}

// SCTP stream parameters
func (d *dataProducerImpl) SCTPStreamParameters() (_ *SctpStreamParameters) {
	return d.data.sctpStreamParameters
}

// DataChannel label
func (d *dataProducerImpl) Label() (_ string) {
	return d.data.lable
}

// DataChannel protocol
func (d *dataProducerImpl) Protocol() (_ string) {
	return d.data.protocol
}

// Whether the DataProducer is paused
func (d *dataProducerImpl) Paused() (_ bool) {
	return d.paused
}

// App custom data
func (d *dataProducerImpl) AppData() (_ AppData) {
	return d.appData
}

func (d *dataProducerImpl) SetAppData(appData AppData) {
	d.appData = appData
}

// Observer
func (d *dataProducerImpl) Observer() (_ DataProducerObserver) {
	return d.observer
}

// Close the DataProducer
func (d *dataProducerImpl) Close() {
	if d.closed {
		return
	}
	d.closed = true
	d.channel.RemoveAllListeners(events.EventName(d.dataProducerId))
	d.channel.Request(Request.MethodTRANSPORT_CLOSE_DATAPRODUCER, &Request.BodyT{
		Type:  Request.BodyTransport_CloseDataProducerRequest,
		Value: &transport.CloseDataProducerRequestT{DataProducerId: d.dataProducerId},
	}, d.transportId)

	d.Emit("@close", DataProducerEvents{AtClose: struct{}{}})
	d.observer.Emit("close", DataProducerObserverEvents{Close: struct{}{}})
}

// Transport was closed
func (d *dataProducerImpl) TransportClosed() {
	if d.closed {
		return
	}
	d.closed = true
	d.channel.RemoveAllListeners(events.EventName(d.dataProducerId))
	d.Emit("transportclose", DataProducerEvents{TransportClose: struct{}{}})
	d.observer.Emit("close", DataProducerObserverEvents{Close: struct{}{}})
}

// Dump DataProducer
func (d *dataProducerImpl) Dump() (_ DataProducerDump, _ error) {
	d.channel.Request(Request.MethodDATAPRODUCER_DUMP, nil, d.dataProducerId)
	return DataProducerDump{}, nil
}

// Get DataProducer stats
func (d *dataProducerImpl) GetStats() (_ []DataProducerStat, _ error) {
	d.channel.Request(Request.MethodDATAPRODUCER_GET_STATS, nil, d.dataProducerId)
	return nil, nil
}

// Pause the DataProducer
func (d *dataProducerImpl) Pause() (_ error) {

	d.channel.Request(Request.MethodDATAPRODUCER_PAUSE, nil, d.dataProducerId)
	wapPaused := d.paused
	d.paused = true
	if !wapPaused {
		d.observer.Emit("pause", DataProducerObserverEvents{Pause: struct{}{}})
	}
	return nil
}

// Resume the DataProducer
func (d *dataProducerImpl) Resume() (_ error) {
	d.channel.Request(Request.MethodDATAPRODUCER_RESUME, nil, d.dataProducerId)
	wapPaused := d.paused
	d.paused = false
	if wapPaused {
		d.observer.Emit("resume", DataProducerObserverEvents{Resume: struct{}{}})
	}
	return nil
}

// Send data (just valid for DataProducers created on a DirectTransport)
func (d *dataProducerImpl) Send(message []byte, isString bool, subchannels []uint16, requiredSubchannel uint16) {
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
	d.channel.Notify(Notification.EventDATAPRODUCER_SEND, &Notification.BodyT{
		Type: Notification.BodyDataProducer_SendNotification,
		Value: &dataproducer.SendNotificationT{
			Ppid:               uint32(ppid),
			Data:               message,
			Subchannels:        subchannels,
			RequiredSubchannel: &requiredSubchannel,
		},
	}, d.dataProducerId)
}
func (d *dataProducerImpl) handleWorkerNotifications() {}
func (d *dataProducerImpl) handleListenerError()       {}
