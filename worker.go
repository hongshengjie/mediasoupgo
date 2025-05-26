package mediasoupgo

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"

	"mediasoupgo/FBS/Request"
	transport "mediasoupgo/FBS/Transport"
	worker "mediasoupgo/FBS/Worker"
	"mediasoupgo/events"
	"mediasoupgo/smap"
)

var _ Worker = &workerImpl{}

type workerImpl struct {
	child            *os.Process
	pid              int
	pids             string
	channel          *Channel
	closed           atomic.Bool
	died             atomic.Bool
	subprocessClosed atomic.Bool
	appData          WorkerAppData
	webRtcServers    *smap.Map[WebRtcServer, struct{}]
	routers          *smap.Map[Router, struct{}]
	observer         WorkerObserver
	events.EventEmmiter[WorkerEvents]
}

func NewWorker(settings *WorkerSettings) Worker {
	logLevel := settings.LogLevel
	logTags := settings.LogTags
	rtcMinPort := settings.RTCMinPort
	rtcMaxPort := settings.RTCMaxPort
	dtlsCertificateFile := settings.DTLSCertificateFile
	dtlsPrivateKeyFile := settings.DTLSPrivateKeyFile
	libwebrtcFieldTrials := settings.LibwebrtcFieldTrials
	disableLiburing := settings.DisableLiburing

	workerBin := os.Getenv("MEDIASOUP_WORKER_BIN")
	cmd := exec.Command(workerBin)
	if logLevel != nil {
		cmd.Args = append(cmd.Args, "--logLevel="+string(*logLevel))
	}
	for _, logTag := range logTags {
		cmd.Args = append(cmd.Args, "--logTag="+string(logTag))
	}
	if rtcMinPort != nil {
		cmd.Args = append(cmd.Args, "--rtcMinPort="+fmt.Sprintf("%d", *rtcMinPort))
	}
	if rtcMaxPort != nil {
		cmd.Args = append(cmd.Args, "--rtcMaxPort="+fmt.Sprintf("%d", *rtcMaxPort))
	}
	if dtlsCertificateFile != nil && dtlsPrivateKeyFile != nil {
		cmd.Args = append(
			cmd.Args,
			"--dtlsPrivateKeyFile="+*dtlsPrivateKeyFile,
			"--dtlsCertificateFile="+*dtlsCertificateFile,
		)
	}
	if libwebrtcFieldTrials != nil {
		cmd.Args = append(cmd.Args, "--libwebrtcFieldTrials="+*libwebrtcFieldTrials)
	}
	if disableLiburing != nil {
		cmd.Args = append(cmd.Args, "disableLiburing="+fmt.Sprintf("%v", *disableLiburing))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil
	producerReader, producerWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	consumerReader, consumerWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	cmd.ExtraFiles = []*os.File{producerReader, consumerWriter}
	channel, spawnDone := NewChannel(producerWriter, consumerReader)

	cmd.Env = []string{"MEDIASOUP_VERSION=" + os.Getenv("MEDIASOUP_VERSION")}
	slog.Debug("cmd", "args", cmd.Args, "env", cmd.Env)
	if err := cmd.Start(); err != nil {
		return nil
	}

	w := &workerImpl{
		channel:          channel,
		closed:           atomic.Bool{},
		died:             atomic.Bool{},
		subprocessClosed: atomic.Bool{},
		appData:          WorkerAppData{},
		webRtcServers:    smap.New[WebRtcServer, struct{}](),
		routers:          smap.New[Router, struct{}](),
		observer:         events.New[WorkerObserverEvents](),
		EventEmmiter:     events.New[WorkerEvents](),
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			code := cmd.ProcessState.ExitCode()
			var started bool
			select {
			case _, ok := <-spawnDone:
				started = !ok
			default:
				started = false
			}
			if !started {
				if code == 42 {
					w.Close()
					w.Emit("@failure", WorkerEvents{AtFailure: events.NewEvent1(errors.New(""))})
				} else {
					w.Close()
					w.Emit("@failure", WorkerEvents{AtFailure: events.NewEvent1(errors.New(""))})
				}
			} else {
				w.WorkerDied(errors.New(""))
			}
		}
	}()
	select {
	case <-spawnDone:

		channel.SetPid(cmd.Process.Pid)
		w.pid = cmd.Process.Pid
		w.pids = strconv.Itoa(cmd.Process.Pid)
		w.child = cmd.Process
		w.Emit("@success", WorkerEvents{AtSuccess: struct{}{}})
		break
		// case <-time.After(time.Second * 10):
		// 	return nil
	}
	w.handleListernerError()
	return w
}

