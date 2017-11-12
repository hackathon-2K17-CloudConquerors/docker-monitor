package monitor

const (
	// NginxLatestQuery to get latest entry from DB
	NginxLatestQuery = "SELECT * FROM NginxContainerEvents ORDER BY DESC LIMIT 1"
	// NginxContainerIDQuery to get all entries from DB
	NginxContainerIDQuery = "SELECT * FROM NginxContainerEvents"
	// HttpdLatestQuery to get latest entry from DB
	HttpdLatestQuery = "SELECT * FROM HttpdContainerEvents ORDER BY DESC LIMIT 1"
	// HttpdContainerIDQuery to get all entries from DB
	HttpdContainerIDQuery = "SELECT * FROM HttpdContainerEvents"
	// PostgresLatestQuery to get latest entry from DB
	PostgresLatestQuery = "SELECT * FROM PostgresContainerEvents ORDER BY DESC LIMIT 1"
	// PostgresContainerIDQuery to get all entries from DB
	PostgresContainerIDQuery = "SELECT * FROM PostgresContainerEvents"
)

const (
	// DefaultLocalhost server endpoint
	DefaultLocalhost = "http://localhost:8088/monitor?containerid="
)
