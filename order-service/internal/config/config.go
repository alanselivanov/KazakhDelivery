package config

type ServerConfig struct {
	Port string `yaml:"port"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
	Timeout  int    `yaml:"timeout"`
}

type NATSConfig struct {
	URL     string `yaml:"url"`
	Cluster string `yaml:"cluster"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	NATS    NATSConfig    `yaml:"nats"`
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
		NATS: NATSConfig{
			URL:     "nats://localhost:4222",
			Cluster: "microservices",
		},
	}
}
