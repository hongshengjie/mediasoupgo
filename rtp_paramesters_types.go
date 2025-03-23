package mediasoupgo

// RtpCapabilities defines what mediasoup or an endpoint can receive at media level.
type RtpCapabilities struct {
	// Supported media and RTX codecs. Optional.
	Codecs []RtpCodecCapability
	// Supported RTP header extensions. Optional.
	HeaderExtensions []RtpHeaderExtension
}

// MediaKind represents media kind ('audio' or 'video').
type MediaKind string

const (
	MediaKindAudio MediaKind = "audio"
	MediaKindVideo MediaKind = "video"
)

// RtpCodecCapability provides information on the capabilities of a codec within the RTP capabilities.
// The list of media codecs supported by mediasoup and their settings is defined in supportedRtpCapabilities.
// Exactly one RtpCodecCapability will be present for each supported combination of parameters that requires
// a distinct value of preferredPayloadType.
type RtpCodecCapability struct {
	// Media kind.
	Kind MediaKind
	// The codec MIME media type/subtype (e.g. 'audio/opus', 'video/VP8').
	MimeType string
	// The preferred RTP payload type. Optional. If given, must be in the 96-127 range.
	PreferredPayloadType *int
	// Codec clock rate expressed in Hertz.
	ClockRate int
	// The number of channels supported (e.g. two for stereo). Just for audio. Default 1. Optional.
	Channels *int
	// Codec specific parameters. Some parameters (such as 'packetization-mode' and 'profile-level-id' in H264
	// or 'profile-id' in VP9) are critical for codec matching. Optional.
	Parameters map[string]interface{}
	// Transport layer and codec-specific feedback messages for this codec. Optional.
	RtcpFeedback []RtcpFeedback
}

// RtpHeaderExtensionDirection represents the direction of RTP header extension.
type RtpHeaderExtensionDirection string

const (
	RtpHeaderExtensionDirectionSendRecv RtpHeaderExtensionDirection = "sendrecv"
	RtpHeaderExtensionDirectionSendOnly RtpHeaderExtensionDirection = "sendonly"
	RtpHeaderExtensionDirectionRecvOnly RtpHeaderExtensionDirection = "recvonly"
	RtpHeaderExtensionDirectionInactive RtpHeaderExtensionDirection = "inactive"
)

// RtpHeaderExtension provides information relating to supported header extensions.
// The list of RTP header extensions supported by mediasoup is defined in supportedRtpCapabilities.
// mediasoup does not currently support encrypted RTP header extensions.
type RtpHeaderExtension struct {
	// Media kind.
	Kind MediaKind
	// The URI of the RTP header extension, as defined in RFC 5285.
	Uri RtpHeaderExtensionUri
	// The preferred numeric identifier that goes in the RTP packet. Must be unique.
	PreferredId int
	// If true, it is preferred that the value in the header be encrypted as per RFC 6904. Default false. Optional.
	PreferredEncrypt *bool
	// If 'sendrecv', mediasoup supports sending and receiving this RTP extension.
	// 'sendonly' means that mediasoup can send (but not receive) it. 'recvonly' means that mediasoup can receive
	// (but not send) it. Optional.
	Direction RtpHeaderExtensionDirection
}

// RtpParameters describes RTP send and receive parameters.
// The RTP send parameters describe a media stream received by mediasoup from an endpoint.
// The RTP receive parameters describe a media stream as sent by mediasoup to an endpoint.
type RtpParameters struct {
	// The MID RTP extension value as defined in the BUNDLE specification. Optional.
	Mid *string
	// Media and RTX codecs in use.
	Codecs []RtpCodecParameters
	// RTP header extensions in use. Optional.
	HeaderExtensions []RtpHeaderExtensionParameters
	// Transmitted RTP streams and their settings. Optional.
	Encodings []RtpEncodingParameters
	// Parameters used for RTCP. Optional.
	Rtcp *RtcpParameters
}

