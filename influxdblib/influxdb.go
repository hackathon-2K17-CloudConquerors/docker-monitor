package influxdblib

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

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

//DataAdder interface has all the methods required to interact with influxdb api
type DataAdder interface {
	CreateDB(string) error
	AddData(tags map[string]string, fields map[string]interface{}) error
	ExecuteQuery(query string, dbname string) (*client.Response, error)
}

// NewDBConnection is used to create a new client and return influxdb handle
func NewDBConnection(user string, pass string, addr string, db string, insecureSkipVerify bool) (*Influxdb, error) {
	zap.L().Debug("Initializing InfluxDBConnection")
	httpClient, err := createHTTPClient(user, pass, addr, insecureSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("Error parsing url %s", err)
	}
	_, _, err = httpClient.Ping(time.Second * 0)
	if err != nil {
		return nil, fmt.Errorf("Unable to create InfluxDB http client %s", err)
	}

	dbConnection := &Influxdb{
		httpClient: httpClient,
		database:   db,
		stopWorker: make(chan struct{}),
	}

	worker := newWorker(dbConnection.stopWorker, dbConnection)
	dbConnection.worker = worker

	// Attempt to create the Database. Silently fail if it already exists.
	if err := dbConnection.CreateDB(db); err != nil {
		return nil, fmt.Errorf("Error: Creating Database: %s", err)
	}

	return dbConnection, nil
}

func createHTTPClient(user string, pass string, addr string, InsecureSkipVerify bool) (client.Client, error) {

	httpClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               addr,
		Username:           user,
		Password:           pass,
		Timeout:            time.Minute,
		InsecureSkipVerify: InsecureSkipVerify,
	})
	if err != nil {
		return nil, err
	}

	return httpClient, nil
}

// CreateDB is used to create a new databases given name
func (d *Influxdb) CreateDB(dbname string) error {
	zap.L().Info("Creating database", zap.String("db", dbname))

	_, err := d.ExecuteQuery("CREATE DATABASE "+dbname, "")
	if err != nil {
		return err
	}

	return nil
}

// ExecuteQuery is used to execute a query given a database name
func (d *Influxdb) ExecuteQuery(query string, dbname string) (*client.Response, error) {

	q := client.Query{
		Command:  query,
		Database: dbname,
		Chunked:  false,
	}

	response, err := d.httpClient.Query(q)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Start is used to start listening for data
func (d *Influxdb) Start() error {
	zap.L().Info("Starting InfluxDB worker")

	go d.worker.startWorker()

	return nil
}

// Stop is used to stop and return from listen goroutine
func (d *Influxdb) Stop() error {
	zap.L().Info("Stopping InfluxDB worker")

	d.stopWorker <- struct{}{}
	d.httpClient.Close()

	return nil
}

// AddData is used to add data to the batch
func (d *Influxdb) AddData(tags map[string]string, fields map[string]interface{}) error {
	zap.L().Debug("Calling AddData", zap.Any("tags", tags), zap.Any("fields", fields))
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  d.database,
		Precision: "us",
	})
	if err != nil {
		return fmt.Errorf("Couldn't add data, error creating batchpoint: %s", err)
	}

	switch tags["ContainerName"] {
	case NGINX:
		pt, err := client.NewPoint("NginxContainerEvents", tags, fields, time.Now())
		if err != nil {
			return fmt.Errorf("Couldn't add ContainerEvent: %s", err)
		}
		bp.AddPoint(pt)
	case HTTPD:
		pt, err := client.NewPoint("HttpdContainerEvents", tags, fields, time.Now())
		if err != nil {
			return fmt.Errorf("Couldn't add FlowEvent: %s", err)
		}
		bp.AddPoint(pt)
	case POSTGRES:
		pt, err := client.NewPoint("PostgresContainerEvents", tags, fields, time.Now())
		if err != nil {
			return fmt.Errorf("Couldn't add FlowEvent: %s", err)
		}
		bp.AddPoint(pt)
	case UNKNOWN:
		pt, err := client.NewPoint("UnknownContainerEvents", tags, fields, time.Now())
		if err != nil {
			return fmt.Errorf("Couldn't add FlowEvent: %s", err)
		}
		bp.AddPoint(pt)
	}

	if err := d.httpClient.Write(bp); err != nil {
		return fmt.Errorf("Couldn't add data: %s", err)
	}

	return nil
}

// CollectContainerEvent implements trireme collector interface
func (d *Influxdb) CollectContainerEvent(container *Container) {
	d.Lock()
	d.worker.addEvent(
		&workerEvent{
			containerName: container.ContainerName,
			container:     container,
		},
	)
	d.Unlock()
}
