package config

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "50053",
		},
	}
}
