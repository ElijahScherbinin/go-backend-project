package config

import "user-service/internal/config/yml_config"

type ymlServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ServerConfig struct {
	host string
	port int
}

// GetHost возвращает хост сервера
func (config *ServerConfig) GetHost() string {
	return config.host
}

// GetPort возвращает порт сервера
func (config *ServerConfig) GetPort() int {
	return config.port
}

func newServerConfig(yml *yml_config.YMLServerConfig) ServerConfig {
	return ServerConfig{
		host: yml.Host,
		port: yml.Port,
	}
}
