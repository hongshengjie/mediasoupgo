package mediasoupgo

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"slices"
	"strings"

	fbsRtpParameters "mediasoupgo/FBS/RtpParameters"
	"mediasoupgo/h264"
	"mediasoupgo/ptr"
)

type RtpCodecsEncodingsMapping struct {
	Codecs []struct {
		PayloadType       byte
		MappedPayloadType byte
	}
	Encodings []struct {
		Ssrc            *uint32
		Rid             *string
		ScalabilityMode *string
		MappedSsrc      uint32
	}
}

var DynamicPayloadTypes = []byte{
	100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114,
	115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 96, 97, 98,
	99,
}

// ValidateRtpCapabilities validates RtpCapabilities
func ValidateRtpCapabilities(caps *RtpCapabilities) error {
	if caps == nil {
		return errors.New("caps is not an object")
	}

	if caps.Codecs == nil {
		caps.Codecs = []*RtpCodecCapability{}
	}

	for _, codec := range caps.Codecs {
		if err := ValidateRtpCodecCapability(codec); err != nil {
			return err
		}
	}

	if caps.HeaderExtensions == nil {
		caps.HeaderExtensions = []*RtpHeaderExtension{}
	}

	for _, ext := range caps.HeaderExtensions {
		if err := ValidateRtpHeaderExtension(ext); err != nil {
			return err
		}
	}

	return nil
}

// ValidateRtpParameters validates RtpParameters
func ValidateRtpParameters(params *RtpParameters) error {
	if params == nil {
		return errors.New("params is not an object")
	}

	// if params.MID != nil && len(*params.MID) == 0 {
	// 	return errors.New("params.mid is not a string")
	// }

	if len(params.Codecs) == 0 {
		return errors.New("missing params.codecs")
	}

	for _, codec := range params.Codecs {
		if err := ValidateRtpCodecParameters(codec); err != nil {
			return err
		}
	}

	if params.HeaderExtensions == nil {
		params.HeaderExtensions = []*RtpHeaderExtensionParameters{}
	}

	for _, ext := range params.HeaderExtensions {
		if err := ValidateRtpHeaderExtensionParameters(ext); err != nil {
			return err
		}
	}

	if params.Encodings == nil {
		params.Encodings = []*RtpEncodingParameters{}
	}

	for _, encoding := range params.Encodings {
		if err := ValidateRtpEncodingParameters(encoding); err != nil {
			return err
		}
	}

	if params.RTCP == nil {
		params.RTCP = &RtcpParameters{}
	}

	return ValidateRtcpParameters(params.RTCP)
}

// ValidateSctpStreamParameters validates SctpStreamParameters
func ValidateSctpStreamParameters(params *SctpStreamParameters) error {
	if params == nil {
		return errors.New("params is not an object")
	}

	// if params.StreamID == 0 {
	// 	return errors.New("missing params.streamId")
	// }

	orderedGiven := params.Ordered != nil
	if params.Ordered == nil {
		ordered := true
		params.Ordered = &ordered
	}

	// if params.MaxPacketLifeTime != nil && *params.MaxPacketLifeTime < 0 {
	// 	return errors.New("invalid params.maxPacketLifeTime")
	// }

	// if params.MaxRetransmits != nil && *params.MaxRetransmits < 0 {
	// 	return errors.New("invalid params.maxRetransmits")
	// }

	if params.MaxPacketLifeTime != nil && params.MaxRetransmits != nil {
		return errors.New("cannot provide both maxPacketLifeTime and maxRetransmits")
	}

	if orderedGiven && *params.Ordered &&
		(params.MaxPacketLifeTime != nil || params.MaxRetransmits != nil) {
		return errors.New("cannot be ordered with maxPacketLifeTime or maxRetransmits")
	}

	if !orderedGiven && (params.MaxPacketLifeTime != nil || params.MaxRetransmits != nil) {
		ordered := false
		params.Ordered = &ordered
	}

	return nil
}

