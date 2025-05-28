package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"

	"mediasoupgo"
	"mediasoupgo/internal/ptr"
	"mediasoupgo/internal/smap"
)

var (
	roomHub = cmap.New[*Room]()
	connHub = cmap.New[*ParticipantConnection]()
)

type Room struct {
	RoomId       string
	Router       mediasoupgo.Router
	webrtcserver mediasoupgo.WebRtcServer
	clients      *smap.Map[string, []mediasoupgo.Producer]

	fnlock          sync.RWMutex
	producerAdded   []func(p string, pp mediasoupgo.Producer)
	producerRemoved []func(p string, pp mediasoupgo.Producer)
}

func (r *Room) AddProducer(pid string, p mediasoupgo.Producer) {
	pp, ok := r.clients.Get(pid)
	if ok {
		pp = append(pp, p)
		r.clients.Set(pid, pp)
	} else {
		var list []mediasoupgo.Producer
		list = append(list, p)
		r.clients.Set(pid, list)
	}

	r.fnlock.RLock()
	fns := r.producerAdded
	r.fnlock.RUnlock()
	for _, f := range fns {
		f(pid, p)
	}
}

func (r *Room) RemovePaticipant(pid string) {
	pp, ok := r.clients.Get(pid)
	if ok {
		r.fnlock.RLock()
		fns := r.producerRemoved
		r.fnlock.RUnlock()
		for _, p := range pp {
			for _, f := range fns {
				f(pid, p)
			}
		}
	}
	r.clients.Delete(pid)
}

type PidProducerId struct {
	pid        string
	producerId string
}

func (r *Room) GetAllProducer() []*PidProducerId {
	var list []*PidProducerId
	r.clients.Range(func(key string, value []mediasoupgo.Producer) bool {
		for _, v := range value {
			list = append(list, &PidProducerId{pid: key, producerId: v.ID()})
		}
		return true
	})
	return list
}

// /
type ParticipantConnection struct {
	ParticipantId      string
	conn               *websocket.Conn
	room               *Room
	produceTransport   mediasoupgo.WebRtcTransport
	consumeTransport   mediasoupgo.WebRtcTransport
	clientCapabilities *mediasoupgo.RtpCapabilities

	producers *smap.Map[string, mediasoupgo.Producer]
	consumers *smap.Map[string, mediasoupgo.Consumer]
}

func (p *ParticipantConnection) write(d any) error {
	dd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return p.conn.Write(context.Background(), websocket.MessageText, dd)
}

func (p *ParticipantConnection) Close() {
	p.room.RemovePaticipant(p.ParticipantId)
}

