package mediasoupgo

type MediaKind string

const (
	Audio MediaKind = "audio"
	Vidoe MediaKind = "video"
)

type RtpCapabilities struct {
	Codes            []*RtpCodecCapability
	HeaderExtensions []*RtpHeaderExtension
}

type RtpCodecCapability struct {
	Kind                 MediaKind
	MimeType             string
	PreferredPayloadType int
	ClockRate            int
	Channels             int
	Parameters           any
	RtcpFeedback         []*RtcpFeedback
}

type RtpHeaderExtensionDirection string

const (
	SendRecv RtpHeaderExtensionDirection = "sendrecv"
	SendOnly RtpHeaderExtensionDirection = "sendonly"
	RecvOnly RtpHeaderExtensionDirection = "recvonly"
	Inactive RtpHeaderExtensionDirection = "inacitve"
)

type RtpHeaderExtension struct {
	Kind             MediaKind
	Uri              RtpHeaderExtensionUri
	PreferredId      int
	PreferredEncrypt bool
	Direciton        RtpHeaderExtensionDirection
}

type RtpParameters struct {
	Mid              string
	Codecs           []*RtpCodecParameters
	HeaderExtensions []*RtpHeaderExtensionParameters
	Encodings        []*RtpEncodingParameters
	Rtcp             *RtcpParameters
}

type RtpCodecParameters struct {
	MimeType     string
	PayloadType  int
	ClockRate    int
	Channels     int
	Parameters   any
	RtcpFeedback []RtcpFeedback
}

type RtcpFeedback struct {
	Type      string
	Parameter string
}

type RtpEncodingParameters struct {
	Ssrc                  int
	Rid                   string
	CodecPaylodtype       int
	Rtx                   struct{ Ssrc int }
	Dtx                   bool
	ScalabilityMode       string
	ScaleResolutionDownBy int
	MaxBitrate            int
}

type RtpHeaderExtensionUri string

const (
	ExtensionUriMid                    RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:mid"
	ExtensionUriRtpStreamId            RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id"
	ExtensionUriRepairedRtpStreamId    RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id"
	ExtensionUriFrameMarkingDraft07    RtpHeaderExtensionUri = "http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07"
	ExtensionUriFrameMarking           RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:framemarking"
	ExtensionUriAudioLevel             RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:ssrc-audio-level"
	ExtensionUriVedioORientation       RtpHeaderExtensionUri = "urn:3gpp:video-orientation"
	ExtensionUriTimeOffset             RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:toffset"
	ExtensionUriTransportWideCcDraft01 RtpHeaderExtensionUri = "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01"
	ExtensionUriAbsSendTime            RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time"
	ExtensionUriAbsCaptureTime         RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/abs-capture-time"
	ExtensionUriPlayoutDelay           RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/playout-delay"
)

type RtpHeaderExtensionParameters struct {
	URI        RtpHeaderExtensionUri
	Id         int
	Encrypt    bool
	Parameters any
}

type RtcpParameters struct {
	Cname       string
	ReducedSize bool
}
