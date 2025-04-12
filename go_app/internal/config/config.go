package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port    string
	Host    string
	Mode    string // "debug" or "release"
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "0.0.0.0",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "gitlab",
		},
	}
}
