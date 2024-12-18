package mediasoupgo

type Router struct{}

func (r *Router) Dump() {}

func (r *Router) Close() {}

func (r *Router) CreateWebRtcTransport() {}

func (r *Router) CreatePlainTransport() {}

func (r *Router) CreatePipeTransport() {}

func (r *Router) CreateDirectTransport() {}

func (r *Router) PipeToRouter() {}

func (r *Router) AddPipeTransportPair() {}

func (r *Router) CreateActiveSpeakerObserver() {}

func (r *Router) CreateAudioLevelObserver() {}
