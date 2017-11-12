package monitor

import "net/http"

// MonitorManipulator implements Monitor struct
type MonitorManipulator interface {
	StartMonitor(interval int)
	ShowWebpage(w http.ResponseWriter, r *http.Request)
}