func (p *ParticipantConnection) read() error {
	var err error
	defer func() {
		p.Close()
		slog.Error("end websocket read ", "error", err)
	}()
	for {
		var t websocket.MessageType
		var d []byte

		t, d, err = p.conn.Read(context.TODO())
		if err != nil {
			return err
		}
		fmt.Println(t, string(d))

		var m ClientMessage
		if err := json.Unmarshal(d, &m); err != nil {
			return err
		}
		switch m.Action {
		case "Init":
			var init ClientMessageInit
			if err := json.Unmarshal(d, &init); err != nil {
				return err
			}
			p.clientCapabilities = &init.RtpCapabilities
		case "ConnectProducerTransport":
			var cpt ClientMessageConnectProducerTransport
			if err := json.Unmarshal(d, &cpt); err != nil {
				return err
			}
			if err := p.produceTransport.Connect(cpt.DtlsParameters); err != nil {
				slog.Error("connect ", "error", err)
			}
			p.write(&ServerMessage{Action: "ConnectedProducerTransport"})
		case "Produce":
			var pp ClientMessageProduce
			if err := json.Unmarshal(d, &pp); err != nil {
				return err
			}
			ppp, err := p.produceTransport.Produce(&mediasoupgo.ProducerOptions{
				Kind:          pp.Kind,
				RTPParameters: pp.RtpParameters,
			})
			if err != nil {
				slog.Error("produce ", "error", err)
				continue
			}
			p.producers.Set(ppp.ID(), ppp)
			p.write(&ServerMessageProduced{
				Action: "Produced",
				ID:     ProducerId(ppp.ID()),
			})
			p.room.AddProducer(p.ParticipantId, ppp)
		case "ConnectConsumerTransport":
			var cct ClientMessageConnectConsumerTransport
			if err := json.Unmarshal(d, &cct); err != nil {
				return err
			}
			if err := p.consumeTransport.Connect(cct.DtlsParameters); err != nil {
				slog.Error("consumeTransport connect error", "error", err)
				continue
			}
			p.write(&ServerMessage{Action: "ConnectedConsumerTransport"})
		case "Consume":
			var c ClientMessageConsume
			if err := json.Unmarshal(d, &c); err != nil {
				return err
			}
			consumer, err := p.consumeTransport.Consume(&mediasoupgo.ConsumerOptions{
				ProducerID:      string(c.ProducerID),
				RTPCapabilities: *p.clientCapabilities,
				Paused:          ptr.Bool(false),
			})
			if err != nil {
				slog.Error("consume error", "error", err)
				continue
			}
			p.consumers.Set(consumer.ID(), consumer)
			p.write(&ServerMessageConsumed{
				Action:        "Consumed",
				ID:            ConsumerId(consumer.ID()),
				ProducerID:    c.ProducerID,
				Kind:          consumer.Kind(),
				RtpParameters: consumer.RTPParameters(),
			})
		case "ConsumerResume":
			var cr ClientMessageConsumerResume
			if err := json.Unmarshal(d, &cr); err != nil {
				return err
			}
			c, ok := p.consumers.Get(string(cr.ID))

			if ok {
				c.Resume()
			}
		}

	}
}

