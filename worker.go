package meidsoupgo

import (
	"os"
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
