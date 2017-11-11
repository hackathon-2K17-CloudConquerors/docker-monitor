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

func banner(version, revision string) {
	fmt.Printf(`

         _                                                                                 _    _
        | |   ____    ____   _   _    ___    __ ___           _ _______    ____   _____   (_)  | |__    ____    __ ___
     __ | |  / __ \  /  __| / | / /  / __ \ |  '___|         | /__  __ \  / __ \  | __ \  | |  | ___|  / __ \  |  '___|
    / _|| | / /  \ | | /    | |/ /  |  __ / | |       ___    | | | | | | / /  \ | | | | | | |  | |    / /  \ | | |
   | |_ | | | \__/ | | \__  |  _ \  | \___  | |      |___|   | | | | | | | \__/ | | | | | | |  | |__  | \__/ | | |
   \_____ /  \____/  \____| \_| \_\  \____| |_|              |_| |_| |_|  \____/  |_| |_| |_|  \____|  \____/  |_|

_______________________________________________________________
             %s - %s
                                          🚀  by CloudConquerors

`, version, revision)
}

func main() {

	banner("hi", "hi")
	cfg, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatal("Error parsing configuration", err)
	}

	err = setLogs(cfg.LogFormat, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Error setting up logs: %s", err)
	}

	zap.L().Debug("Config used", zap.Any("Config", cfg))

	influxInstance, err := influxdblib.NewDBConnection(cfg.InfluxUsername, cfg.InfluxPassword, cfg.InfluxURL, cfg.InfluxDBName, cfg.DBSkipTLS)
	if err != nil {
		log.Fatal("Error initiating connection to DB", err)
	}
	// Start worker
	influxInstance.Start()

	dockerClient, err := initDockerClient(constants.DefaultDockerSocketType, constants.DefaultDockerSocket)
	if err != nil {
		log.Fatal("Error initializing docker client", err)
	}

	monitorInstance := monitor.NewMonitor(cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPServer, influxInstance, cfg.InfluxDBName, cfg.SMTPUser, cfg.RecipientAddress)

	go pollNginxContainer(influxInstance, dockerClient, cfg.PollInterval, monitorInstance)
	go pollHttpdContainer(influxInstance, dockerClient, cfg.PollInterval, monitorInstance)
	go pollPostgresContainer(influxInstance, dockerClient, cfg.PollInterval, monitorInstance)

	go func() {
		err := startMonitorServer(cfg.ListenAddress, monitorInstance)
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
	for range time.Tick(time.Second * time.Duration(10)) {
		var isNginxRunning bool
		containers, err := dockerClientInstance.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			log.Fatal("Error retrieving docker client", err)
		}

		for _, container := range containers {
			if strings.Contains(container.Image, influxdblib.NGINX) {
				isNginxRunning = true
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerName: influxdblib.NGINX,
					ContainerID:   container.ID[:10],
					ImageName:     container.Image,
					Status:        influxdblib.ContainerStart,
				})
			} else {
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerID: container.ID[:10],
					ImageName:   container.Image,
					Status:      influxdblib.ContainerStart,
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
					ContainerName: influxdblib.HTTPD,
					ContainerID:   container.ID[:10],
					ImageName:     container.Image,
					Status:        influxdblib.ContainerStart,
				})
			} else {
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerID: container.ID[:10],
					ImageName:   container.Image,
					Status:      influxdblib.ContainerStart,
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
					ContainerName: influxdblib.POSTGRES,
					ContainerID:   container.ID[:10],
					ImageName:     container.Image,
					Status:        influxdblib.ContainerStart,
				})
			} else {
				influxInstance.CollectContainerEvent(&influxdblib.Container{
					ContainerID: container.ID[:10],
					ImageName:   container.Image,
					Status:      influxdblib.ContainerStart,
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

func startMonitorServer(listenAddress string, monitor monitor.MonitorManipulator) error {
	mux := http.NewServeMux()

	go monitor.StartMonitor(10)

	//mux.HandleFunc("/", monitor.StartContainer)

	handler := cors.Default().Handler(mux)

	err := http.ListenAndServe(listenAddress, handler)
	if err != nil {
		return fmt.Errorf("ListenAndServe: %s", err)
	}

	zap.L().Info("Server Listening at", zap.Any("port", listenAddress))
	return nil
}

func initDockerClient(socketType string, socketAddress string) (*dockerClient.Client, error) {

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
