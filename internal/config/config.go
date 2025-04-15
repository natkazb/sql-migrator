package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	Sql    SQLConf    `yaml:"sql"`
	Path string `yaml:"path"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}


type SQLConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbName"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Driver   string `yaml:"driver"`
}

func NewConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error in opening file %s: %w", filePath, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("error in decoding %s: %w", filePath, err)
	}
	return config, nil
}
