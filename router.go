package mediasoupgo

import (
	"mediasoupgo/FBS/Request"
	"mediasoupgo/FBS/Response"
)

type Router struct {
	routerId string
	channel  *Channel
}

func (r *Router) Dump() (*Response.ResponseT, error) {
	return r.channel.Request(Request.MethodROUTER_DUMP, &Request.BodyT{Type: Request.BodyNONE}, r.routerId)
}

func (r *Router) Close() {}

func (r *Router) CreateWebRtcTransport() {}

func (r *Router) CreatePlainTransport() {}

func (r *Router) CreatePipeTransport() {}

func (r *Router) CreateDirectTransport() {}

func (r *Router) PipeToRouter() {}

func (r *Router) AddPipeTransportPair() {}

func (r *Router) CreateActiveSpeakerObserver() {}

func (r *Router) CreateAudioLevelObserver() {}
