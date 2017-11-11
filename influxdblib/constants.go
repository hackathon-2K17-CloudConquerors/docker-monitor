package influxdblib

const (
	// NGINX conatiner name
	NGINX = "nginx"
	// HTTPD container name
	HTTPD = "httpd"
	// POSTGRES container name
	POSTGRES = "postgres"

	UNKNOWN = "unknown"
)

const (
	ContainerStart = "start"
	ContainerStop  = "stop"
)

const (
	NginxMeasurementName    = "NginxContainerEvents"
	HttpdMeasurementName    = "HttpdContainerEvents"
	PostgresMeasurementName = "PostgresContainerEvents"
)