// GenerateRouterRtpCapabilities generates RTP capabilities for the Router
func GenerateRouterRtpCapabilities(mediaCodecs []*RtpCodecCapability) (*RtpCapabilities, error) {
	// Assuming supportedRtpCapabilities is defined globally or passed as parameter
	if err := ValidateRtpCapabilities(&SupportedRtpCapabilities); err != nil {
		return nil, err
	}

	clonedSupportedRtpCapabilities := cloneRtpCapabilities(SupportedRtpCapabilities)
	dynamicPayloadTypes := slices.Clone(DynamicPayloadTypes)
	caps := RtpCapabilities{
		Codecs:           []*RtpCodecCapability{},
		HeaderExtensions: clonedSupportedRtpCapabilities.HeaderExtensions,
	}

	for _, mediaCodec := range mediaCodecs {
		if err := ValidateRtpCodecCapability(mediaCodec); err != nil {
			return nil, err
		}
		// fmt.Println("support", *clonedSupportedRtpCapabilities.Codecs[0])
		matchedSupportedCodec := findMatchingCodec(
			clonedSupportedRtpCapabilities.Codecs,
			mediaCodec,
		)
		if matchedSupportedCodec == nil {
			return nil, fmt.Errorf(
				"media codec not supported [mimeType:%s]",
				mediaCodec.MimeType,
			)
		}
		//		fmt.Println("match", *matchedSupportedCodec)
		codec := cloneRtpCodecCapability(*matchedSupportedCodec)

		if mediaCodec.PreferredPayloadType != nil {
			codec.PreferredPayloadType = mediaCodec.PreferredPayloadType
			if idx := indexOf(dynamicPayloadTypes, *codec.PreferredPayloadType); idx != -1 {
				dynamicPayloadTypes = append(
					dynamicPayloadTypes[:idx],
					dynamicPayloadTypes[idx+1:]...)
			}
		} else if codec.PreferredPayloadType == nil {
			if len(dynamicPayloadTypes) == 0 {
				return nil, errors.New("cannot allocate more dynamic codec payload types")
			}
			pt := dynamicPayloadTypes[0]
			dynamicPayloadTypes = dynamicPayloadTypes[1:]
			codec.PreferredPayloadType = &pt
		}

		for _, c := range caps.Codecs {
			if c.PreferredPayloadType != nil && codec.PreferredPayloadType != nil &&
				c.PreferredPayloadType == codec.PreferredPayloadType {
				return nil, errors.New("duplicated codec.preferredPayloadType")
			}
		}

		codec.Parameters = mergeParameters(codec.Parameters, mediaCodec.Parameters)
		caps.Codecs = append(caps.Codecs, &codec)

		if codec.Kind == VideoMediaKind {
			if len(dynamicPayloadTypes) == 0 {
				return nil, errors.New(
					"cannot allocate more dynamic codec payload types",
				)
			}
			pt := dynamicPayloadTypes[0]
			dynamicPayloadTypes = dynamicPayloadTypes[1:]
			rtxCodec := &RtpCodecCapability{
				Kind:                 codec.Kind,
				MimeType:             string(codec.Kind) + "/rtx",
				PreferredPayloadType: &pt,
				ClockRate:            codec.ClockRate,
				Parameters: map[string]any{
					"apt": codec.PreferredPayloadType,
				},
			}
			caps.Codecs = append(caps.Codecs, rtxCodec)
		}
	}

	return &caps, nil
}

// GetProducerRtpParametersMapping generates codec payloads and encodings mapping
func GetProducerRtpParametersMapping(
	params *RtpParameters,
	caps *RtpCapabilities,
) (*RtpCodecsEncodingsMapping, error) {
	rtpMapping := &RtpCodecsEncodingsMapping{}
	codecToCapCodec := make(map[*RtpCodecParameters]*RtpCodecCapability)

	for _, codec := range params.Codecs {
		if IsRtxCodec(codec) {
			continue
		}
		var matchedCapCodec *RtpCodecCapability

		for _, capCodec := range caps.Codecs {
			fmt.Println("type two", *capCodec, *codec)

			if MatchCodecs(codec, capCodec, true, true) {
				matchedCapCodec = capCodec
				break
			}
		}
		if matchedCapCodec == nil {
			return nil, fmt.Errorf(
				"unsupported codec [mimeType:%s, payloadType:%d]",
				codec.MimeType,
				codec.PayloadType,
			)
		}
		codecToCapCodec[codec] = matchedCapCodec
	}

	for _, codec := range params.Codecs {
		if !IsRtxCodec(codec) {
			continue
		}

		if IsRtxCodec(codec) {
			continue
		}
		apt, ok := codec.Parameters["apt"]
		if !ok {
			return nil, fmt.Errorf(
				"missing apt parameter for RTX PT %d",
				codec.PayloadType,
			)
		}
		apt1, _ := toUint8(apt)

		associatedMediaCodec := findAssociatedMediaCodec(params.Codecs, apt1)
		if associatedMediaCodec == nil {
			return nil, fmt.Errorf(
				"missing media codec for RTX PT %d",
				codec.PayloadType,
			)
		}
		capMediaCodec := codecToCapCodec[associatedMediaCodec]
		associatedCapRtxCodec := findRtxCodec(caps.Codecs, *capMediaCodec.PreferredPayloadType)
		if associatedCapRtxCodec == nil {
			return nil, fmt.Errorf(
				"no RTX codec for capability codec PT %d",
				capMediaCodec.PreferredPayloadType,
			)
		}
		codecToCapCodec[codec] = associatedCapRtxCodec
	}

	for codec, capCodec := range codecToCapCodec {
		rtpMapping.Codecs = append(rtpMapping.Codecs,
			struct{ PayloadType, MappedPayloadType byte }{
				PayloadType:       codec.PayloadType,
				MappedPayloadType: *capCodec.PreferredPayloadType,
			})
	}

	mappedSsrc := generateRandomNumber()
	for _, encoding := range params.Encodings {
		mappedEncoding := struct {
			Ssrc                 *uint32
			Rid, ScalabilityMode *string
			MappedSsrc           uint32
		}{
			MappedSsrc: mappedSsrc,
		}
		mappedSsrc++
		if encoding.RID != nil {
			mappedEncoding.Rid = encoding.RID
		}
		if encoding.SSRC != nil {
			mappedEncoding.Ssrc = encoding.SSRC
		}
		if encoding.ScalabilityMode != nil {
			mappedEncoding.ScalabilityMode = encoding.ScalabilityMode
		}
		rtpMapping.Encodings = append(rtpMapping.Encodings, mappedEncoding)
	}

	return rtpMapping, nil
}

