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

type RedisConfig struct {
	URI      string `yaml:"uri"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	TTL      int    `yaml:"ttl"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	NATS    NATSConfig    `yaml:"nats"`
	Redis   RedisConfig   `yaml:"redis"`
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "50051",
		},
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "inventory_service",
			Timeout:  10,
		},
		NATS: NATSConfig{
			URL:     "nats://localhost:4222",
			Cluster: "microservices",
		},
		Redis: RedisConfig{
			URI:      "localhost:6379",
			Password: "",
			DB:       0,
			TTL:      60, // Cache TTL in minutes
		},
	}
}
