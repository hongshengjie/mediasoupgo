package main

import (
	"encoding/json"

	"mediasoupgo"
)

type (
	TransportId   string
	RoomId        string
	ProducerId    string
	ConsumerId    string
	ParticipantId string
)

// ServerMessage represents messages sent from server to client.
type ServerMessage struct {
	Action string `json:"action"`
}

type TransportOptions struct {
	ID             TransportId                 `json:"id"`
	DtlsParameters mediasoupgo.DtlsParameters  `json:"dtlsParameters"`
	IceCandidates  []*mediasoupgo.IceCandidate `json:"iceCandidates"`
	IceParameters  mediasoupgo.IceParameters   `json:"iceParameters"`
}

type (
	// Specific server message types
	ServerMessageInit struct {
		Action                   string                      `json:"action"`
		RoomID                   RoomId                      `json:"roomId"`
		ConsumerTransportOptions TransportOptions            `json:"consumerTransportOptions"`
		ProducerTransportOptions TransportOptions            `json:"producerTransportOptions"`
		RouterRtpCapabilities    mediasoupgo.RtpCapabilities `json:"routerRtpCapabilities"`
	}
)

type ServerMessageProducerAdded struct {
	Action        string        `json:"action"`
	ParticipantID ParticipantId `json:"participantId"`
	ProducerID    ProducerId    `json:"producerId"`
}

type ServerMessageProducerRemoved struct {
	Action        string        `json:"action"`
	ParticipantID ParticipantId `json:"participantId"`
	ProducerID    ProducerId    `json:"producerId"`
}

type ServerMessageProduced struct {
	Action string     `json:"action"`
	ID     ProducerId `json:"id"`
}

type ServerMessageConsumed struct {
	Action        string                    `json:"action"`
	ID            ConsumerId                `json:"id"`
	ProducerID    ProducerId                `json:"producerId"`
	Kind          mediasoupgo.MediaKind     `json:"kind"`
	RtpParameters mediasoupgo.RtpParameters `json:"rtpParameters"`
}

// ClientMessage represents messages sent from client to server.
type ClientMessage struct {
	Action string `json:"action"`
}

// Specific client message types
type ClientMessageInit struct {
	RtpCapabilities mediasoupgo.RtpCapabilities `json:"rtpCapabilities"`
	Action          string                      `json:"action"`
}

type ClientMessageConnectProducerTransport struct {
	DtlsParameters mediasoupgo.DtlsParameters `json:"dtlsParameters"`
	Action         string                     `json:"action"`
}

type ClientMessageProduce struct {
	Kind          mediasoupgo.MediaKind     `json:"kind"`
	RtpParameters mediasoupgo.RtpParameters `json:"rtpParameters"`
	Action        string                    `json:"action"`
}

type ClientMessageConnectConsumerTransport struct {
	DtlsParameters mediasoupgo.DtlsParameters `json:"dtlsParameters"`
	Action         string                     `json:"action"`
}

type ClientMessageConsume struct {
	ProducerID ProducerId `json:"producerId"`
	Action     string     `json:"action"`
}

type ClientMessageConsumerResume struct {
	ID     ConsumerId `json:"id"`
	Action string     `json:"action"`
}

// MarshalJSON for ServerMessage to handle tagged union
func ParseServerMessage(d []byte) error {
	var m ServerMessage
	if err := json.Unmarshal(d, &m); err != nil {
		return err
	}
	switch m.Action {
	case "Init":
		var init ServerMessageInit
		if err := json.Unmarshal(d, &init); err != nil {
			return err
		}
	case "ProducerAdded":
		var pa ServerMessageProducerAdded
		if err := json.Unmarshal(d, &pa); err != nil {
			return err
		}
	case "ProducerRemoved":
		var pr ServerMessageProducerRemoved
		if err := json.Unmarshal(d, &pr); err != nil {
			return err
		}
	case "ConnectedProducerTransport":
	case "Produced":
		var p ServerMessageProduced
		if err := json.Unmarshal(d, &p); err != nil {
			return err
		}
	case "ConnectedConsumerTransport":
	case "Consumed":
		var c ServerMessageConsumed
		if err := json.Unmarshal(d, &c); err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON for ClientMessage to handle tagged union
func ParseClientMessage(d []byte) error {
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
	case "ConnectProducerTransport":
		var cpt ClientMessageConnectProducerTransport
		if err := json.Unmarshal(d, &cpt); err != nil {
			return err
		}
	case "Produce":
		var p ClientMessageProduce
		if err := json.Unmarshal(d, &p); err != nil {
			return err
		}
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

	return nil
}
