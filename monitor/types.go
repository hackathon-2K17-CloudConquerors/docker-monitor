package monitor

import (
	"sync"
	"time"

	"github.com/docker-monitor/influxdblib"
	dockerClient "github.com/docker/docker/client"
	"github.com/domodwyer/mailyak"
)

// Monitor is holds data for monitoring
type Monitor struct {
	influx    *influxdblib.Influxdb
	dbname    string
	eventMap  map[string]struct{}
	email     *email
	container *influxdblib.Container
	docker    *dockerClient.Client
	sync.RWMutex
}

type email struct {
	smtpUser   string
	smtpPass   string
	smtpServer string
	to         string
	email      *mailyak.MailYak
	sync.Mutex
}

type emailcontent struct {
	time      time.Time
	container *influxdblib.Container
}