func main() {
	worker := mediasoupgo.NewWorker(&mediasoupgo.WorkerSettings{
		LogLevel: ptr.To(mediasoupgo.Debug),
		LogTags:  mediasoupgo.DefaultLogTags(),
	})
	defer worker.Close()

	server, err := worker.CreateWebRtcServer(
		&mediasoupgo.WebRtcServerOption{
			ListenInfos: []*mediasoupgo.TransportListenInfo{
				{
					Protocol:         mediasoupgo.UDPTransportProtocol,
					IP:               "127.0.0.1",
					AnnouncedAddress: ptr.String("127.0.0.1"),
					Port:             ptr.To[uint16](2233),
					PortRange:        nil,
					Flags:            nil,
					SendBufferSize:   ptr.To(uint32(1024 * 100)),
					RecvBufferSize:   ptr.To(uint32(1024 * 100)),
				},
			},
		})
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		roomId := r.URL.Query().Get("roomId")
		if roomId == "" {
			roomId = uuid.NewString()
		}
		c, err := websocket.Accept(
			w,
			r,
			&websocket.AcceptOptions{OriginPatterns: []string{"*"}},
		)
		if err != nil {
			panic(err)
		}
		var ok bool
		var room *Room
		if room, ok = roomHub.Get(roomId); !ok {
			router, err := worker.CreateRouter(&mediasoupgo.RouterOption{
				MediaCodecs: []*mediasoupgo.RtpCodecCapability{
					{
						PreferredPayloadType: ptr.To[byte](111),
						Kind:                 mediasoupgo.AudioMediaKind,
						MimeType:             "audio/opus",
						ClockRate:            48000,
						Channels:             ptr.To[byte](2),
						RTCPFeedback: []*mediasoupgo.RtcpFeedback{
							{Type: "nack"},
							{Type: "transport-cc"},
						},
					},

					{
						MimeType:  "video/VP8",
						ClockRate: 90000,
						RTCPFeedback: []*mediasoupgo.RtcpFeedback{
							{Type: "nack"},
							{Type: "nack", Parameter: ptr.To("pli")},
							{Type: "ccm", Parameter: ptr.To("fir")},
							{Type: "goog-remb"},
							{Type: "transport-cc"},
						},
					},
				},
				AppData: mediasoupgo.RouterAppData{},
			})
			if err != nil {
				panic(err)
			}
			room = &Room{RoomId: roomId, Router: router, webrtcserver: server, clients: smap.New[string, []mediasoupgo.Producer]()}
			roomHub.Set(roomId, room)
		}
		p := &ParticipantConnection{
			ParticipantId: uuid.NewString(),
			conn:          c,
			room:          room,
			producers:     smap.New[string, mediasoupgo.Producer](),
			consumers:     smap.New[string, mediasoupgo.Consumer](),
		}
		p.produceTransport, err = room.Router.CreateWebRtcTransport(
			&mediasoupgo.WebRtcTransportOptions{
				EnableUdp:                       ptr.To(true),
				EnableTcp:                       ptr.To(true),
				PreferUdp:                       nil,
				PreferTcp:                       nil,
				IceConsentTimeout:               nil,
				InitialAvailableOutgoingBitrate: nil,
				EnableSctp:                      ptr.To(true),
				NumSctpStreams:                  nil,
				MaxSctpMessageSize:              nil,
				SctpSendBufferSize:              nil,
				AppData:                         map[string]any{},
				WebRtcTransportListen: &mediasoupgo.WebRtcTransportListen{
					WebRtcServer: server,
				},
			})

		p.consumeTransport, err = room.Router.CreateWebRtcTransport(
			&mediasoupgo.WebRtcTransportOptions{
				EnableUdp:                       ptr.To(true),
				EnableTcp:                       ptr.To(true),
				PreferUdp:                       nil,
				PreferTcp:                       nil,
				IceConsentTimeout:               nil,
				InitialAvailableOutgoingBitrate: nil,
				EnableSctp:                      ptr.To(true),
				NumSctpStreams:                  nil,
				MaxSctpMessageSize:              nil,
				SctpSendBufferSize:              nil,
				AppData:                         map[string]any{},
				WebRtcTransportListen: &mediasoupgo.WebRtcTransportListen{
					WebRtcServer: server,
				},
			})
		p.write(&ServerMessageInit{
			Action: "Init",
			RoomID: RoomId(room.RoomId),
			ConsumerTransportOptions: TransportOptions{
				ID:             TransportId(p.consumeTransport.ID()),
				DtlsParameters: p.consumeTransport.DtlsParameters(),
				IceCandidates:  p.consumeTransport.IceCandidates(),
				IceParameters:  p.consumeTransport.IceParameters(),
			},
			ProducerTransportOptions: TransportOptions{
				ID:             TransportId(p.produceTransport.ID()),
				DtlsParameters: p.produceTransport.DtlsParameters(),
				IceCandidates:  p.produceTransport.IceCandidates(),
				IceParameters:  p.produceTransport.IceParameters(),
			},
			RouterRtpCapabilities: *p.room.Router.RtpCapabilities(),
		})
		p.room.fnlock.Lock()
		p.room.producerAdded = append(p.room.producerAdded, func(pr string, pp mediasoupgo.Producer) {
			if pr != p.ParticipantId {
				p.write(&ServerMessageProducerAdded{
					Action:        "ProducerAdded",
					ParticipantID: ParticipantId(pr),
					ProducerID:    ProducerId(pp.ID()),
				})
			}
		})

		p.room.producerRemoved = append(p.room.producerRemoved, func(pr string, pp mediasoupgo.Producer) {
			if pr != p.ParticipantId {
				p.write(&ServerMessageProducerRemoved{
					Action:        "ProducerRemoved",
					ParticipantID: ParticipantId(pr),
					ProducerID:    ProducerId(pp.ID()),
				})
			}
		})
		p.room.fnlock.Unlock()
		producers := p.room.GetAllProducer()
		for _, v := range producers {
			p.write(&ServerMessageProducerAdded{
				Action:        "ProducerAdded",
				ParticipantID: ParticipantId(v.pid),
				ProducerID:    ProducerId(v.producerId),
			})
		}
		go p.read()
		connHub.Set(p.ParticipantId, p)
	})
	http.ListenAndServe(":3000", nil)
}
