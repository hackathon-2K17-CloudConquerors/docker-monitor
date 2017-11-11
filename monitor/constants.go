package monitor

const (
	NginxLatestQuery = "SELECT * FROM NginxContainerEvents ORDER BY DESC LIMIT 1"

	NginxContainerIDQuery = "SELECT * FROM NginxContainerEvents"

	HttpdLatestQuery = "SELECT * FROM HttpdContainerEvents ORDER BY DESC LIMIT 1"

	HttpdContainerIDQuery = "SELECT * FROM HttpdContainerEvents"

	PostgresLatestQuery = "SELECT * FROM PostgresContainerEvents ORDER BY DESC LIMIT 1"

	PostgresContainerIDQuery = "SELECT * FROM PostgresContainerEvents"
)
