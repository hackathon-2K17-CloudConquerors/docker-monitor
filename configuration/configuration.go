package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	flag "github.com/spf13/pflag"
)

// Configuration stuct is used to populate the various fields used by trireme-statistics
type Configuration struct {
	ListenAddress string

	InfluxUsername string
	InfluxPassword string
	InfluxDBName   string
	InfluxURL      string
	DBSkipTLS      bool

	PollInterval int

	LogFormat string
	LogLevel  string
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

// LoadConfiguration will load the configuration struct
func LoadConfiguration() (*Configuration, error) {
	flag.Usage = usage
	flag.String("ListenAddress", "", "Server Address [Default: 8080]")
	flag.String("LogLevel", "", "Log level. Default to info (trace//debug//info//warn//error//fatal)")
	flag.String("LogFormat", "", "Log Format. Default to human")

	flag.String("InfluxUsername", "", "Username of the database [default: sibi]")
	flag.String("InfluxPassword", "", "Password of the database [default: sibi]")
	flag.String("InfluxDBName", "", "Name of the database [default: flowDB]")
	flag.String("InfluxURL", "", "URI to connect to DB [default: http://influxdb:8086]")
	flag.Bool("DBSkipTLS", true, "Is valid TLS required for the DB server ? [default: true]")

	flag.Int("PollInterval", 20, "Time interval to poll containers [default: 5m]")

	// Setting up default configuration
	viper.SetDefault("ListenAddress", ":8080")
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFormat", "human")

	viper.SetDefault("InfluxUsername", "aporeto")
	viper.SetDefault("InfluxPassword", "aporeto")
	viper.SetDefault("InfluxDBName", "flowDB")
	viper.SetDefault("InfluxURL", "http://influxdb:8086")
	viper.SetDefault("DBSkipTLS", true)

	viper.SetDefault("PollInterval", 5)

	// Binding ENV variables
	// Each config will be of format TRIREME_XYZ as env variable, where XYZ
	// is the upper case config.
	viper.SetEnvPrefix("TRIREME")
	viper.AutomaticEnv()

	// Binding CLI flags.
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	var config Configuration

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %s", err)
	}

	return &config, nil
}
