package config

type AuthConfig struct {
	Enabled            bool
	Username           string
	Password           string
	SessionExpireHours int64
}
