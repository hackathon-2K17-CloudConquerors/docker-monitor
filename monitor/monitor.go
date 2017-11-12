package monitor

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker-monitor/influxdblib"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"go.uber.org/zap"
)

// NewMonitor is a handler for processing emails and monitoring
func NewMonitor(smtpUser string, smtpPass string, smtpServer string, docker *dockerClient.Client, influxdb *influxdblib.Influxdb, dbname string, to string) MonitorManipulator {

	return &Monitor{
		influx:   influxdb,
		dbname:   dbname,
		docker:   docker,
		eventMap: make(map[string]struct{}),
		email:    newEmailProcessor(smtpUser, smtpPass, smtpServer, to),
	}
}

// StartContainer is a server function from a web service
func (m *Monitor) StartContainer(w http.ResponseWriter, r *http.Request) {

	containerID := r.FormValue("containerid")
	zap.L().Info("Starting Container", zap.Any("container", containerID))
	if err := m.docker.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{}); err != nil {
		http.Error(w, err.Error(), 3)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Successfully started container "+containerID)
}

// StartMonitor starts goroutines to qurey DB and to listen for containers
func (m *Monitor) StartMonitor(interval int) {
	zap.L().Info("Started Monitoring", zap.Any("interval", interval))
	for range time.Tick(time.Second * time.Duration(interval)) {
		go m.checkAndSendEmailForNginxContainer()
		go m.checkAndSendEmailForHttpdContainer()
		go m.checkAndSendEmailForPostgresContainer()
	}
	return
}