// GetConsumableRtpParameters generates RTP parameters for Consumers
func GetConsumableRtpParameters(
	kind string,
	params RtpParameters,
	caps RtpCapabilities,
	rtpMapping RtpCodecsEncodingsMapping,
) *RtpParameters {
	consumableParams := &RtpParameters{
		RTCP: &RtcpParameters{},
	}

	for _, codec := range params.Codecs {
		if IsRtxCodec(codec) {
			continue
		}
		consumableCodecPt := findMappedPayloadType(rtpMapping.Codecs, codec.PayloadType)
		matchedCapCodec := findCodecByPayloadType(caps.Codecs, consumableCodecPt)
		consumableCodec := RtpCodecParameters{
			MimeType:     matchedCapCodec.MimeType,
			PayloadType:  *matchedCapCodec.PreferredPayloadType,
			ClockRate:    matchedCapCodec.ClockRate,
			Channels:     matchedCapCodec.Channels,
			Parameters:   codec.Parameters,
			RTCPFeedback: matchedCapCodec.RTCPFeedback,
		}
		consumableParams.Codecs = append(consumableParams.Codecs, &consumableCodec)

		consumableCapRtxCodec := findRtxCodec(caps.Codecs, consumableCodec.PayloadType)
		if consumableCapRtxCodec != nil {
			consumableRtxCodec := RtpCodecParameters{
				MimeType:     consumableCapRtxCodec.MimeType,
				PayloadType:  *consumableCapRtxCodec.PreferredPayloadType,
				ClockRate:    consumableCapRtxCodec.ClockRate,
				Parameters:   consumableCapRtxCodec.Parameters,
				RTCPFeedback: consumableCapRtxCodec.RTCPFeedback,
			}
			consumableParams.Codecs = append(consumableParams.Codecs, &consumableRtxCodec)
		}
	}

	for _, capExt := range caps.HeaderExtensions {
		if capExt.Kind != MediaKind(kind) ||
			(*capExt.Direction != "" && *capExt.Direction != "sendrecv" && *capExt.Direction != "sendonly") {
			continue
		}
		consumableExt := RtpHeaderExtensionParameters{
			URI:        capExt.URI,
			ID:         capExt.PreferredID,
			Encrypt:    capExt.PreferredEncrypt,
			Parameters: map[string]interface{}{},
		}
		consumableParams.HeaderExtensions = append(
			consumableParams.HeaderExtensions,
			&consumableExt,
		)
	}

	consumableEncodings := slices.Clone(params.Encodings)
	for i, encoding := range consumableEncodings {
		encoding.SSRC = &rtpMapping.Encodings[i].MappedSsrc
		encoding.RID = nil
		encoding.RTX = nil
		consumableParams.Encodings = append(consumableParams.Encodings, encoding)
	}

	if params.RTCP != nil {
		consumableParams.RTCP.CNAME = params.RTCP.CNAME
	}
	reducedSize := true
	consumableParams.RTCP.ReducedSize = ptr.To(reducedSize)

	return consumableParams
}

// CanConsume checks if RTP capabilities can consume a Producer
func CanConsume(consumableParams *RtpParameters, caps *RtpCapabilities) (bool, error) {
	if err := ValidateRtpCapabilities(caps); err != nil {
		return false, err
	}

	matchingCodecs := []*RtpCodecParameters{}
	for _, codec := range consumableParams.Codecs {
		var findmatch bool
		for _, capcodec := range caps.Codecs {
			if MatchCodecs(capcodec, codec, true, false) {
				findmatch = true
				break
			}
		}
		if !findmatch {
			continue
		}
		matchingCodecs = append(matchingCodecs, codec)
	}

	if len(matchingCodecs) == 0 || IsRtxCodec(matchingCodecs[0]) {
		return false, nil
	}
	return true, nil
}

