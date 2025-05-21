package config

import (
	"log"
	"strings"
	"user-service/internal/config/yml_config"

	"github.com/spf13/viper"
)

type AppConfig struct {
	server   ServerConfig
	database DatabaseConfig
}

func (config *AppConfig) GetServerConfig() ServerConfig {
	return config.server
}

func (config *AppConfig) GetDatabaseConfig() DatabaseConfig {
	return config.database
}

func newAppConfig(yml *yml_config.YMLAppConfig) *AppConfig {
	return &AppConfig{
		server:   newServerConfig(&yml.Server),
		database: newDatabaseConfig(&yml.Database),
	}
}

func LoadConfig(path string) (*AppConfig, error) {
	viperInstance := viper.New()
	viperInstance.SetConfigName(path)
	viperInstance.AddConfigPath(".")
	viperInstance.AutomaticEnv()
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viperInstance.ReadInConfig(); err != nil {
		log.Println("Ошибка чтения YML!", err)
		return nil, err
	}

	var ymlConfig yml_config.YMLAppConfig
	if err := viperInstance.Unmarshal(&ymlConfig); err != nil {
		log.Println("Ошибка конвертации YML в yml_config.YMLConfig!", err)
		return nil, err
	}

	return newAppConfig(&ymlConfig), nil
}
