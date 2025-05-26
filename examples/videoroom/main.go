package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"

	"mediasoupgo"
	"mediasoupgo/ptr"
)

var (
	roomHub = cmap.New[*Room]()
	connHub = cmap.New[*ParticipantConnection]()
)

type (
	Room struct {
		RoomId string
		Router mediasoupgo.Router

		webrtcserver mediasoupgo.WebRtcServer
	}
	ParticipantConnection struct {
		ParticipantId      string
		conn               *websocket.Conn
		room               *Room
		produceTransport   mediasoupgo.WebRtcTransport
		consumeTransport   mediasoupgo.WebRtcTransport
		clientCapabilities *mediasoupgo.RtpCapabilities
		producers          map[string]mediasoupgo.Producer
	}
)

func (p *ParticipantConnection) write(d any) error {
	dd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return p.conn.Write(context.Background(), websocket.MessageText, dd)
}

func (p *ParticipantConnection) read() error {
	var err error
	defer func() {
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
			p.producers[ppp.ID()] = ppp
			p.write(&ServerMessageProduced{ID: ProducerId(ppp.ID()), Action: "Produced"})
		case "ConnectConsumerTransport":
			var cct ClientMessageConnectConsumerTransport
			if err := json.Unmarshal(d, &cct); err != nil {
				return err
			}
		case "Consume":
			var c ClientMessageConsume
			if err := json.Unmarshal(d, &c); err != nil {
				return err
			}
		case "ConsumerResume":
			var cr ClientMessageConsumerResume
			if err := json.Unmarshal(d, &cr); err != nil {
				return err
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
			room = &Room{RoomId: roomId, Router: router, webrtcserver: server}
			roomHub.Set(roomId, room)
		}
		p := &ParticipantConnection{
			ParticipantId: uuid.NewString(),
			conn:          c,
			room:          room,
			producers:     map[string]mediasoupgo.Producer{},
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
		go p.read()
		connHub.Set(p.ParticipantId, p)
	})
	http.ListenAndServe(":3000", nil)
}