// GetConsumerRtpParameters generates RTP parameters for a specific Consumer
func GetConsumerRtpParameters(
	consumableRtpParameters RtpParameters,
	remoteRtpCapabilities RtpCapabilities,
	pipe bool,
	enableRtx bool,
) (*RtpParameters, error) {
	consumerParams := &RtpParameters{
		RTCP: consumableRtpParameters.RTCP,
	}

	for _, capCodec := range remoteRtpCapabilities.Codecs {
		if err := ValidateRtpCodecCapability(capCodec); err != nil {
			return nil, err
		}
	}

	consumableCodecs := slices.Clone(consumableRtpParameters.Codecs)
	rtxSupported := false

	for _, codec := range consumableCodecs {
		if !enableRtx && IsRtxCodec(codec) {
			continue
		}
		var matchedCapCodec *RtpCodecCapability
		for _, capcodec := range remoteRtpCapabilities.Codecs {
			if MatchCodecs(capcodec, codec, true, false) {
				matchedCapCodec = capcodec
			}
		}
		if matchedCapCodec == nil {
			continue
		}
		codec.RTCPFeedback = filterRtcpFeedback(matchedCapCodec.RTCPFeedback, enableRtx)
		consumerParams.Codecs = append(consumerParams.Codecs, codec)
	}

	for i := len(consumerParams.Codecs) - 1; i >= 0; i-- {
		codec := consumerParams.Codecs[i]
		if IsRtxCodec(codec) {
			aptx, _ := toUint8(codec.Parameters["apt"])
			apt := byte(aptx)
			if findAssociatedMediaCodec(consumerParams.Codecs, apt) != nil {
				rtxSupported = true
			} else {
				consumerParams.Codecs = slices.Delete(consumerParams.Codecs, i, i+1)
			}
		}
	}

	if len(consumerParams.Codecs) == 0 || IsRtxCodec(consumerParams.Codecs[0]) {
		return nil, errors.New("no compatible media codecs")
	}

	consumerParams.HeaderExtensions = filterHeaderExtensions(
		consumableRtpParameters.HeaderExtensions,
		remoteRtpCapabilities.HeaderExtensions,
	)

	if hasTransportCC(consumerParams.HeaderExtensions) {
		for i := range consumerParams.Codecs {
			consumerParams.Codecs[i].RTCPFeedback = filterFeedback(
				consumerParams.Codecs[i].RTCPFeedback,
				"goog-remb",
			)
		}
	} else if hasAbsSendTime(consumerParams.HeaderExtensions) {
		for i := range consumerParams.Codecs {
			consumerParams.Codecs[i].RTCPFeedback = filterFeedback(consumerParams.Codecs[i].RTCPFeedback, "transport-cc")
		}
	} else {
		for i := range consumerParams.Codecs {
			consumerParams.Codecs[i].RTCPFeedback = filterFeedback(consumerParams.Codecs[i].RTCPFeedback, "transport-cc", "goog-remb")
		}
	}

	if !pipe {
		consumerEncoding := &RtpEncodingParameters{}
		consumerEncoding.SSRC = ptr.To(generateRandomNumber())
		if rtxSupported {
			rtxSsrc := *consumerEncoding.SSRC + 1
			consumerEncoding.RTX = &RTX{SSRC: rtxSsrc}
		}

		encodingWithScalabilityMode := findEncodingWithScalabilityMode(
			consumableRtpParameters.Encodings,
		)
		scalabilityMode := ""
		if encodingWithScalabilityMode != nil {
			scalabilityMode = *encodingWithScalabilityMode.ScalabilityMode
		}

		if len(consumableRtpParameters.Encodings) > 1 {
			temporalLayers := parseScalabilityMode(scalabilityMode).TemporalLayers
			scalabilityMode = fmt.Sprintf(
				"L%dT%d",
				len(consumableRtpParameters.Encodings),
				temporalLayers,
			)
		}

		if scalabilityMode != "" {
			consumerEncoding.ScalabilityMode = &scalabilityMode
		}

		maxBitrate := maxEncodingMaxBitrate(consumableRtpParameters.Encodings)
		if maxBitrate > 0 {
			consumerEncoding.MaxBitrate = &maxBitrate
		}

		consumerParams.Encodings = append(consumerParams.Encodings, consumerEncoding)
	} else {
		consumableEncodings := slices.Clone(consumableRtpParameters.Encodings)
		baseSsrc := generateRandomNumber()
		baseRtxSsrc := generateRandomNumber()
		for i := range consumableEncodings {
			consumableEncodings[i].SSRC = ptr.To(baseSsrc + uint32(i))
			if enableRtx {
				consumableEncodings[i].RTX = &RTX{SSRC: baseRtxSsrc + uint32(i)}
			} else {
				consumableEncodings[i].RTX = nil
			}
			consumerParams.Encodings = append(consumerParams.Encodings, consumableEncodings[i])
		}
	}

	return consumerParams, nil
}

