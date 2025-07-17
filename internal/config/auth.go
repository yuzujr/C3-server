package config

type AuthConfig struct {
	Enabled            bool
	Username           string
	Password           string
	SessionSecret      string
	SessionExpireHours string
}
