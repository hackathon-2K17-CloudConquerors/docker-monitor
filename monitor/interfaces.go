package monitor

type MonitorManipulator interface {
	StartMonitor(interval int)
}
