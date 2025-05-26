package mediasoupgo

type RtpCapabilities struct {
	Codecs           []*RtpCodecCapability `json:"codecs"`
	HeaderExtensions []*RtpHeaderExtension `json:"headerExtensions"`
}

type MediaKind string

const (
	AudioMediaKind MediaKind = "audio"
	VideoMediaKind MediaKind = "video"
)

type RtpCodecCapability struct {
	Kind                 MediaKind       `json:"kind"`
	MimeType             string          `json:"mimeType"`
	PreferredPayloadType *byte           `json:"preferredPayloadType"`
	ClockRate            uint32          `json:"clockRate"`
	Channels             *byte           `json:"channels"`
	Parameters           map[string]any  `json:"parameters"`
	RTCPFeedback         []*RtcpFeedback `json:"rtcpFeedback"`
}

type RtpHeaderExtensionDirection string

const (
	SendRecvRtpHeaderExtensionDirection RtpHeaderExtensionDirection = "sendrecv"
	SendOnlyRtpHeaderExtensionDirection RtpHeaderExtensionDirection = "sendonly"
	RecvOnlyRtpHeaderExtensionDirection RtpHeaderExtensionDirection = "recvonly"
	InactiveRtpHeaderExtensionDirection RtpHeaderExtensionDirection = "inactive"
)

type RtpHeaderExtension struct {
	Kind             MediaKind                    `json:"kind"`
	URI              RtpHeaderExtensionUri        `json:"uri"`
	PreferredID      byte                         `json:"preferredId"`
	PreferredEncrypt *bool                        `json:"preferredEncrypt"`
	Direction        *RtpHeaderExtensionDirection `json:"direction"`
}

type RtpParameters struct {
	MID              *string                         `json:"mid"`
	Codecs           []*RtpCodecParameters           `json:"codecs"`
	HeaderExtensions []*RtpHeaderExtensionParameters `json:"headerExtensions"`
	Encodings        []*RtpEncodingParameters        `json:"encodings"`
	RTCP             *RtcpParameters                 `json:"rtcp"`
}

type RtpCodecParameters struct {
	MimeType     string          `json:"mimeType"`
	PayloadType  byte            `json:"payloadType"`
	ClockRate    uint32          `json:"clockRate"`
	Channels     *byte           `json:"channels"`
	Parameters   map[string]any  `json:"parameters"`
	RTCPFeedback []*RtcpFeedback `json:"rtcpFeedback"`
}

type RtcpFeedback struct {
	Type      string  `json:"type"`
	Parameter *string `json:"parameter"`
}

type RtpEncodingParameters struct {
	SSRC                  *uint32  `json:"ssrc"`
	RID                   *string  `json:"rid"`
	CodecPayloadType      *byte    `json:"codecPayloadType"`
	RTX                   *RTX     `json:"rtx"`
	DTX                   *bool    `json:"dtx"`
	ScalabilityMode       *string  `json:"scalabilityMode"`
	ScaleResolutionDownBy *float64 `json:"scaleResolutionDownBy"`
	MaxBitrate            *uint32  `json:"maxBitrate"`
}

type RTX struct {
	SSRC uint32 `json:"ssrc"`
}

type RtpHeaderExtensionUri string

const (
	MIDRtpHeaderExtensionUri                 RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:mid"
	RTPStreamIDRtpHeaderExtensionUri         RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id"
	RepairedRTPStreamIDRtpHeaderExtensionUri RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id"
	FrameMarking07RtpHeaderExtensionUri      RtpHeaderExtensionUri = "http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07"
	FrameMarkingRtpHeaderExtensionUri        RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:framemarking"
	SSRCAudioLevelRtpHeaderExtensionUri      RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:ssrc-audio-level"
	VideoOrientationRtpHeaderExtensionUri    RtpHeaderExtensionUri = "urn:3gpp:video-orientation"
	ToffsetRtpHeaderExtensionUri             RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:toffset"
	TransportWideCC01RtpHeaderExtensionUri   RtpHeaderExtensionUri = "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01"
	AbsSendTimeRtpHeaderExtensionUri         RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time"
	AbsCaptureTimeRtpHeaderExtensionUri      RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/abs-capture-time"
	PlayoutDelayRtpHeaderExtensionUri        RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/playout-delay"
)

type RtpHeaderExtensionParameters struct {
	URI        RtpHeaderExtensionUri
	ID         byte
	Encrypt    *bool
	Parameters map[string]any
}

type RtcpParameters struct {
	CNAME       *string
	ReducedSize *bool
}
