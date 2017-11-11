package monitor

import (
	"sync"

	"github.com/docker-monitor/influxdblib"
	"github.com/domodwyer/mailyak"
)

type Monitor struct {
	influx      *influxdblib.Influxdb
	dbname      string
	eventMap    map[string]struct{}
	containerID string
	email       *email
}

type email struct {
	smtpUser   string
	smtpPass   string
	smtpServer string
	from       string
	to         string
	email      *mailyak.MailYak
	sync.Mutex
}

type emailcontent struct {
	ContainerID   string
	ImageName     string
	ContainerName string
	Status        string
	time          string
}
