package mediasoupgo

import (
	"strings"

	fbsRtpParameters "mediasoupgo/internal/FBS/RtpParameters"
	FbsWebRtcTransport "mediasoupgo/internal/FBS/WebRtcTransport"
)

func ToFbsRtpParameters(a *RtpParameters) *fbsRtpParameters.RtpParametersT {
	ret := &fbsRtpParameters.RtpParametersT{}
	if a.MID != nil {
		ret.Mid = *a.MID
	}
	for _, v := range a.Codecs {
		var rtcfeedbacks []*fbsRtpParameters.RtcpFeedbackT
		for _, x := range v.RTCPFeedback {
			fb := &fbsRtpParameters.RtcpFeedbackT{
				Type: x.Type,
			}
			if x.Parameter != nil {
				fb.Parameter = *x.Parameter
			}
		}
		c := &fbsRtpParameters.RtpCodecParametersT{
			MimeType:     v.MimeType,
			PayloadType:  v.PayloadType,
			ClockRate:    v.ClockRate,
			Channels:     v.Channels,
			Parameters:   ToFBSParameters(v.Parameters),
			RtcpFeedback: rtcfeedbacks,
		}

		ret.Codecs = append(ret.Codecs, c)
	}
	if a.RTCP != nil {
		ret.Rtcp = &fbsRtpParameters.RtcpParametersT{}
		if a.RTCP.CNAME != nil {
			ret.Rtcp.Cname = *a.RTCP.CNAME
		}
		if a.RTCP.ReducedSize != nil {
			ret.Rtcp.ReducedSize = *a.RTCP.ReducedSize
		}
	}
	for _, v := range a.HeaderExtensions {
		header := &fbsRtpParameters.RtpHeaderExtensionParametersT{
			Uri:        ToFBSHeaderExtensionUri(v.URI),
			Id:         v.ID,
			Parameters: ToFBSParameters(v.Parameters),
		}
		if v.Encrypt != nil {
			header.Encrypt = *v.Encrypt
		}
		ret.HeaderExtensions = append(ret.HeaderExtensions, header)

	}
	for _, v := range a.Encodings {
		enc := &fbsRtpParameters.RtpEncodingParametersT{
			Ssrc:             v.SSRC,
			CodecPayloadType: v.CodecPayloadType,
			MaxBitrate:       v.MaxBitrate,
		}
		if v.RID != nil {
			enc.Rid = *v.RID
		}
		if v.DTX != nil {
			enc.Dtx = *v.DTX
		}
		if v.ScalabilityMode != nil {
			enc.ScalabilityMode = *v.ScalabilityMode
		}
		if v.RTX != nil {
			enc.Rtx = &fbsRtpParameters.RtxT{Ssrc: v.RTX.SSRC}
		}
		ret.Encodings = append(ret.Encodings, enc)
	}
	return ret
}

func ToFBSRtpEncodingParameters(p []*RtpEncodingParameters) []*fbsRtpParameters.RtpEncodingParametersT {
	var ret []*fbsRtpParameters.RtpEncodingParametersT
	for _, v := range p {
		item := &fbsRtpParameters.RtpEncodingParametersT{
			Ssrc:             v.SSRC,
			CodecPayloadType: v.CodecPayloadType,
			MaxBitrate:       v.MaxBitrate,
		}
		if v.RID != nil {
			item.Rid = *v.RID
		}
		if v.RTX != nil {
			item.Rtx = &fbsRtpParameters.RtxT{Ssrc: v.RTX.SSRC}
		}
		if v.DTX != nil {
			item.Dtx = *v.DTX
		}
		if v.ScalabilityMode != nil {
			item.ScalabilityMode = *v.ScalabilityMode
		}
		ret = append(ret, item)
	}
	return ret
}

