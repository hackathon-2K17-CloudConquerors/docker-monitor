package main

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/docker-monitor/configuration"
	"github.com/docker-monitor/influxdblib"
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

	influxInstance.Start()

	go pollNginxContainer(influxInstance, cfg.PollInterval)
	go pollHttpdContainer(influxInstance, cfg.PollInterval)
	go pollPostgresContainer(influxInstance, cfg.PollInterval)

	for {

	}

}

func pollNginxContainer(influxInstance *influxdblib.Influxdb, interval int) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		influxInstance.CollectContainerEvent(influxdblib.NGINX)
	}
}

func pollHttpdContainer(influxInstance *influxdblib.Influxdb, interval int) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		influxInstance.CollectContainerEvent(influxdblib.HTTPD)
	}
}

func pollPostgresContainer(influxInstance *influxdblib.Influxdb, interval int) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		influxInstance.CollectContainerEvent(influxdblib.POSTGRES)
	}
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
