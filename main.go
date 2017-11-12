package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/docker-monitor/configuration"
	"github.com/docker-monitor/constants"
	"github.com/docker-monitor/influxdblib"
	"github.com/docker-monitor/monitor"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func banner() {
	fmt.Printf(`

         _                                                                                  _    _
        | |   ____    _____   _   _    ___    __ ___           _ _______    ____   _____   (_)  | |__    ____    __ ___
     __ | |  / __ \  /  ___| / | / /  / __ \ |  '___|         | /__  __ \  / __ \  | __ \  | |  | ___|  / __ \  |  '___|
    / _|| | / /  \ | | |     | |/ /  |  __ / | |       ___    | | | | | | / /  \ | | | | | | |  | |    / /  \ | | |
   | |_ | | | \__/ | | |___  |  _ \  | \___  | |      |___|   | | | | | | | \__/ | | | | | | |  | |__  | \__/ | | |
    \____ /  \____/  \_____| \_| \_\  \____| |_|              |_| |_| |_|  \____/  |_| |_| |_|  \____|  \____/  |_|

_______________________________________________________________

                                          🚀  by CloudConquerors

`)
}

func main() {

	banner()
	// LoadConfiguration loads the config from flags or env variables
	cfg, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatal("Error parsing configuration", err)
	}

	err = setLogs(cfg.LogFormat, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Error setting up logs: %s", err)
	}

	zap.L().Debug("Config used", zap.Any("Config", cfg))

	// Create a new DB instance
	influxInstance, err := influxdblib.NewDBConnection(
		cfg.InfluxUsername,
		cfg.InfluxPassword,
		cfg.InfluxURL,
		cfg.InfluxDBName,
		cfg.DBSkipTLS,
	)
	if err != nil {
		log.Fatal("Error initiating connection to DB", err)
	}
	// Start influx worker
	influxInstance.Start()

	// Create docker client
	dockerClient, err := initDockerClient(constants.DefaultDockerSocketType, constants.DefaultDockerSocket)
	if err != nil {
		log.Fatal("Error initializing docker client", err)
	}

	// Create new monitor instance
	monitorInstance := monitor.NewMonitor(
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.SMTPServer,
		dockerClient,
		influxInstance,
		cfg.InfluxDBName,
		cfg.RecipientAddress,
	)

	// Start looking for containers and add to DB
	go pollNginxContainer(influxInstance, dockerClient, cfg.PollInterval, monitorInstance)
	go pollHttpdContainer(influxInstance, dockerClient, cfg.PollInterval, monitorInstance)
	go pollPostgresContainer(influxInstance, dockerClient, cfg.PollInterval, monitorInstance)

	// Start server
	go func() {
		err := startMonitorServer(cfg.ListenAddress, monitorInstance, cfg.MonitorInterval)
		if err != nil {
			zap.L().Fatal("Error: Connecting to GraphServer", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	zap.L().Info("Everything started. Waiting for Stop signal")
	// Waiting for a Sig
	<-c

}

func pollNginxContainer(influxInstance *influxdblib.Influxdb, dockerClientInstance *dockerClient.Client, interval int, email monitor.MonitorManipulator) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		var isNginxRunning bool
		containers, err := dockerClientInstance.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			log.Fatal("Error retrieving docker client", err)
		}

		for _, container := range containers {
			if strings.Contains(container.Image, influxdblib.NGINX) {
				isNginxRunning = true
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerName:    influxdblib.NGINX,
					ContainerID:      container.ID[:10],
					ImageName:        container.Image,
					ContainerState:   container.State,
					ContainerStatus:  container.Status,
					ContainerCreated: container.Created,
					ContainerNetwork: container.HostConfig.NetworkMode,
					Status:           influxdblib.ContainerStart,
				})
			} else {
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerID:      container.ID[:10],
					ImageName:        container.Image,
					ContainerState:   container.State,
					ContainerStatus:  container.Status,
					ContainerCreated: container.Created,
					ContainerNetwork: container.HostConfig.NetworkMode,
					Status:           influxdblib.ContainerStart,
				})
			}
		}

		if !isNginxRunning {
			influxInstance.CollectContainerEvent(&influxdblib.Container{
				ContainerName: influxdblib.NGINX,
				Status:        influxdblib.ContainerStop,
			})
		}
	}
	return
}

