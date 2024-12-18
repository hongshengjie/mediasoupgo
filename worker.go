package meidsoupgo

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/kataras/go-events"
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
	pid int

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

	cmd.Env = []string{"MEDIASOUP_VERSION=" + "3.15.2"}
	cmd.Run()
	go func() {
		cmd.Wait()
	}()
	w := &CoreWorker{
		child:            cmd.Process,
		pid:              cmd.Process.Pid,
		channel:          NewChannel(producerWriter, consumerReader, cmd.Process.Pid),
		closed:           false,
		died:             false,
		subprocessClosed: false,
		appData:          WorkerAppData{},
		webRtcServers:    map[*WebRtcServer]struct{}{},
		routers:          map[*Router]struct{}{},
		observer: func(*WorkerObserverEvents) {
		},
	}
	w.channel.Once(events.EventName(strconv.Itoa(w.pid)), func(i ...interface{}) {})
	return w
}
func (w *CoreWorker) Close()              {}
func (w *CoreWorker) Dump()               {}
func (w *CoreWorker) GetResourceUsage()   {}
func (w *CoreWorker) UpdateSettings()     {}
func (w *CoreWorker) CreateRouter()       {}
func (w *CoreWorker) CreateWebRtcServer() {}
func (w *CoreWorker) WrokerDied()         {}

