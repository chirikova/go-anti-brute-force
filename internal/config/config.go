package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger Logger
	DB     DB
	GRPC   GRPC
}

type Logger struct {
	Level      string
	OutputPath string `yaml:"outputPath"`
}

type DB struct {
	Host     string
	Port     string
	DBName   string `yaml:"dbName"`
	User     string
	Password string
}

type GRPC struct {
	Host string
	Port string
}

func InitConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("error closing config file %s", err)
		}
	}(file)

	return New(file)
}

func New(file *os.File) (*Config, error) {
	config := &Config{}

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (db *DB) BuildDSN() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=disable",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.DBName,
	)

	return dsn
}
