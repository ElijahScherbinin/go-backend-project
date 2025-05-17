package yml_config

type YMLAppConfig struct {
	Server   YMLServerConfig   `yaml:"server"`
	Database YMLDatabaseConfig `yaml:"database"`
}