// GetPipeConsumerRtpParameters generates RTP parameters for a pipe Consumer
func GetPipeConsumerRtpParameters(
	consumableRtpParameters *RtpParameters,
	enableRtx bool,
) *RtpParameters {
	consumerParams := &RtpParameters{
		Codecs:           []*RtpCodecParameters{},
		HeaderExtensions: []*RtpHeaderExtensionParameters{},
		Encodings:        []*RtpEncodingParameters{},
		RTCP:             consumableRtpParameters.RTCP,
	}

	consumableCodecs := slices.Clone(consumableRtpParameters.Codecs)
	for _, codec := range consumableCodecs {
		if !enableRtx && IsRtxCodec(codec) {
			continue
		}
		codec.RTCPFeedback = filterPipeRtcpFeedback(codec.RTCPFeedback, enableRtx)
		consumerParams.Codecs = append(consumerParams.Codecs, codec)
	}

	consumerParams.HeaderExtensions = filterPipeHeaderExtensions(
		consumableRtpParameters.HeaderExtensions,
	)

	consumableEncodings := slices.Clone(consumableRtpParameters.Encodings)
	baseSsrc := generateRandomNumber()
	baseRtxSsrc := generateRandomNumber()
	for i := range consumableEncodings {
		consumableEncodings[i].SSRC = ptr.To(baseSsrc + uint32(i))
		if enableRtx {
			rtxSsrc := baseRtxSsrc + uint32(i)
			consumableEncodings[i].RTX = &RTX{SSRC: rtxSsrc}
		} else {
			consumableEncodings[i].RTX = nil
		}
		consumerParams.Encodings = append(consumerParams.Encodings, consumableEncodings[i])
	}

	return consumerParams
}

// IsRtxCodec checks if codec is RTX
func IsRtxCodec(codec interface{}) bool {
	var mimeType string
	switch c := codec.(type) {
	case *RtpCodecCapability:
		mimeType = c.MimeType
	case *RtpCodecParameters:
		mimeType = c.MimeType
	default:
		return false
	}
	return regexp.MustCompile(`.+\/rtx$`).MatchString(strings.ToLower(mimeType))
}

// MatchCodecs compares codecs for compatibility
func MatchCodecs(aCodec, bCodec interface{}, strict, modify bool) bool {
	var aMimeType, bMimeType string
	var aClockRate, bClockRate uint32
	var aChannels, bChannels *byte
	var aParams, bParams map[string]interface{}

	switch a := aCodec.(type) {
	case *RtpCodecCapability:
		aMimeType = strings.ToLower(a.MimeType)
		aClockRate = a.ClockRate
		aChannels = a.Channels
		aParams = a.Parameters
	case *RtpCodecParameters:
		aMimeType = strings.ToLower(a.MimeType)
		aClockRate = a.ClockRate
		aChannels = a.Channels
		aParams = a.Parameters
	default:
		return false
	}

	switch b := bCodec.(type) {
	case *RtpCodecCapability:
		bMimeType = strings.ToLower(b.MimeType)
		bClockRate = b.ClockRate
		bChannels = b.Channels
		bParams = b.Parameters
	case *RtpCodecParameters:
		bMimeType = strings.ToLower(b.MimeType)
		bClockRate = b.ClockRate
		bChannels = b.Channels
		bParams = b.Parameters
	default:
		return false
	}
	if aMimeType != bMimeType || aClockRate != bClockRate {
		fmt.Println(1)
		return false
	}
	if aChannels != nil && bChannels == nil {

		fmt.Println(2)
		return false
	}
	if aChannels == nil && bChannels != nil {

		fmt.Println(3)
		return false
	}
	if aChannels != nil && bChannels != nil && *aChannels != *bChannels {

		fmt.Println(4)
		return false
	}

	switch aMimeType {
	case "audio/multiopus":
		if aNumStreams, aOk := aParams["num_streams"]; aOk {
			if bNumStreams, bOk := bParams["num_streams"]; !bOk || aNumStreams != bNumStreams {
				return false
			}
		}
		if aCoupledStreams, aOk := aParams["coupled_streams"]; aOk {
			if bCoupledStreams, bOk := bParams["coupled_streams"]; !bOk ||
				aCoupledStreams != bCoupledStreams {
				return false
			}
		}
	case "video/h264":
	case "video/h264-svc":
		if strict {
			aPacketizationMode := getParameter(aParams, "packetization-mode", 0)
			bPacketizationMode := getParameter(bParams, "packetization-mode", 0)
			if aPacketizationMode != bPacketizationMode {
				return false
			}
			if !h264.IsSameProfile(aParams, bParams) {
				return false
			}
			selectedProfileLevelId, err := h264.GenerateProfileLevelIdStringForAnswer(
				aParams,
				bParams,
			)
			if err != nil {
				return false
			}
			if modify {
				if selectedProfileLevelId != "" {
					aParams["profile-level-id"] = selectedProfileLevelId
				} else {
					delete(aParams, "profile-level-id")
				}
			}
		}
	case "video/vp9":
		if strict {
			aProfileId := getParameter(aParams, "profile-id", 0)
			bProfileId := getParameter(bParams, "profile-id", 0)
			if aProfileId != bProfileId {
				return false
			}
		}
	}

	return true
}

