package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Application struct {
		GRPCPort string `yaml:"grpc-port"`
		HTTPPort string `yaml:"http-port"`
		Version  string `yaml:"version"`
	} `yaml:"application"`
}

var AppConfig *Config

func NewConfig(path string) (*Config, error) {
	configPath, err := ParseFlags(path)
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	bindEnvironmentVariables(config)
	return config, nil
}

func bindEnvironmentVariables(config *Config) {
	viper.AutomaticEnv()
	//port := viper.GetString("port")
	//if strings.TrimSpace(port) != "" {
	//	config.Application.Port = port
	//}
	//dsn := viper.GetString("dsn")
	//fmt.Println("dsn--------------------------------")
	//fmt.Println(dsn)
	//if strings.TrimSpace(dsn) != "" {
	//	config.Application.DataSource = dsn
	//}
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags(path string) (string, error) {
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", path, "path to config file")
	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}
