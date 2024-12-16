package meidsoupgo

import (
	"fmt"
	"os"
	"os/exec"
)

type WebRtcServer struct{}

type Router struct{}

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

func NewCoreWorker(logLevel string, logTags []string, rtcMinPort, rtcMaxPort uint16, dtlsCertificateFile, dtlsPrivateKeyFile, libwebrtcFieldTrials string, disableLiburing bool) {
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
}
