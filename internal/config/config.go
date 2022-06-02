package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Level   string
	Storage string
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HttpConf
	Rabbit  RabbitConf
}

type StorageConf struct {
	Type Storage `json:"type"`
	Dsn  string  `json:"dsn"`
}

type HttpConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RabbitConf struct {
	Url      string `json:"url"`
	Queue    string `json:"queue"`
	Exchange string `json:"exchange"`
}

type LoggerConf struct {
	Level            string        `json:"level"`
	Encoding         string        `json:"encoding"`
	OutputPaths      []string      `json:"outputPaths"`
	ErrorOutputPaths []string      `json:"errorOutputPaths"`
	EncoderConfig    encoderConfig `json:"encoderConfig"`
}

type encoderConfig struct {
	MessageKey   string `json:"messageKey"`
	LevelKey     string `json:"levelKey"`
	LevelEncoder string `json:"levelEncoder"`
}

const SQL Storage = "sql"

func NewConfig() Config {
	return Config{}
}

func LoadConfig(path string) (*Config, error) {
	resultConfig, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}

	config := NewConfig()
	err = json.Unmarshal(resultConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid unmarshal config %s:%w", path, err)
	}

	return &config, nil
}
