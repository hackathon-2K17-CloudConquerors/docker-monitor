package monitor

import (
	"time"

	"github.com/docker-monitor/influxdblib"
	"go.uber.org/zap"
)

func NewMonitor(smtpUser string, smtpPass string, smtpServer string, influxdb *influxdblib.Influxdb, dbname string, from string, to string) MonitorManipulator {

	return &Monitor{
		influx:   influxdb,
		dbname:   dbname,
		eventMap: make(map[string]struct{}),
		email:    newEmailProcessor(smtpUser, smtpPass, smtpServer, from, to),
	}
}

func (m *Monitor) StartMonitor(interval int) {
	for range time.Tick(time.Second * time.Duration(10)) {
		go m.checkAndSendEmailForNginxContainer()
		go m.checkAndSendEmailForHttpdContainer()
		go m.checkAndSendEmailForPostgresContainer()
	}
	return
}

func (m *Monitor) checkAndSendEmailForNginxContainer() {

	res, err := m.influx.ExecuteQuery(NginxLatestQuery, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	if len(res.Results[0].Series) > 0 {
		if res.Results[0].Series[0].Name == influxdblib.NginxMeasurementName {
			for _, nginxAttr := range res.Results[0].Series[0].Values {
				timestamp, _, containerName, event := m.extractValues(nginxAttr)
				key := containerName + event
				if _, ok := m.eventMap[key]; !ok {
					if event == influxdblib.ContainerStop {
						m.email.sendEmail(&emailcontent{
							ContainerID:   m.executeContainerQuery(NginxContainerIDQuery),
							ContainerName: containerName,
							Status:        event,
							time:          timestamp,
						})
						m.eventMap[key] = struct{}{}
					}
				}
			}
		}
	}
	return
}

func (m *Monitor) checkAndSendEmailForHttpdContainer() {

	res, err := m.influx.ExecuteQuery(HttpdLatestQuery, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	if len(res.Results[0].Series) > 0 {
		if res.Results[0].Series[0].Name == influxdblib.HttpdMeasurementName {
			for _, httpdAttr := range res.Results[0].Series[0].Values {
				timestamp, _, containerName, event := m.extractValues(httpdAttr)
				key := containerName + event
				if _, ok := m.eventMap[key]; !ok {
					if event == influxdblib.ContainerStop {

						m.email.sendEmail(&emailcontent{
							ContainerID:   m.executeContainerQuery(HttpdContainerIDQuery),
							ContainerName: containerName,
							Status:        event,
							time:          timestamp,
						})
						m.eventMap[key] = struct{}{}
					}
				}
			}
		}
	}
	return
}

func (m *Monitor) checkAndSendEmailForPostgresContainer() {
	res, err := m.influx.ExecuteQuery(PostgresLatestQuery, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	if len(res.Results[0].Series) > 0 {
		if res.Results[0].Series[0].Name == influxdblib.PostgresMeasurementName {
			for _, postgresAttr := range res.Results[0].Series[0].Values {
				timestamp, _, containerName, event := m.extractValues(postgresAttr)
				key := containerName + event
				if _, ok := m.eventMap[key]; !ok {
					if event == influxdblib.ContainerStop {
						m.email.sendEmail(&emailcontent{
							ContainerID:   m.executeContainerQuery(PostgresContainerIDQuery),
							ContainerName: containerName,
							Status:        event,
							time:          timestamp,
						})
						m.eventMap[key] = struct{}{}
					}
				}
			}
		}
	}
	return
}

func (m *Monitor) executeContainerQuery(query string) string {

	switch query {
	case NginxContainerIDQuery:
		m.retirieveContainerID(NginxContainerIDQuery)
	case HttpdContainerIDQuery:
		m.retirieveContainerID(HttpdContainerIDQuery)
	case PostgresContainerIDQuery:
		m.retirieveContainerID(PostgresContainerIDQuery)
	}

	return m.containerID
}

func (m *Monitor) retirieveContainerID(query string) {

	res, err := m.influx.ExecuteQuery(query, m.dbname)
	if err != nil {
		zap.L().Error("Error: Retrieving container events from DB", zap.Error(err))
	}

	for _, container := range res.Results[0].Series[0].Values {
		_, containerID, _, event := m.extractValues(container)
		if event == influxdblib.ContainerStart && containerID != "" {
			m.containerID = containerID
			return
		}
	}
}

func (m *Monitor) extractValues(container []interface{}) (string, string, string, string) {
	m.email.Lock()
	var timestamp, containerID, containerName, event string
	if value := container[0]; value != nil {
		timestamp = value.(string)
	}
	if value := container[1]; value != nil {
		containerID = value.(string)
	}
	if value := container[2]; value != nil {
		containerName = value.(string)
	}
	if value := container[4]; value != nil {
		event = value.(string)
	}
	m.email.Unlock()
	return timestamp, containerID, containerName, event
}