func (w *workerImpl) Observer() WorkerObserver {
	return w.observer
}

func (w *workerImpl) SubprocessClosed() bool {
	return w.subprocessClosed.Load()
}

func (w *workerImpl) Pid() int {
	return w.pid
}

func (w *workerImpl) Died() bool {
	return w.died.Load()
}

func (w *workerImpl) WorkerDied(err error) {
	if w.closed.Load() {
		return
	}
	w.died.Store(true)
	w.closed.Store(true)

	w.channel.Close()
	w.routers.Range(func(key Router, value struct{}) bool {
		key.WorkerClosed()
		return true
	})
	w.webRtcServers.Range(func(key WebRtcServer, value struct{}) bool {
		key.WorkerClosed()
		return true
	})
	w.Emit("died", WorkerEvents{Died: events.NewEvent1(err)})
	w.observer.Emit("close", WorkerObserverEvents{Close: struct{}{}})
}

func (w *workerImpl) Close() {
	if w.closed.Load() {
		return
	}
	w.closed.Store(true)
	w.child.Kill()
	w.channel.Close()

	w.routers.Range(func(key Router, value struct{}) bool {
		key.WorkerClosed()
		return true
	})
	w.webRtcServers.Range(func(key WebRtcServer, value struct{}) bool {
		key.WorkerClosed()
		return true
	})
	w.observer.Emit("close", WorkerObserverEvents{Close: struct{}{}})
}

func (w *workerImpl) Closed() bool {
	return w.closed.Load()
}

func (w *workerImpl) AppData() WorkerAppData {
	return w.appData
}

func (w *workerImpl) SetAppData(data WorkerAppData) {
	w.appData = data
}

func (w *workerImpl) Dump() (*WorkerDump, error) {
	resp, err := w.channel.Request(
		Request.MethodWORKER_DUMP,
		&Request.BodyT{Type: Request.BodyNONE},
		w.pids,
	)
	if err != nil {
		return nil, err
	}
	dump := resp.Body.Value.(*worker.DumpResponseT)
	ret := &WorkerDump{
		PID:             dump.Pid,
		WebRTCServerIDs: dump.WebRtcServerIds,
		RouterIDs:       dump.RouterIds,
	}
	if dump.ChannelMessageHandlers != nil {
		ret.ChannelMessageHandlers = struct {
			ChannelRequestHandlers      []string
			ChannelNotificationHandlers []string
		}{
			ChannelRequestHandlers:      dump.ChannelMessageHandlers.ChannelRequestHandlers,
			ChannelNotificationHandlers: dump.ChannelMessageHandlers.ChannelNotificationHandlers,
		}
	}
	if dump.Liburing != nil {
		ret.Liburing = &struct {
			SQEProcessCount   uint64
			SQEMissCount      uint64
			UserDataMissCount uint64
		}{
			SQEProcessCount:   (dump.Liburing.SqeProcessCount),
			SQEMissCount:      (dump.Liburing.SqeMissCount),
			UserDataMissCount: (dump.Liburing.UserDataMissCount),
		}
	}
	return ret, nil
}

func (w *workerImpl) GetResourceUsage() (*WorkerResourceUsage, error) {
	resp, err := w.channel.Request(
		Request.MethodWORKER_GET_RESOURCE_USAGE,
		&Request.BodyT{Type: Request.BodyNONE},
		w.pids,
	)
	if err != nil {
		return nil, err
	}
	u := resp.Body.Value.(*worker.ResourceUsageResponseT)
	ret := &WorkerResourceUsage{
		RUUtime:    u.RuUtime,
		RUStime:    u.RuStime,
		RUMaxrss:   u.RuMaxrss,
		RUIxrss:    u.RuIxrss,
		RUIdrss:    u.RuIdrss,
		RUIsrss:    u.RuIsrss,
		RUMinflt:   u.RuMinflt,
		RUMajflt:   u.RuMajflt,
		RUNswap:    u.RuNswap,
		RUInblock:  u.RuInblock,
		RUOublock:  u.RuOublock,
		RUMsgsnd:   u.RuMsgsnd,
		RUMsgrcv:   u.RuMsgrcv,
		RUNsignals: u.RuNsignals,
		RUNvcsw:    u.RuNvcsw,
		RUNivcsw:   u.RuNivcsw,
	}
	return ret, nil
}

