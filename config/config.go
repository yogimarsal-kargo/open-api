package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var cfgPath string

type Config struct {
	Service   ServiceConfig  `yaml:"service"`
	Database  DatabaseConfig `yaml:"database"`
	Consumers struct {
		ConsumerA ConsumerConfig `yaml:"consumer_a"`
		ConsumerB ConsumerConfig `yaml:"consumer_b"`
	} `yaml:"consumers"`
	NewRelic    NewRelicConfig    `yaml:"newrelic"`
	Log         LogConfig         `yaml:"log"`
	FeatureFlag FeatureFlagConfig `yaml:"feature_flag"`
}

type ServiceConfig struct {
	Name string `yaml:"name"`
}
type ConsumerConfig struct {
	Topic       string `yaml:"topic"`
	Concurrency int    `yaml:"concurrency"`
	MaxInFlight int    `yaml:"max-in-flight"`
}

type DatabaseConfig struct {
	Name     string `env:"DB_NAME" env-description:"Database name"`
	Port     string `env:"DB_PORT" env-description:"Database port"`
	Host     string `env:"DB_HOST" env-description:"Database host"`
	Username string `env:"DB_USERNAME" env-description:"Username of DB"`
	Password string `env:"DB_PASSWORD" env-description:"Password for the corresponding Username DB"`
}

type FeatureFlagConfig struct {
	AppName string `env:"FEATURE_FLAG_APP_NAME" env-description:"Feature Flag App Name"`
	URL     string `env:"FEATURE_FLAG_URL" env-description:"Feature Flag Url"`
	Token   string `env:"FEATURE_FLAG_TOKEN" env-description:"Feature Flag Token"`
}

type NewRelicConfig struct {
	AppName                  string `env:"NEW_RELIC_APPNAME" env-description:"New Relic application name"`
	LicenseKey               string `env:"NEW_RELIC_LICENSE" env-description:"New Relic license key"`
	DistributedTracerEnabled bool   `yaml:"distributed-tracer-enabled"`
}

type LogConfig struct {
	LogLevel string `yaml:"level"`
}

func InitializeConfig() Config {
	initFlag()

	var cfg Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return cfg
}

func initFlag() {
	flag.StringVar(&cfgPath, "config", "files/etc/app_config/config.dev.yaml", "location of config file")
	flag.Parse()
}

func PrintEnvDesc(cfg Config) {

	strConfig, err := cleanenv.GetDescription(&cfg, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(strConfig)
}
