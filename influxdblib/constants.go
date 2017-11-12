package influxdblib

const (
	// NGINX conatiner name
	NGINX = "nginx"
	// HTTPD container name
	HTTPD = "httpd"
	// POSTGRES container name
	POSTGRES = "postgres"
	// UNKNOWN container name
	UNKNOWN = "unknown"
)

const (
	// ContainerStart is start event
	ContainerStart = "start"
	// ContainerStop is start event
	ContainerStop = "stop"
)

const (
	// NginxMeasurementName is a DB measurement name
	NginxMeasurementName = "NginxContainerEvents"
	// HttpdMeasurementName is a DB measurement name
	HttpdMeasurementName = "HttpdContainerEvents"
	// PostgresMeasurementName is a DB measurement name
	PostgresMeasurementName = "PostgresContainerEvents"
)
