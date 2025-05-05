package config

type ServerConfig struct {
	Port string `yaml:"port"`
}

type ServicesConfig struct {
	Inventory string `yaml:"inventory"`
	Order     string `yaml:"order"`
	User      string `yaml:"user"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Services ServicesConfig `yaml:"services"`
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		Services: ServicesConfig{
			Inventory: "localhost:50051",
			Order:     "localhost:50052",
			User:      "localhost:50053",
		},
	}
}
