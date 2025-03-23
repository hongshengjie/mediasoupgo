package mediasoupgo

// SupportedRtpCapabilities defines the supported RTP capabilities.
var SupportedRtpCapabilities = RtpCapabilities{
	Codecs: []RtpCodecCapability{
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/opus",
			ClockRate: 48000,
			Channels:  intPtr(2),
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/multiopus",
			ClockRate: 48000,
			Channels:  intPtr(4),
			Parameters: map[string]interface{}{
				"channel_mapping": "0,1,2,3",
				"num_streams":     2,
				"coupled_streams": 2,
			},
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/multiopus",
			ClockRate: 48000,
			Channels:  intPtr(6),
			Parameters: map[string]interface{}{
				"channel_mapping": "0,4,1,2,3,5",
				"num_streams":     4,
				"coupled_streams": 2,
			},
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/multiopus",
			ClockRate: 48000,
			Channels:  intPtr(8),
			Parameters: map[string]interface{}{
				"channel_mapping": "0,6,1,2,3,4,5,7",
				"num_streams":     5,
				"coupled_streams": 3,
			},
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:                 MediaKindAudio,
			MimeType:             "audio/PCMU",
			PreferredPayloadType: intPtr(0),
			ClockRate:            8000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:                 MediaKindAudio,
			MimeType:             "audio/PCMA",
			PreferredPayloadType: intPtr(8),
			ClockRate:            8000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/ISAC",
			ClockRate: 32000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/ISAC",
			ClockRate: 16000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:                 MediaKindAudio,
			MimeType:             "audio/G722",
			PreferredPayloadType: intPtr(9),
			ClockRate:            8000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/iLBC",
			ClockRate: 8000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/SILK",
			ClockRate: 24000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/SILK",
			ClockRate: 16000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/SILK",
			ClockRate: 12000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/SILK",
			ClockRate: 8000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "transport-cc"},
			},
		},
		{
			Kind:                 MediaKindAudio,
			MimeType:             "audio/CN",
			PreferredPayloadType: intPtr(13),
			ClockRate:            32000,
		},
		{
			Kind:                 MediaKindAudio,
			MimeType:             "audio/CN",
			PreferredPayloadType: intPtr(13),
			ClockRate:            16000,
		},
		{
			Kind:                 MediaKindAudio,
			MimeType:             "audio/CN",
			PreferredPayloadType: intPtr(13),
			ClockRate:            8000,
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/telephone-event",
			ClockRate: 48000,
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/telephone-event",
			ClockRate: 32000,
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/telephone-event",
			ClockRate: 16000,
		},
		{
			Kind:      MediaKindAudio,
			MimeType:  "audio/telephone-event",
			ClockRate: 8000,
		},
		{
			Kind:      MediaKindVideo,
			MimeType:  "video/VP8",
			ClockRate: 90000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "nack", Parameter: stringPtr("pli")},
				{Type: "ccm", Parameter: stringPtr("fir")},
				{Type: "goog-remb"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindVideo,
			MimeType:  "video/VP9",
			ClockRate: 90000,
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "nack", Parameter: stringPtr("pli")},
				{Type: "ccm", Parameter: stringPtr("fir")},
				{Type: "goog-remb"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindVideo,
			MimeType:  "video/H264",
			ClockRate: 90000,
			Parameters: map[string]interface{}{
				"level-asymmetry-allowed": 1,
			},
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "nack", Parameter: stringPtr("pli")},
				{Type: "ccm", Parameter: stringPtr("fir")},
				{Type: "goog-remb"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindVideo,
			MimeType:  "video/H264-SVC",
			ClockRate: 90000,
			Parameters: map[string]interface{}{
				"level-asymmetry-allowed": 1,
			},
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "nack", Parameter: stringPtr("pli")},
				{Type: "ccm", Parameter: stringPtr("fir")},
				{Type: "goog-remb"},
				{Type: "transport-cc"},
			},
		},
		{
			Kind:      MediaKindVideo,
			MimeType:  "video/H265",
			ClockRate: 90000,
			Parameters: map[string]interface{}{
				"level-asymmetry-allowed": 1,
			},
			RtcpFeedback: []RtcpFeedback{
				{Type: "nack"},
				{Type: "nack", Parameter: stringPtr("pli")},
				{Type: "ccm", Parameter: stringPtr("fir")},
				{Type: "goog-remb"},
				{Type: "transport-cc"},
			},
		},
	},
	HeaderExtensions: []RtpHeaderExtension{
		{
			Kind:             MediaKindAudio,
			Uri:              RtpHeaderExtensionUriSdesMid,
			PreferredId:      1,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriSdesMid,
			PreferredId:      1,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriSdesRtpStreamId,
			PreferredId:      2,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionRecvOnly,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriSdesRepairedRtpStreamId,
			PreferredId:      3,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionRecvOnly,
		},
		{
			Kind:             MediaKindAudio,
			Uri:              RtpHeaderExtensionUriAbsSendTime,
			PreferredId:      4,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriAbsSendTime,
			PreferredId:      4,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindAudio,
			Uri:              RtpHeaderExtensionUriTransportWideCcExtensions,
			PreferredId:      5,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionRecvOnly,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriTransportWideCcExtensions,
			PreferredId:      5,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriFrameMarkingDraft,
			PreferredId:      6,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriFrameMarking,
			PreferredId:      7,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindAudio,
			Uri:              RtpHeaderExtensionUriSsrcAudioLevel,
			PreferredId:      10,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriVideoOrientation,
			PreferredId:      11,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriToffset,
			PreferredId:      12,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindAudio,
			Uri:              RtpHeaderExtensionUriAbsCaptureTime,
			PreferredId:      13,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriAbsCaptureTime,
			PreferredId:      13,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindAudio,
			Uri:              RtpHeaderExtensionUriPlayoutDelay,
			PreferredId:      14,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
		{
			Kind:             MediaKindVideo,
			Uri:              RtpHeaderExtensionUriPlayoutDelay,
			PreferredId:      14,
			PreferredEncrypt: boolPtr(false),
			Direction:        RtpHeaderExtensionDirectionSendRecv,
		},
	},
}

// Helper functions to create pointer values
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
