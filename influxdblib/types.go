package influxdblib

import (
	"sync"

	"github.com/influxdata/influxdb/client/v2"
)

//Influxdb inplements a DataAdder interface for influxDB
type Influxdb struct {
	httpClient client.Client
	database   string

	stopWorker chan struct{}
	worker     *worker
	sync.Mutex
}

// Container holds data of the containers
type Container struct {
	ContainerID      string
	ContainerName    string
	ContainerStatus  string
	ContainerCreated int64
	ContainerNetwork string
	ContainerState   string
	ImageName        string
	Status           string
}

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