func pollHttpdContainer(influxInstance *influxdblib.Influxdb, dockerClientInstance *dockerClient.Client, interval int, email monitor.MonitorManipulator) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		var isHttpdRunning bool
		containers, err := dockerClientInstance.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			log.Fatal("Error retrieving docker client", err)
		}

		for _, container := range containers {
			if strings.Contains(container.Image, influxdblib.HTTPD) {
				isHttpdRunning = true
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerName:    influxdblib.HTTPD,
					ContainerID:      container.ID[:10],
					ImageName:        container.Image,
					ContainerState:   container.State,
					ContainerStatus:  container.Status,
					ContainerCreated: container.Created,
					ContainerNetwork: container.HostConfig.NetworkMode,
					Status:           influxdblib.ContainerStart,
				})
			} else {
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerID:      container.ID[:10],
					ImageName:        container.Image,
					ContainerState:   container.State,
					ContainerStatus:  container.Status,
					ContainerCreated: container.Created,
					ContainerNetwork: container.HostConfig.NetworkMode,
					Status:           influxdblib.ContainerStart,
				})
			}
		}

		if !isHttpdRunning {
			influxInstance.CollectContainerEvent(&influxdblib.Container{
				ContainerName: influxdblib.HTTPD,
				Status:        influxdblib.ContainerStop,
			})
		}
	}
	return
}

func pollPostgresContainer(influxInstance *influxdblib.Influxdb, dockerClientInstance *dockerClient.Client, interval int, email monitor.MonitorManipulator) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		var isPostgresRunning bool
		containers, err := dockerClientInstance.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			log.Fatal("Error retrieving docker client", err)
		}

		for _, container := range containers {
			if strings.Contains(container.Image, influxdblib.POSTGRES) {
				isPostgresRunning = true
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerName:    influxdblib.POSTGRES,
					ContainerID:      container.ID[:10],
					ImageName:        container.Image,
					ContainerState:   container.State,
					ContainerStatus:  container.Status,
					ContainerCreated: container.Created,
					ContainerNetwork: container.HostConfig.NetworkMode,
					Status:           influxdblib.ContainerStart,
				})
			} else {
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerID:      container.ID[:10],
					ImageName:        container.Image,
					ContainerState:   container.State,
					ContainerStatus:  container.Status,
					ContainerCreated: container.Created,
					ContainerNetwork: container.HostConfig.NetworkMode,
					Status:           influxdblib.ContainerStart,
				})
			}
		}

		if !isPostgresRunning {
			influxInstance.CollectContainerEvent(&influxdblib.Container{
				ContainerName: influxdblib.POSTGRES,
				Status:        influxdblib.ContainerStop,
			})
		}
	}
	return
}

func startMonitorServer(listenAddress string, monitor monitor.MonitorManipulator, interval int) error {
	mux := http.NewServeMux()

	// start processing emails
	go monitor.StartMonitor(interval)

	mux.HandleFunc("/monitor", monitor.StartContainer)

	handler := cors.Default().Handler(mux)

	err := http.ListenAndServe(listenAddress, handler)
	if err != nil {
		return fmt.Errorf("ListenAndServe: %s", err)
	}

	zap.L().Info("Server Listening at", zap.Any("port", listenAddress))
	return nil
}

func initDockerClient(socketType string, socketAddress string) (*dockerClient.Client, error) {
	zap.L().Info("Initializing Docker Client", zap.Any("socket", socketType))
	var socket string

	switch socketType {
	case "tcp":
		socket = "https://" + socketAddress

	case "unix":
		// Sanity check that this path exists
		if _, oserr := os.Stat(socketAddress); os.IsNotExist(oserr) {
			return nil, oserr
		}
		socket = "unix://" + socketAddress

	default:
		return nil, fmt.Errorf("Bad socket type %s", socketType)
	}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-dockerClient-1.0"}
	dockerClient, err := dockerClient.NewClient(socket, constants.DockerClientVersion, nil, defaultHeaders)

	if err != nil {
		return nil, fmt.Errorf("Error creating Docker Client %s", err.Error())
	}

	return dockerClient, nil
}

// setLogs setups Zap to log at the specified log level and format
func setLogs(logFormat, logLevel string) error {
	var zapConfig zap.Config

	switch logFormat {
	case "json":
		zapConfig = zap.NewProductionConfig()
		zapConfig.DisableStacktrace = true
	default:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.DisableStacktrace = true
		zapConfig.DisableCaller = true
		zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set the logger
	switch logLevel {
	case "trace":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}
