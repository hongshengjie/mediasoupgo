package mediasoupgo

type RtpStreamRecvStats struct {
	BaseRtpStreamStats
	Type           string
	Jitter         float64
	PacketCount    int
	ByteCount      int
	Bitrate        float64
	BitrateByLayer BitrateByLayer
}

type RtpStreamSendStats struct {
	BaseRtpStreamStats
	Type        string
	PacketCount int
	ByteCount   int
	Bitrate     float64
}

type BaseRtpStreamStats struct {
	Timestamp            float64
	Ssrc                 int
	RtxSsrc              *int
	Rid                  *string
	Kind                 string
	MimeType             string
	PacketsLost          int
	FractionLost         float64
	PacketsDiscarded     int
	PacketsRetransmitted int
	PacketsRepaired      int
	NackCount            int
	NackPacketCount      int
	PliCount             int
	FirCount             int
	Score                float64
	RoundTripTime        *float64
	RtxPacketsDiscarded  *int
}

type BitrateByLayer map[string]float64