func SerializeRtpMapping(a *RtpCodecsEncodingsMapping) *fbsRtpParameters.RtpMappingT {
	ret := &fbsRtpParameters.RtpMappingT{}
	for _, v := range a.Codecs {
		ret.Codecs = append(ret.Codecs, &fbsRtpParameters.CodecMappingT{
			PayloadType:       v.PayloadType,
			MappedPayloadType: v.MappedPayloadType,
		})
	}

	for _, v := range a.Encodings {
		i := &fbsRtpParameters.EncodingMappingT{
			Ssrc:       v.Ssrc,
			MappedSsrc: v.MappedSsrc,
		}
		if v.Rid != nil {
			i.Rid = *v.Rid
		}
		if v.ScalabilityMode != nil {
			i.ScalabilityMode = *v.ScalabilityMode
		}
		ret.Encodings = append(ret.Encodings, i)
	}
	return ret
}

// ValidateRtpCodecCapability validates RtpCodecCapability
func ValidateRtpCodecCapability(codec *RtpCodecCapability) error {
	if codec == nil {
		return errors.New("codec is not an object")
	}

	// if codec.MimeType == "" {
	// 	return errors.New("missing codec.mimeType")
	// }

	mimeTypeRegex := regexp.MustCompile(`^(audio|video)/(.+)`)
	mimeTypeMatch := mimeTypeRegex.FindStringSubmatch(codec.MimeType)
	if len(mimeTypeMatch) == 0 {
		return errors.New("invalid codec.mimeType")
	}
	// fmt.Println(mimeTypeMatch[0], mimeTypeMatch[1])
	codec.Kind = MediaKind(strings.ToLower(mimeTypeMatch[1]))

	// if codec.PreferredPayloadType != nil && *codec.PreferredPayloadType < 0 {
	// 	return errors.New("invalid codec.preferredPayloadType")
	// }

	// if codec.ClockRate == 0 {
	// 	return errors.New("missing codec.clockRate")
	// }

	if codec.Kind == AudioMediaKind {
		if codec.Channels == nil {
			codec.Channels = ptr.To[byte](1)
		}
	} else {
		codec.Channels = nil
	}

	if codec.Parameters == nil {
		codec.Parameters = make(map[string]interface{})
	}

	for key, value := range codec.Parameters {
		if value == nil {
			codec.Parameters[key] = ""
			value = ""
		}
		if _, ok := value.(string); !ok {
			if ok := isNumber(value); !ok {
				return fmt.Errorf("invalid codec parameter [key:%s, value:%v]", key, value)
			}
		}
		if key == "apt" {
			if ok := isNumber(value); !ok {
				return errors.New("invalid codec apt parameter")
			}
		}
	}

	if codec.RTCPFeedback == nil {
		codec.RTCPFeedback = []*RtcpFeedback{}
	}

	for _, fb := range codec.RTCPFeedback {
		if err := ValidateRtcpFeedback(fb); err != nil {
			return err
		}
	}

	return nil
}

// ValidateRtcpFeedback validates RtcpFeedback
func ValidateRtcpFeedback(fb *RtcpFeedback) error {
	if fb == nil {
		return errors.New("fb is not an object")
	}

	// if fb.Type == "" {
	// 	return errors.New("missing fb.type")
	// }

	if fb.Parameter == nil {
		fb.Parameter = ptr.To("")
	}

	return nil
}

// ValidateRtpHeaderExtension validates RtpHeaderExtension
func ValidateRtpHeaderExtension(ext *RtpHeaderExtension) error {
	if ext == nil {
		return errors.New("ext is not an object")
	}

	if ext.Kind != AudioMediaKind && ext.Kind != VideoMediaKind {
		return errors.New("invalid ext.kind")
	}

	if ext.URI == "" {
		return errors.New("missing ext.uri")
	}

	if ext.PreferredID == 0 {
		return errors.New("missing ext.preferredId")
	}

	// if ext.PreferredEncrypt == nil {
	// 	encrypt := false
	// 	ext.PreferredEncrypt = &encrypt
	// }

	if ext.Direction == nil {
		ext.Direction = ptr.To(RtpHeaderExtensionDirection("sendrecv"))
	}

	return nil
}