func (m *Monitor) checkAndSendEmailForNginxContainer() {
	m.Lock()
	defer m.Unlock()
	res, err := m.influx.ExecuteQuery(NginxLatestQuery, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	if len(res.Results[0].Series) > 0 {
		if res.Results[0].Series[0].Name == influxdblib.NginxMeasurementName {
			for _, nginxAttr := range res.Results[0].Series[0].Values {
				timestamp, allContAttr := m.extractValues(nginxAttr)
				key := allContAttr.ContainerName + allContAttr.Status
				if _, ok := m.eventMap[key]; !ok {
					if allContAttr.Status == influxdblib.ContainerStop {
						zap.L().Info("Container Stopped...Sending email", zap.Any("container", allContAttr.ContainerName))
						containerAttr, err := m.executeContainerQuery(NginxContainerIDQuery)
						if err != nil {
							zap.L().Error("Error: Empty Attributes", zap.Error(err))
						} else {
							m.email.sendEmail(&emailcontent{container: &influxdblib.Container{
								ContainerID:      containerAttr.ContainerID,
								ContainerName:    containerAttr.ContainerName,
								Status:           containerAttr.Status,
								ContainerNetwork: containerAttr.ContainerNetwork,
								ContainerState:   containerAttr.ContainerState,
								ContainerStatus:  containerAttr.ContainerStatus,
								ImageName:        containerAttr.ImageName,
							},
								time: timestamp,
							})
							m.eventMap[key] = struct{}{}
						}
					} else {
						if len(m.eventMap) > 0 {
							key := allContAttr.ContainerName + influxdblib.ContainerStop
							delete(m.eventMap, key)
						}
					}
				}
			}
		}
	}
	return
}

func (m *Monitor) checkAndSendEmailForHttpdContainer() {
	m.Lock()
	defer m.Unlock()
	res, err := m.influx.ExecuteQuery(HttpdLatestQuery, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	if len(res.Results[0].Series) > 0 {
		if res.Results[0].Series[0].Name == influxdblib.HttpdMeasurementName {
			for _, httpdAttr := range res.Results[0].Series[0].Values {
				timestamp, allContAttr := m.extractValues(httpdAttr)
				key := allContAttr.ContainerName + allContAttr.Status
				if _, ok := m.eventMap[key]; !ok {
					if allContAttr.Status == influxdblib.ContainerStop {
						zap.L().Info("Container Stopped...Sending email", zap.Any("container", allContAttr.ContainerName))
						containerAttr, err := m.executeContainerQuery(HttpdContainerIDQuery)
						if err != nil {
							zap.L().Error("Error: Empty Attributes", zap.Error(err))
						} else {
							m.email.sendEmail(&emailcontent{container: &influxdblib.Container{
								ContainerID:      containerAttr.ContainerID,
								ContainerName:    containerAttr.ContainerName,
								Status:           containerAttr.Status,
								ContainerNetwork: containerAttr.ContainerNetwork,
								ContainerState:   containerAttr.ContainerState,
								ContainerStatus:  containerAttr.ContainerStatus,
								ImageName:        containerAttr.ImageName,
							},
								time: timestamp,
							})
							m.eventMap[key] = struct{}{}
						}
					} else {
						if len(m.eventMap) > 0 {
							key := allContAttr.ContainerName + influxdblib.ContainerStop
							delete(m.eventMap, key)
						}
					}
				}
			}
		}
	}
	return
}

func (m *Monitor) checkAndSendEmailForPostgresContainer() {
	m.Lock()
	defer m.Unlock()
	res, err := m.influx.ExecuteQuery(PostgresLatestQuery, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	if len(res.Results[0].Series) > 0 {
		if res.Results[0].Series[0].Name == influxdblib.PostgresMeasurementName {
			for _, postgresAttr := range res.Results[0].Series[0].Values {
				timestamp, allContAttr := m.extractValues(postgresAttr)
				key := allContAttr.ContainerName + allContAttr.Status
				if _, ok := m.eventMap[key]; !ok {
					if allContAttr.Status == influxdblib.ContainerStop {
						zap.L().Info("Container Stopped...Sending email", zap.Any("container", allContAttr.ContainerName))
						containerAttr, err := m.executeContainerQuery(PostgresContainerIDQuery)
						if err != nil {
							zap.L().Error("Error: Empty Attributes", zap.Error(err))
						} else {
							m.email.sendEmail(&emailcontent{container: &influxdblib.Container{
								ContainerID:      containerAttr.ContainerID,
								ContainerName:    containerAttr.ContainerName,
								Status:           containerAttr.Status,
								ContainerNetwork: containerAttr.ContainerNetwork,
								ContainerState:   containerAttr.ContainerState,
								ContainerStatus:  containerAttr.ContainerStatus,
								ImageName:        containerAttr.ImageName,
							},
								time: timestamp,
							})
							m.eventMap[key] = struct{}{}
						}
					} else {
						if len(m.eventMap) > 0 {
							key := allContAttr.ContainerName + influxdblib.ContainerStop
							delete(m.eventMap, key)
						}
					}
				}
			}
		}
	}
	return
}

func (m *Monitor) executeContainerQuery(query string) (*influxdblib.Container, error) {

	var err error
	switch query {
	case NginxContainerIDQuery:
		err = m.retirieveContainerID(NginxContainerIDQuery)
	case HttpdContainerIDQuery:
		err = m.retirieveContainerID(HttpdContainerIDQuery)
	case PostgresContainerIDQuery:
		err = m.retirieveContainerID(PostgresContainerIDQuery)
	}

	return m.container, err
}

func (m *Monitor) retirieveContainerID(query string) error {

	res, err := m.influx.ExecuteQuery(query, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	for _, container := range res.Results[0].Series[0].Values {
		_, contAttr := m.extractValues(container)
		if contAttr.Status == influxdblib.ContainerStart && contAttr.ContainerID != "" {
			m.container = &influxdblib.Container{
				ContainerID:      contAttr.ContainerID,
				ContainerName:    contAttr.ContainerName,
				Status:           contAttr.Status,
				ContainerNetwork: contAttr.ContainerNetwork,
				ContainerState:   contAttr.ContainerState,
				ContainerStatus:  contAttr.ContainerStatus,
				ImageName:        contAttr.ImageName,
			}
			return nil
		}
	}
	return fmt.Errorf("Error: Container stopped before processing")
}

func (m *Monitor) extractValues(container []interface{}) (time.Time, *influxdblib.Container) {

	var contAttr influxdblib.Container
	var timestamp string
	if value := container[0]; value != nil {
		timestamp = value.(string)
	}
	if value := container[2]; value != nil {
		contAttr.ContainerID = value.(string)
	}
	if value := container[3]; value != nil {
		contAttr.ContainerName = value.(string)
	}
	if value := container[4]; value != nil {
		contAttr.ContainerNetwork = value.(string)
	}
	if value := container[5]; value != nil {
		contAttr.ContainerState = value.(string)
	}
	if value := container[6]; value != nil {
		contAttr.ContainerStatus = value.(string)
	}
	if value := container[7]; value != nil {
		contAttr.ImageName = value.(string)
	}
	if value := container[8]; value != nil {
		contAttr.Status = value.(string)
	}
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		zap.L().Error("Error: Parsing time", zap.Error(err))
	}

	return parsedTime, &contAttr
}
