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

	PollInterval    int
	MonitorInterval int

	SMTPUser         string
	SMTPPassword     string
	SMTPServer       string
	RecipientAddress string

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
	flag.String("ListenAddress", "8088", "Server Address [Default: 8088]")
	flag.String("LogLevel", "", "Log level. Default to info (trace//debug//info//warn//error//fatal)")
	flag.String("LogFormat", "", "Log Format. Default to human")

	flag.String("InfluxUsername", "", "Username of the database [default: sibi]")
	flag.String("InfluxPassword", "", "Password of the database [default: sibi]")
	flag.String("InfluxDBName", "", "Name of the database [default: flowDB]")
	flag.String("InfluxURL", "", "URI to connect to DB [default: http://influxdb:8086]")
	flag.Bool("DBSkipTLS", true, "Is valid TLS required for the DB server ? [default: true]")

	flag.Int("PollInterval", 30, "Time interval to poll containers [default: 5m]")
	flag.Int("MonitorInterval", 10, "Time interval to watch for stopped containers [default: 5m]")

	flag.String("SMTPUser", "", "Username to connect to SMTP server [default: sibi]")
	flag.String("SMTPPassword", "", "Password to connect to SMTP server [default: sibi]")
	flag.String("SMTPServer", "", "SMTP server to send email")
	flag.String("RecipientAddress", "", "email of the receiver")

	// Setting up default configuration
	viper.SetDefault("ListenAddress", ":8088")
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFormat", "human")

	viper.SetDefault("InfluxUsername", "sibi")
	viper.SetDefault("InfluxPassword", "sibi")
	viper.SetDefault("InfluxDBName", "containers")
	viper.SetDefault("InfluxURL", "http://influxdb:8086")
	viper.SetDefault("DBSkipTLS", true)

	viper.SetDefault("PollInterval", 30)
	viper.SetDefault("MonitorInterval", 10)

	viper.SetDefault("SMTPUser", "cloudconquerors@gmail.com")
	viper.SetDefault("SMTPPassword", "sibicramesh")
	viper.SetDefault("SMTPServer", "smtp.gmail.com:587")
	viper.SetDefault("RecipientAddress", "cloudconquerors@gmail.com")

	// Binding ENV variables
	// Each config will be of format MONITOR_XYZ as env variable, where XYZ
	// is the upper case config.
	viper.SetEnvPrefix("MONITOR")
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
