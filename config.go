package main

import (
	"encoding/json"
	"io/ioutil"
)

type DBConfig struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Type     string `json:"type, omitempty"`
	DBName   string `json:"dbname, omitempty"`
}

type Config struct {
	DB        *DBConfig `json:"db,omitempty"`
	TableName string    `json:"tablename,omitempty"`
	Step      int       `json:"step,omitempty"`
}

func NewDBConfigDefault() *DBConfig {
	cfg := new(DBConfig)
	cfg.Host = "127.0.0.1"
	cfg.Port = 3306
	cfg.User = "root"
	cfg.Password = ""
	cfg.Type = "mysql"
	cfg.DBName = "service"
	return cfg
}

func NewConfigDefault() *Config {
	cfg := new(Config)
	cfg.Step = 1000
	cfg.TableName = "id_gen"

	cfg.DB = NewDBConfigDefault()
	return cfg
}

func NewConfigFromFile(filename string) *Config {
	data, _ := ioutil.ReadFile(filename)
	return NewConfig(data)
}

func NewConfig(data []byte) *Config {
	cfg := NewConfigDefault()
	json.Unmarshal(data, cfg)
	return cfg
}
