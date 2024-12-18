package mediasoupgo

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/google/uuid"

	"mediasoupgo/FBS/Request"
	"mediasoupgo/FBS/Response"
	Transportfsb "mediasoupgo/FBS/Transport"
	"mediasoupgo/FBS/Worker"
)

type (
	WorkerObserverEvents struct{}
	WorkerObserver       func(*WorkerObserverEvents)
	WorkerAppData        struct{}
)

type CoreWorker struct {
	// mediasoup-worker child process.
	child *os.Process

	// Worker process PID.
	pid  int
	pids string

	// Channel instance.
	channel *Channel

	// Closed flag.
	closed bool

	// Died dlag.
	died bool

	// Worker subprocess closed flag.
	subprocessClosed bool

	// Custom app data.
	appData WorkerAppData

	// WebRtcServers set.
	webRtcServers map[*WebRtcServer]struct{}

	// Routers set.
	routers map[*Router]struct{}

	// Observer instance.
	observer WorkerObserver
}

func NewCoreWorker(logLevel string, logTags []string, rtcMinPort, rtcMaxPort uint16, dtlsCertificateFile, dtlsPrivateKeyFile, libwebrtcFieldTrials string, disableLiburing bool) *CoreWorker {
	workerBin := os.Getenv("MEDIASOUP_WORKER_BIN")

	cmd := exec.Command(workerBin)
	if logLevel != "" {
		cmd.Args = append(cmd.Args, "--logLevel="+logLevel)
	}
	for _, logTag := range logTags {
		cmd.Args = append(cmd.Args, "--logTag="+logTag)
	}
	if rtcMinPort != 0 {
		cmd.Args = append(cmd.Args, "--rtcMinPort="+fmt.Sprintf("%d", rtcMinPort))
	}
	if rtcMaxPort != 0 {
		cmd.Args = append(cmd.Args, "--rtcMaxPort="+fmt.Sprintf("%d", rtcMaxPort))
	}
	if dtlsCertificateFile != "" && dtlsPrivateKeyFile != "" {
		cmd.Args = append(cmd.Args, "--dtlsPrivateKeyFile="+dtlsPrivateKeyFile, "--dtlsCertificateFile="+dtlsCertificateFile)
	}
	if libwebrtcFieldTrials != "" {
		cmd.Args = append(cmd.Args, "--libwebrtcFieldTrials="+libwebrtcFieldTrials)
	}
	if disableLiburing {
		cmd.Args = append(cmd.Args, "disableLiburing=true")
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
	cmd.Env = []string{"MEDIASOUP_VERSION=" + "3.15.2"}
	if err := cmd.Start(); err != nil {
		return nil
	}
	go func() {
		cmd.Wait()
	}()
	select {
	case <-spawnDone:
		break
	case <-time.After(time.Second * 10):
		return nil
	}
	channel.SetPid(cmd.Process.Pid)
	w := &CoreWorker{
		child:            cmd.Process,
		pid:              cmd.Process.Pid,
		pids:             strconv.Itoa(cmd.Process.Pid),
		channel:          channel,
		closed:           false,
		died:             false,
		subprocessClosed: false,
		appData:          WorkerAppData{},
		webRtcServers:    map[*WebRtcServer]struct{}{},
		routers:          map[*Router]struct{}{},
		observer:         func(*WorkerObserverEvents) {},
	}
	return w
}

func (w *CoreWorker) Close() {}

func (w *CoreWorker) Dump() (*Response.ResponseT, error) {
	return w.channel.Request(Request.MethodWORKER_DUMP, &Request.BodyT{Type: Request.BodyNONE}, w.pids)
}

func (w *CoreWorker) GetResourceUsage() (*Response.ResponseT, error) {
	return w.channel.Request(Request.MethodWORKER_GET_RESOURCE_USAGE, &Request.BodyT{Type: Request.BodyNONE}, w.pids)
}

func (w *CoreWorker) UpdateSettings() {}

func (w *CoreWorker) CreateRouter() {}

func (w *CoreWorker) CreateWebRtcServer(listenInfos []*Transportfsb.ListenInfoT) (*WebRtcServer, error) {
	id := uuid.NewString()
	body := &Request.BodyT{Type: Request.BodyWorker_CreateWebRtcServerRequest, Value: &Worker.CreateWebRtcServerRequestT{WebRtcServerId: id, ListenInfos: listenInfos}}
	_, err := w.channel.Request(Request.MethodWORKER_CREATE_WEBRTCSERVER, body, w.pids)
	if err != nil {
		return nil, err
	}
	return &WebRtcServer{WebRtcServerId: id}, nil
}

func (w *CoreWorker) WrokerDied() {}