// ValidateRtpCodecParameters validates RtpCodecParameters
func ValidateRtpCodecParameters(codec *RtpCodecParameters) error {
	if codec == nil {
		return errors.New("codec is not an object")
	}

	// if codec.MimeType == "" {
	// 	return errors.New("missing codec.mimeType")
	// }

	mimeTypeRegex := regexp.MustCompile(`^(audio|video)/(.+)`)
	mimeTypeMatch := mimeTypeRegex.FindStringSubmatch(codec.MimeType)
	if len(mimeTypeMatch) == 0 {
		return errors.New("invalid codec.mimeType")
	}

	// if codec.PayloadType == 0 {
	// 	return errors.New("missing codec.payloadType")
	// }
	//
	// if codec.ClockRate == 0 {
	// 	return errors.New("missing codec.clockRate")
	// }

	kind := MediaKind(strings.ToLower(mimeTypeMatch[1]))
	if kind == AudioMediaKind {
		if codec.Channels == nil {
			codec.Channels = ptr.To[byte](1)
		}
	} else {
		codec.Channels = nil
	}

	if codec.Parameters == nil {
		codec.Parameters = make(map[string]interface{})
	}

	for key, value := range codec.Parameters {
		if value == nil {
			codec.Parameters[key] = ""
			value = ""
		}
		if _, ok := value.(string); !ok {
			if ok := isNumber(value); !ok {
				return fmt.Errorf("invalid codec parameter [key:%s, value:%v]", key, value)
			}
		}
		if key == "apt" {
			if ok := isNumber(value); !ok {
				return errors.New("invalid codec apt parameter")
			}
		}
	}

	if codec.RTCPFeedback == nil {
		codec.RTCPFeedback = []*RtcpFeedback{}
	}

	for _, fb := range codec.RTCPFeedback {
		if err := ValidateRtcpFeedback(fb); err != nil {
			return err
		}
	}

	return nil
}

// ValidateRtpHeaderExtensionParameters validates RtpHeaderExtensionParameters
func ValidateRtpHeaderExtensionParameters(ext *RtpHeaderExtensionParameters) error {
	if ext == nil {
		return errors.New("ext is not an object")
	}

	// if ext.URI == "" {
	// 	return errors.New("missing ext.uri")
	// }
	//
	// if ext.ID == 0 {
	// 	return errors.New("missing ext.id")
	// }

	if ext.Encrypt == nil {
		ext.Encrypt = ptr.To(false)
	}

	if ext.Parameters == nil {
		ext.Parameters = make(map[string]interface{})
	}

	for key, value := range ext.Parameters {
		if value == nil {
			ext.Parameters[key] = ""
			value = ""
		}
		if _, ok := value.(string); !ok {
			if ok := isNumber(value); !ok {
				return errors.New("invalid header extension parameter")
			}
		}
	}

	return nil
}

// ValidateRtpEncodingParameters validates RtpEncodingParameters
func ValidateRtpEncodingParameters(encoding *RtpEncodingParameters) error {
	if encoding == nil {
		return errors.New("encoding is not an object")
	}

	// if encoding.SSRC != 0 && encoding.SSRC < 0 {
	// 	return errors.New("invalid encoding.ssrc")
	// }

	// if encoding.RID == "" {
	// 	return errors.New("invalid encoding.rid")
	// }

	if encoding.RTX != nil && encoding.RTX.SSRC < 0 {
		return errors.New("missing encoding.rtx.ssrc")
	}

	if encoding.DTX == nil {
		encoding.DTX = ptr.To(false)
	}

	return nil
}

// ValidateRtcpParameters validates RtcpParameters
func ValidateRtcpParameters(rtcp *RtcpParameters) error {
	if rtcp == nil {
		return errors.New("rtcp is not an object")
	}
	//
	// if rtcp.CNAME == nil {
	// 	return errors.New("invalid rtcp.cname")
	// }

	if rtcp.ReducedSize == nil {
		rtcp.ReducedSize = ptr.To(true)
	}

	return nil
}

// Helper functions (placeholders for actual implementations)
func cloneRtpCapabilities(caps RtpCapabilities) RtpCapabilities {
	return caps // Implement deep clone
}

func cloneRtpCodecCapability(codec RtpCodecCapability) RtpCodecCapability {
	return codec // Implement deep clone
}

func findMatchingCodec(
	codecs []*RtpCodecCapability,
	mediaCodec *RtpCodecCapability,
) *RtpCodecCapability {
	for _, codec := range codecs {
		if MatchCodecs(mediaCodec, codec, true, false) {
			return codec
		}
	}
	return nil
}

func findCodecByPayloadType(codecs []*RtpCodecCapability, pt byte) *RtpCodecCapability {
	for _, c := range codecs {
		if c.PreferredPayloadType != nil && *c.PreferredPayloadType == pt {
			return c
		}
	}
	return nil
}

