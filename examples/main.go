package main

import (
	"fmt"
	"log/slog"
	"time"

	"mediasoupgo"
	"mediasoupgo/internal/ptr"
)

func main() {
	w := mediasoupgo.NewWorker(&mediasoupgo.WorkerSettings{
		LogLevel:             ptr.To(mediasoupgo.Debug),
		LogTags:              mediasoupgo.DefaultLogTags(),
		RTCMinPort:           nil,
		RTCMaxPort:           nil,
		DTLSCertificateFile:  nil,
		DTLSPrivateKeyFile:   nil,
		LibwebrtcFieldTrials: nil,
		DisableLiburing:      nil,
	})
	defer w.Close()
	dumpresp, err := w.Dump()
	slog.Info("dump", "resp", dumpresp, "err", err)
	usage, err := w.GetResourceUsage()
	slog.Info("GetResourceUsage", "resp", usage, "err", err)

	w.UpdateSettings(
		&mediasoupgo.WorkerUpdateableSettings{
			LogLevel: string(mediasoupgo.Debug),
			LogTags:  []string{string(mediasoupgo.Message)},
		},
	)
	server, err := w.CreateWebRtcServer(
		&mediasoupgo.WebRtcServerOption{
			ListenInfos: []*mediasoupgo.TransportListenInfo{
				{
					Protocol:         mediasoupgo.UDPTransportProtocol,
					IP:               "0.0.0.0",
					AnnouncedAddress: ptr.String("127.0.0.1"),
					Port:             ptr.To[uint16](2233),
					PortRange:        nil,
					Flags:            nil,
					SendBufferSize:   ptr.To(uint32(1024 * 100)),
					RecvBufferSize:   ptr.To(uint32(1024 * 100)),
				},
			},
		})

	slog.Info("createWebrtcserver", "server", server, "error", err)
	router, err := w.CreateRouter(&mediasoupgo.RouterOption{
		MediaCodecs: []*mediasoupgo.RtpCodecCapability{
			{
				PreferredPayloadType: nil,
				Kind:                 mediasoupgo.AudioMediaKind,
				MimeType:             "audio/opus",
				ClockRate:            48000,
				Channels:             ptr.To[byte](2),
				RTCPFeedback: []*mediasoupgo.RtcpFeedback{
					{Type: "nack"},
					{Type: "transport-cc"},
				},
			},
		},
		AppData: mediasoupgo.RouterAppData{},
	})
	slog.Info("createRouter", "router", router, "err", err)
	// router
	webrtcTrans, err := router.CreateWebRtcTransport(
		&mediasoupgo.WebRtcTransportOptions{
			EnableUdp:                       ptr.To(true),
			EnableTcp:                       ptr.To(true),
			PreferUdp:                       nil,
			PreferTcp:                       nil,
			IceConsentTimeout:               nil,
			InitialAvailableOutgoingBitrate: nil,
			EnableSctp:                      ptr.To(true),
			NumSctpStreams:                  nil,
			MaxSctpMessageSize:              nil,
			SctpSendBufferSize:              nil,
			AppData:                         map[string]any{},
			WebRtcTransportListen: &mediasoupgo.WebRtcTransportListen{
				WebRtcServer: server,
			},
		},
	)

	fmt.Println(webrtcTrans, err)

	e := webrtcTrans.Connect(mediasoupgo.DtlsParameters{
		Role: ptr.To(mediasoupgo.ClientDtlsRole),
		Fingerprints: []mediasoupgo.DtlsFingerprint{{
			Algorithm: mediasoupgo.SHA256FingerprintAlgorithm,
			Value:     "A3:31:B9:80:7C:DD:36:0C:E3:8A:10:7C:D9:F0:40:33:80:EA:78:68:54:63:D4:6D:DA:42:05:CC:A7:B3:56:E8",
		}},
	})
	fmt.Println(e)
	producer, err := webrtcTrans.Produce(&mediasoupgo.ProducerOptions{
		ID:   nil,
		Kind: mediasoupgo.AudioMediaKind,
		RTPParameters: mediasoupgo.RtpParameters{
			MID: ptr.String("0"),
			Codecs: []*mediasoupgo.RtpCodecParameters{
				{
					PayloadType: 111,
					MimeType:    "audio/opus",
					ClockRate:   48000,
					Channels:    ptr.To[byte](2),
					Parameters: map[string]any{
						"minptime":     10,
						"useinbandfec": 1,
					},
					RTCPFeedback: []*mediasoupgo.RtcpFeedback{
						{Type: "nack"},
						{Type: "transport-cc"},
					},
				},
			},
			HeaderExtensions: []*mediasoupgo.RtpHeaderExtensionParameters{
				{
					URI:     mediasoupgo.SSRCAudioLevelRtpHeaderExtensionUri,
					ID:      1,
					Encrypt: ptr.To(false),
				},
				{
					URI:     mediasoupgo.TransportWideCC01RtpHeaderExtensionUri,
					ID:      2,
					Encrypt: ptr.To(false),
				},
				{
					URI:     mediasoupgo.AbsSendTimeRtpHeaderExtensionUri,
					ID:      3,
					Encrypt: ptr.Bool(false),
				},
				{
					URI:     mediasoupgo.MIDRtpHeaderExtensionUri,
					ID:      4,
					Encrypt: ptr.Bool(false),
				},
			},
			Encodings: []*mediasoupgo.RtpEncodingParameters{
				{
					SSRC:                  ptr.Uint32(1320542560),
					RID:                   ptr.String("0"),
					CodecPayloadType:      ptr.To[byte](111),
					RTX:                   &mediasoupgo.RTX{SSRC: 1320542560},
					DTX:                   ptr.Bool(false),
					ScalabilityMode:       nil,
					ScaleResolutionDownBy: nil,
					MaxBitrate:            ptr.Uint32(2312312312),
				},
			},
			RTCP: &mediasoupgo.RtcpParameters{
				CNAME:       ptr.String("0vcOi94DjwS/4A/x"),
				ReducedSize: ptr.Bool(true),
			},
		},
		AppData: map[string]any{},
	})
	fmt.Println(producer, err)
	time.Sleep(time.Hour)
}
