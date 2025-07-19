package config

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Logging  bool
}
