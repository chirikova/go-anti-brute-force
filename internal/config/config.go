package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  Logger
	DB      DB
	GRPC    GRPC
	Limiter struct {
		Login Limit
		Pass  Limit
		IP    Limit
	}
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

type Limit struct {
	Limit    int64
	Interval time.Duration
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
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.DBName,
	)

	return dsn
}