// RtpCodecParameters provides information on codec settings within the RTP parameters.
// The list of media codecs supported by mediasoup and their settings is defined in supportedRtpCapabilities.
type RtpCodecParameters struct {
	// The codec MIME media type/subtype (e.g. 'audio/opus', 'video/VP8').
	MimeType string
	// The value that goes in the RTP Payload Type Field. Must be unique.
	PayloadType int
	// Codec clock rate expressed in Hertz.
	ClockRate int
	// The number of channels supported (e.g. two for stereo). Just for audio. Default 1. Optional.
	Channels *int
	// Codec-specific parameters available for signaling. Some parameters are critical for codec matching. Optional.
	Parameters map[string]interface{}
	// Transport layer and codec-specific feedback messages for this codec. Optional.
	RtcpFeedback []RtcpFeedback
}

// RtcpFeedback provides information on RTCP feedback messages for a specific codec.
// The list of RTCP feedbacks supported by mediasoup is defined in supportedRtpCapabilities.
type RtcpFeedback struct {
	// RTCP feedback type.
	Type string
	// RTCP feedback parameter. Optional.
	Parameter *string
}

// Rtx represents RTX stream information.
type Rtx struct {
	// The RTX SSRC.
	Ssrc uint32
}

// RtpEncodingParameters provides information relating to an encoding, which represents a media RTP stream
// and its associated RTX stream (if any).
type RtpEncodingParameters struct {
	// The media SSRC. Optional.
	Ssrc *uint32
	// The RID RTP extension value. Must be unique. Optional.
	Rid *string
	// Codec payload type this encoding affects. If unset, first media codec is chosen. Optional.
	CodecPayloadType *int
	// RTX stream information. It must contain a numeric ssrc field indicating the RTX SSRC. Optional.
	Rtx *Rtx
	// It indicates whether discontinuous RTP transmission will be used. Default false. Optional.
	Dtx *bool
	// Number of spatial and temporal layers in the RTP stream (e.g. 'L1T3'). See webrtc-svc. Optional.
	ScalabilityMode *string
	// Others. Optional.
	ScaleResolutionDownBy *float64
	MaxBitrate            *uint32
}

// RtpHeaderExtensionUri represents the URI of RTP header extensions.
type RtpHeaderExtensionUri string

const (
	RtpHeaderExtensionUriSdesMid                   RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:mid"
	RtpHeaderExtensionUriSdesRtpStreamId           RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id"
	RtpHeaderExtensionUriSdesRepairedRtpStreamId   RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id"
	RtpHeaderExtensionUriFrameMarkingDraft         RtpHeaderExtensionUri = "http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07"
	RtpHeaderExtensionUriFrameMarking              RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:framemarking"
	RtpHeaderExtensionUriSsrcAudioLevel            RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:ssrc-audio-level"
	RtpHeaderExtensionUriVideoOrientation          RtpHeaderExtensionUri = "urn:3gpp:video-orientation"
	RtpHeaderExtensionUriToffset                   RtpHeaderExtensionUri = "urn:ietf:params:rtp-hdrext:toffset"
	RtpHeaderExtensionUriTransportWideCcExtensions RtpHeaderExtensionUri = "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01"
	RtpHeaderExtensionUriAbsSendTime               RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time"
	RtpHeaderExtensionUriAbsCaptureTime            RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/abs-capture-time"
	RtpHeaderExtensionUriPlayoutDelay              RtpHeaderExtensionUri = "http://www.webrtc.org/experiments/rtp-hdrext/playout-delay"
)

// RtpHeaderExtensionParameters defines a RTP header extension within the RTP parameters.
// The list of RTP header extensions supported by mediasoup is defined in supportedRtpCapabilities.
// mediasoup does not currently support encrypted RTP header extensions and no parameters are currently considered.
type RtpHeaderExtensionParameters struct {
	// The URI of the RTP header extension, as defined in RFC 5285.
	Uri RtpHeaderExtensionUri
	// The numeric identifier that goes in the RTP packet. Must be unique.
	Id int
	// If true, the value in the header is encrypted as per RFC 6904. Default false. Optional.
	Encrypt *bool
	// Configuration parameters for the header extension. Optional.
	Parameters map[string]interface{}
}

// RtcpParameters provides information on RTCP settings within the RTP parameters.
// If no cname is given in a producer's RTP parameters, the mediasoup transport will choose a random one.
type RtcpParameters struct {
	// The Canonical Name (CNAME) used by RTCP (e.g. in SDES messages). Optional.
	Cname *string
	// Whether reduced size RTCP RFC 5506 is configured (if true) or compound RTCP as specified in RFC 3550 (if false).
	// Default true. Optional.
	ReducedSize *bool
}