func findRtxCodec(codecs []*RtpCodecCapability, apt byte) *RtpCodecCapability {
	for _, c := range codecs {
		if IsRtxCodec(&RtpCodecParameters{MimeType: c.MimeType}) {
			aptVal, _ := toUint8(c.Parameters["apt"])
			if byte(aptVal) == apt {
				return c
			}
		}
	}
	return nil
}

func findAssociatedMediaCodec(codecs []*RtpCodecParameters, apt byte) *RtpCodecParameters {
	for _, c := range codecs {
		if c.PayloadType == apt {
			return c
		}
	}
	return nil
}

func findMappedPayloadType(codecs []struct{ PayloadType, MappedPayloadType byte }, pt byte) byte {
	for _, c := range codecs {
		if c.PayloadType == pt {
			return c.MappedPayloadType
		}
	}
	return 0
}

func mergeParameters(a, b map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}

func generateRandomNumber() uint32 {
	return uint32(rand.Int31n(999999999-100000000+1) + 100000000)
}

func filterRtcpFeedback(fb []*RtcpFeedback, enableRtx bool) []*RtcpFeedback {
	result := []*RtcpFeedback{}
	for _, f := range fb {
		if enableRtx || f.Type != "nack" || (f.Parameter != nil) {
			result = append(result, f)
		}
	}
	return result
}

func filterHeaderExtensions(
	exts []*RtpHeaderExtensionParameters,
	capExts []*RtpHeaderExtension,
) []*RtpHeaderExtensionParameters {
	result := []*RtpHeaderExtensionParameters{}
	for _, ext := range exts {
		for _, capExt := range capExts {
			if ext.ID == capExt.PreferredID &&
				ext.URI == capExt.URI {
				result = append(result, ext)
				break
			}
		}
	}
	return result
}

func hasTransportCC(exts []*RtpHeaderExtensionParameters) bool {
	for _, ext := range exts {
		if ext.URI == TransportWideCC01RtpHeaderExtensionUri {
			return true
		}
	}
	return false
}

func hasAbsSendTime(exts []*RtpHeaderExtensionParameters) bool {
	for _, ext := range exts {
		if ext.URI == AbsSendTimeRtpHeaderExtensionUri {
			return true
		}
	}
	return false
}

func filterFeedback(fb []*RtcpFeedback, types ...string) []*RtcpFeedback {
	result := []*RtcpFeedback{}
	for _, f := range fb {

		exclude := slices.Contains(types, f.Type)
		if !exclude {
			result = append(result, f)
		}
	}
	return result
}

func findEncodingWithScalabilityMode(encodings []*RtpEncodingParameters) *RtpEncodingParameters {
	for _, encoding := range encodings {
		if encoding.ScalabilityMode != nil {
			return encoding
		}
	}
	return nil
}

func maxEncodingMaxBitrate(encodings []*RtpEncodingParameters) uint32 {
	maxBitrate := uint32(0)
	for _, encoding := range encodings {
		if encoding.MaxBitrate != nil && *encoding.MaxBitrate > maxBitrate {
			maxBitrate = *encoding.MaxBitrate
		}
	}
	return maxBitrate
}

func filterPipeRtcpFeedback(fb []*RtcpFeedback, enableRtx bool) []*RtcpFeedback {
	result := []*RtcpFeedback{}
	for _, f := range fb {
		if (f.Type == "nack" && f.Parameter != nil && *f.Parameter == "pli") ||
			(f.Type == "ccm" && f.Parameter != nil && *f.Parameter == "fir") ||
			(enableRtx && f.Type == "nack" && f.Parameter == nil) {
			result = append(result, f)
		}
	}
	return result
}

func filterPipeHeaderExtensions(
	exts []*RtpHeaderExtensionParameters,
) []*RtpHeaderExtensionParameters {
	result := []*RtpHeaderExtensionParameters{}
	for _, ext := range exts {
		if ext.URI != "urn:ietf:params:rtp-hdrext:sdes:mid" &&
			ext.URI != "http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time" &&
			ext.URI != "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01" {
			result = append(result, ext)
		}
	}
	return result
}

func indexOf(slice []byte, val byte) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

func getParameter(params map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if val, ok := params[key]; ok {
		return val
	}
	return defaultValue
}

func isNumber(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}

func toUint8(v interface{}) (uint8, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Uint8:
		return uint8(rv.Uint()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val := rv.Int(); val >= 0 && val <= 255 {
			return uint8(val), nil
		}
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val := rv.Uint(); val <= 255 {
			return uint8(val), nil
		}
	case reflect.Float32, reflect.Float64:
		if val := rv.Float(); val >= 0 && val <= 255 && float64(int(val)) == val {
			return uint8(val), nil
		}
	default:
		return 0, fmt.Errorf("cannot convert %v to uint8", v)
	}
	return 0, fmt.Errorf("value %v out of uint8 range", v)
}
