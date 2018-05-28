package env

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type ServerConfig struct {
	Name string
}

type DatabaseConfig struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func NewDbConfig(name string, host string, port int, user string, password string) (DatabaseConfig, error) {
	db := DatabaseConfig{
		Name:     name,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
	return db, nil
}

func NewServerConfig(name string) (ServerConfig, error) {
	s := ServerConfig{
		Name: name,
	}
	return s, nil
}

func CreateNewDbConfig(name string, host string, port int, user string, password string) DatabaseConfig {
	cfg, err := NewDbConfig(name, host, port, user, password)
	if err != nil {
		fmt.Printf("Fail to create Database Config", err)
	}
	return cfg
}

func CreateNewServerConfig(name string) ServerConfig {
	cfg, err := NewServerConfig(name)
	if err != nil {
		fmt.Printf("Fail to create Database Config", err)
	}
	return cfg
}

func GetConfig() *Config {
	conf, err := ini.Load("files/etc/config/hello.development.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}
	dbSection := conf.Section("Database")
	scfg := CreateNewServerConfig(conf.Section("Server").Key("Name").String())
	dbcfg := CreateNewDbConfig(
		dbSection.Key("dbname").String(),
		dbSection.Key("Host").String(),
		dbSection.Key("Port").MustInt(),
		dbSection.Key("User").String(),
		dbSection.Key("Password").String(),
	)
	return &Config{
		Server:   scfg,
		Database: dbcfg,
	}
}
