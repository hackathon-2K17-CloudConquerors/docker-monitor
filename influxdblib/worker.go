package influxdblib

import (
	"fmt"

	"go.uber.org/zap"
)

// A worker manages the workload for the InfluxDB collector
type worker struct {
	events chan *workerEvent
	stop   chan struct{}
	db     DataAdder
}

// a workerEvent is an event that the worker need to process
type workerEvent struct {
	containerName string
	container     *Container
}

func newWorker(stop chan struct{}, db DataAdder) *worker {
	return &worker{
		events: make(chan *workerEvent, 500),
		stop:   stop,
		db:     db,
	}
}

func (w *worker) addEvent(wevent *workerEvent) {
	select {
	case w.events <- wevent: // Put event in channel unless it is full
		zap.L().Debug("Adding event to InfluxDBProcessingQueue.")
	default:
		zap.L().Warn("Event queue full for InfluxDB. Dropping event.")
	}
}

// startWorker start processing the event for this worker.
// Blocking... Use go.
func (w *worker) startWorker() {
	zap.L().Info("Starting InfluxDBworker")
	for {
		select {
		case event := <-w.events:
			w.processEvent(event)
		case <-w.stop:
			return
		}
	}
}

func (w *worker) processEvent(wevent *workerEvent) {
	zap.L().Debug("Processing event for InfluxDB")
	switch wevent.containerName {
	case NGINX:
		if err := w.doCollectContainerEvent(NGINX, wevent.container); err != nil {
			zap.L().Error("Couldn't process influxDB Request ContainerRequest", zap.Error(err))
		}
	case HTTPD:
		if err := w.doCollectContainerEvent(HTTPD, wevent.container); err != nil {
			zap.L().Error("Couldn't process influxDB Request FlowRequest", zap.Error(err))
		}
	case POSTGRES:
		if err := w.doCollectContainerEvent(POSTGRES, wevent.container); err != nil {
			zap.L().Error("Couldn't process influxDB Request ContainerRequest", zap.Error(err))
		}
	default:
		if err := w.doCollectContainerEvent(UNKNOWN, wevent.container); err != nil {
			zap.L().Error("Couldn't process influxDB Request ContainerRequest", zap.Error(err))
		}
	}
}

// CollectContainerEvent implements trireme collector interface
func (w *worker) doCollectContainerEvent(containerType string, container *Container) error {
	var containerName string

	switch containerType {
	case NGINX:
		containerName = NGINX
	case HTTPD:
		containerName = HTTPD
	case POSTGRES:
		containerName = POSTGRES
	case UNKNOWN:
		containerName = UNKNOWN
	default:
		return fmt.Errorf("Unrecognized container event name %s ", containerType)
	}

	return w.db.AddData(map[string]string{
		"ContainerName": containerName,
	}, map[string]interface{}{
		"ContainerID": container.ContainerID,
		"ImageName":   container.ImageName,
		"Status":      container.Status,
	})
}
