package yml_config

type YMLServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
