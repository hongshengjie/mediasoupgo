package mediasoupgo

// BaseRtpStreamStats represents base statistics for an RTP stream.
type BaseRtpStreamStats struct {
	// Timestamp of the stats.
	Timestamp int64
	// Synchronization source identifier.
	Ssrc uint32
	// RTX synchronization source identifier. Optional.
	RtxSsrc *uint32
	// RTP stream identifier (RID). Optional.
	Rid *string
	// Media kind ("audio" or "video").
	Kind string
	// MIME type of the stream (e.g., "audio/opus").
	MimeType string
	// Number of packets lost.
	PacketsLost int
	// Fraction of packets lost (0 to 1).
	FractionLost float64
	// Number of packets discarded.
	PacketsDiscarded int
	// Number of packets retransmitted.
	PacketsRetransmitted int
	// Number of packets repaired.
	PacketsRepaired int
	// Number of NACK messages sent/received.
	NackCount int
	// Number of packets NACKed.
	NackPacketCount int
	// Number of PLI (Picture Loss Indication) messages.
	PliCount int
	// Number of FIR (Full Intra Request) messages.
	FirCount int
	// Quality score of the stream (0 to 10).
	Score int
	// Round trip time in milliseconds. Optional.
	RoundTripTime *float64
	// Number of RTX packets discarded. Optional.
	RtxPacketsDiscarded *int
}

// RtpStreamRecvStats represents statistics for a received RTP stream.
type RtpStreamRecvStats struct {
	BaseRtpStreamStats
	// Type of the stats (e.g., "inbound-rt specialisedp").
	Type string
	// Jitter in milliseconds.
	Jitter float64
	// Total number of packets received.
	PacketCount int
	// Total number of bytes received.
	ByteCount int
	// Bitrate in bits per second.
	Bitrate float64
	// Bitrate by layer for SVC streams.
	BitrateByLayer map[string]float64
}

// RtpStreamSendStats represents statistics for a sent RTP stream.
type RtpStreamSendStats struct {
	BaseRtpStreamStats
	// Type of the stats (e.g., "outbound-rtp").
	Type string
	// Total number of packets sent.
	PacketCount int
	// Total number of bytes sent.
	ByteCount int
	// Bitrate in bits per second.
	Bitrate float64
}
type BitrateByLayer map[string]float64