func ToFBSParameters(p map[string]any) []*fbsRtpParameters.ParameterT {
	var ret []*fbsRtpParameters.ParameterT
	for k, v := range p {

		value := &fbsRtpParameters.ValueT{}

		switch val := v.(type) {
		case int:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case int8:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case int16:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case int32:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case int64:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case uint:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case uint8:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case uint16:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case uint32:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case uint64:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case float32:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case float64:
			value.Type = fbsRtpParameters.ValueInteger32
			value.Value = &fbsRtpParameters.Integer32T{Value: int32(val)}
		case string:
			value.Type = fbsRtpParameters.ValueString
			value.Value = &fbsRtpParameters.StringT{Value: val}
		case bool:
			value.Type = fbsRtpParameters.ValueBoolean
			var b byte
			if val {
				b = 1
			} else {
				b = 0
			}
			value.Value = &fbsRtpParameters.BooleanT{Value: b}
		case []int32:
			value.Type = fbsRtpParameters.ValueInteger32Array
			value.Value = &fbsRtpParameters.Integer32ArrayT{Value: val}
		case []int:
			vals := make([]int32, 0, len(val))
			for _, y := range val {
				vals = append(vals, int32(y))
			}
			value.Type = fbsRtpParameters.ValueInteger32Array
			value.Value = &fbsRtpParameters.Integer32ArrayT{Value: vals}
		case []int64:
			vals := make([]int32, 0, len(val))
			for _, y := range val {
				vals = append(vals, int32(y))
			}
			value.Type = fbsRtpParameters.ValueInteger32Array
			value.Value = &fbsRtpParameters.Integer32ArrayT{Value: vals}
		}
		i := &fbsRtpParameters.ParameterT{
			Name:  k,
			Value: value,
		}
		ret = append(ret, i)
	}
	return ret
}

func ToFBSHeaderExtensionUri(uri RtpHeaderExtensionUri) fbsRtpParameters.RtpHeaderExtensionUri {
	switch uri {
	case MIDRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriMid

	case RTPStreamIDRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriRtpStreamId

	case RepairedRTPStreamIDRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriRepairRtpStreamId

	case FrameMarking07RtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriFrameMarkingDraft07

	case FrameMarkingRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriFrameMarking

	case SSRCAudioLevelRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriAudioLevel

	case VideoOrientationRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriVideoOrientation

	case ToffsetRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriTimeOffset

	case TransportWideCC01RtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriTransportWideCcDraft01

	case AbsSendTimeRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriAbsSendTime

	case AbsCaptureTimeRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriAbsCaptureTime

	case PlayoutDelayRtpHeaderExtensionUri:
		return fbsRtpParameters.RtpHeaderExtensionUriPlayoutDelay

	default:
		return fbsRtpParameters.EnumValuesRtpHeaderExtensionUri[string(uri)]
	}
}

func ToFbsDtlsFingerprint(item DtlsFingerprint) *FbsWebRtcTransport.FingerprintT {
	var algorithm FbsWebRtcTransport.FingerprintAlgorithm

	switch item.Algorithm {
	case "sha-1":
		algorithm = FbsWebRtcTransport.FingerprintAlgorithmSHA1

	case "sha-224":
		algorithm = FbsWebRtcTransport.FingerprintAlgorithmSHA224

	case "sha-256":
		algorithm = FbsWebRtcTransport.FingerprintAlgorithmSHA256

	case "sha-384":
		algorithm = FbsWebRtcTransport.FingerprintAlgorithmSHA384

	case "sha-512":
		algorithm = FbsWebRtcTransport.FingerprintAlgorithmSHA512

	default:
		key := strings.ToUpper(strings.Join(strings.Split(string(item.Algorithm), "-"), ""))
		algorithm = FbsWebRtcTransport.EnumValuesFingerprintAlgorithm[key]
	}

	// avoid mediasoup crash
	if len(item.Value) == 0 {
		item.Value = "unknown"
	}

	return &FbsWebRtcTransport.FingerprintT{
		Algorithm: algorithm,
		Value:     item.Value,
	}
}

func ToDtlsFingerprint(item *FbsWebRtcTransport.FingerprintT) DtlsFingerprint {
	var algorithm string

	switch item.Algorithm {
	case FbsWebRtcTransport.FingerprintAlgorithmSHA1:
		algorithm = "sha-1"

	case FbsWebRtcTransport.FingerprintAlgorithmSHA224:
		algorithm = "sha-224"

	case FbsWebRtcTransport.FingerprintAlgorithmSHA256:
		algorithm = "sha-256"

	case FbsWebRtcTransport.FingerprintAlgorithmSHA384:
		algorithm = "sha-384"

	case FbsWebRtcTransport.FingerprintAlgorithmSHA512:
		algorithm = "sha-512"

	default:
		algorithm = "sha-" + strings.TrimPrefix(item.Algorithm.String(), "SHA")
	}

	return DtlsFingerprint{
		Algorithm: FingerprintAlgorithm(algorithm),
		Value:     item.Value,
	}
}