func (w *workerImpl) UpdateSettings(req *WorkerUpdateableSettings) error {
	_, err := w.channel.Request(
		Request.MethodWORKER_UPDATE_SETTINGS,
		&Request.BodyT{
			Type: Request.BodyWorker_UpdateSettingsRequest,
			Value: &worker.UpdateSettingsRequestT{
				LogLevel: req.LogLevel,
				LogTags:  req.LogTags,
			},
		},
		w.pids,
	)
	return err
}

func (w *workerImpl) CreateWebRtcServer(req *WebRtcServerOption) (WebRtcServer, error) {
	id := uuid.NewString()
	var listenInfos []*transport.ListenInfoT
	for _, v := range req.ListenInfos {
		var addr string
		if v.AnnouncedIP != nil {
			addr = *v.AnnouncedIP
		}
		if v.AnnouncedAddress != nil {
			addr = *v.AnnouncedAddress
		}
		info := &transport.ListenInfoT{
			Protocol:         transport.EnumValuesProtocol[strings.ToUpper(string(v.Protocol))],
			Ip:               v.IP,
			AnnouncedAddress: addr,
			PortRange:        &transport.PortRangeT{},
			Flags:            &transport.SocketFlagsT{},
		}
		if v.Port != nil {
			info.Port = *v.Port
		}

		if v.PortRange != nil {
			info.PortRange = &transport.PortRangeT{Min: v.PortRange.Min, Max: v.PortRange.Max}
		}
		if v.Flags != nil {
			info.Flags = &transport.SocketFlagsT{
				Ipv6Only:     v.Flags.IPV6Only,
				UdpReusePort: v.Flags.UDPReusePort,
			}
		}
		if v.RecvBufferSize != nil {
			info.RecvBufferSize = *v.RecvBufferSize
		}
		if v.SendBufferSize != nil {
			info.SendBufferSize = *v.SendBufferSize
		}
		listenInfos = append(listenInfos, info)
	}
	r := &worker.CreateWebRtcServerRequestT{WebRtcServerId: id, ListenInfos: listenInfos}
	body := &Request.BodyT{Type: Request.BodyWorker_CreateWebRtcServerRequest, Value: r}
	_, err := w.channel.Request(Request.MethodWORKER_CREATE_WEBRTCSERVER, body, w.pids)
	if err != nil {
		return nil, err
	}
	s := NewWebRtcServer(id, w.channel, req.AppData)
	s.On("@close", func(arg WebRtcServerEvents) {
		w.webRtcServers.Delete(s)
	})
	w.webRtcServers.Set(s, struct{}{})
	w.Observer().Emit("newwebrtcserver", WorkerObserverEvents{Newwebrtcserver: events.NewEvent1(s)})
	return s, nil
}

func (w *workerImpl) CreateRouter(req *RouterOption) (Router, error) {
	routerId := uuid.NewString()
	body := &Request.BodyT{
		Type:  Request.BodyWorker_CreateRouterRequest,
		Value: &worker.CreateRouterRequestT{RouterId: routerId},
	}
	_, err := w.channel.Request(Request.MethodWORKER_CREATE_ROUTER, body, w.pids)
	if err != nil {
		return nil, err
	}
	caps, err := GenerateRouterRtpCapabilities(req.MediaCodecs)
	if err != nil {
		return nil, err
	}
	r, err := NewRouter(routerId, caps, w.channel, RouterAppData(w.appData))
	if err != nil {
		return nil, err
	}
	w.routers.Set(r, struct{}{})
	r.On("@close", func(arg RouterEvents) {
		w.routers.Delete(r)
	})
	w.observer.Emit("newrouter", WorkerObserverEvents{Newrouter: events.NewEvent1(r)})
	return r, nil
}

func (w *workerImpl) handleListernerError() {
	w.On("listenererror", func(arg WorkerEvents) {
		// TODO
	})
}
