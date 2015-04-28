package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type DBConfig struct {
	Host      string `json:"host,omitempty"`
	Port      int    `json:"port,omitempty"`
	User      string `json:"user,omitempty"`
	Password  string `json:"password,omitempty"`
	Type      string `json:"type, omitempty"`
	DBName    string `json:"dbname, omitempty"`
	TableName string `json:"tablename, omitempty"`
}

type Config struct {
	DB   *DBConfig `json:"db"`
	Step int       `json:"step"`
}

func NewDBConfigDefault() *DBConfig {
	cfg := new(DBConfig)
	cfg.Host = "127.0.0.1"
	cfg.Port = 3306
	cfg.User = "root"
	cfg.Password = ""
	cfg.Type = "mysql"
	cfg.DBName = "service"
	cfg.TableName = "id_gen"
	return cfg
}

func NewConfigDefault() *Config {
	cfg := new(Config)
	cfg.Step = 1000
	cfg.DB = NewDBConfigDefault()
	return cfg
}

func NewConfigFromFile(filename string) *Config {
	data, _ := ioutil.ReadFile(filename)
	return NewConfig(data)
}

func NewConfig(data []byte) *Config {
	cfg := NewConfigDefault()
	err := json.Unmarshal(data, cfg)
	if err != nil {
		log.Fatal("config is not a valid json")
	}
	return cfg
}
