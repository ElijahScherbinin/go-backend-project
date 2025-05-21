package config

import "user-service/internal/config/yml_config"

type DatabaseConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	sslmode  string
}

// GetHost возвращает хост базы данных
func (config *DatabaseConfig) GetHost() string {
	return config.host
}

// GetPort возвращает порт базы данных
func (config *DatabaseConfig) GetPort() int {
	return config.port
}

// GetUser возвращает имя пользователя
func (config *DatabaseConfig) GetUser() string {
	return config.user
}

// GetPassword возвращает пароль
func (config *DatabaseConfig) GetPassword() string {
	return config.password
}

// GetDBName возвращает имя базы данных
func (config *DatabaseConfig) GetDBName() string {
	return config.dbname
}

// GetSSLMode возвращает режим SSL
func (config *DatabaseConfig) GetSSLMode() string {
	return config.sslmode
}

func newDatabaseConfig(yml *yml_config.YMLDatabaseConfig) DatabaseConfig {
	return DatabaseConfig{
		host:     yml.Host,
		port:     yml.Port,
		user:     yml.User,
		password: yml.Password,
		dbname:   yml.DBName,
		sslmode:  yml.SSLMode,
	}
}
