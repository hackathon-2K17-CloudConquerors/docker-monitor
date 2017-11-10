package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/docker-monitor/configuration"
	"github.com/docker-monitor/constants"
	"github.com/docker-monitor/influxdblib"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func banner(version, revision string) {
	fmt.Printf(`


	  _____     _
	 |_   _| __(_)_ __ ___ _ __ ___   ___
	   | || '__| | '__/ _ \ '_'' _ \ / _ \
	   | || |  | | | |  __/ | | | | |  __/
	   |_||_|  |_|_|  \___|_| |_| |_|\___|
		STATISTICS

_______________________________________________________________
             %s - %s
                                          🚀  by CloudConquerors

`, version, revision)
}

func main() {
	cfg, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatal("Error parsing configuration", err)
	}

	err = setLogs(cfg.LogFormat, "debug")
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

	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal("Error retrieving docker client", err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	go pollNginxContainer(influxInstance, dockerClient, cfg.PollInterval)
	go pollHttpdContainer(influxInstance, dockerClient, cfg.PollInterval)
	go pollPostgresContainer(influxInstance, dockerClient, cfg.PollInterval)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	zap.L().Info("Everything started. Waiting for Stop signal")
	// Waiting for a Sig
	<-c

}

func pollNginxContainer(influxInstance *influxdblib.Influxdb, dockerClientInstance *dockerClient.Client, interval int) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		containers, err := dockerClientInstance.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			log.Fatal("Error retrieving docker client", err)
		}

		for _, container := range containers {
			fmt.Printf("%s %s\n", container.ID[:10], container.Image)
		}
		//influxInstance.CollectContainerEvent(influxdblib.NGINX)
	}
}

func pollHttpdContainer(influxInstance *influxdblib.Influxdb, dockerClientInstance *dockerClient.Client, interval int) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		influxInstance.CollectContainerEvent(influxdblib.HTTPD)
	}
}

func pollPostgresContainer(influxInstance *influxdblib.Influxdb, dockerClientInstance *dockerClient.Client, interval int) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		influxInstance.CollectContainerEvent(influxdblib.POSTGRES)
	}
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
