package config

type ServerConfig struct {
	Port string `yaml:"port"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
	Timeout  int    `yaml:"timeout"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "50052",
		},
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "order_service",
			Timeout:  10,
		},
	}
}
