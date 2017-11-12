package monitor

import "net/http"

// MonitorManipulator implements Monitor struct
type MonitorManipulator interface {
	StartMonitor(interval int)
	StartContainer(w http.ResponseWriter, r *http.Request)
}
