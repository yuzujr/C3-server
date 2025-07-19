package models

type ClientConfig struct {
	ID         uint            `gorm:"primaryKey" json:"-"`
	ClientID   string          `gorm:"not null"`
	LastUpload string          `json:"lastUpload"`
	Api        ApiConfig       `gorm:"embedded;embeddedPrefix:api_" json:"api"`
	BuildInfo  BuildInfoConfig `gorm:"embedded;embeddedPrefix:build_" json:"build_info"`
}

type ApiConfig struct {
	Hostname               string `json:"hostname"`
	Port                   uint   `json:"port"`
	BasePath               string `json:"base_path"`
	UseSSL                 bool   `json:"use_ssl"`
	SkipSSLVerification    bool   `json:"skip_ssl_verification"`
	IntervalSeconds        uint   `json:"interval_seconds"`
	MaxRetries             uint   `json:"max_retries"`
	RetryDelayMilliseconds uint   `json:"retry_delay_ms"`
	AddToStartup           bool   `json:"add_to_startup"`
	ClientID               string `json:"client_id"`
}

type BuildInfoConfig struct {
	BuildTimestamp string `json:"build_timestamp"`
	Hardcoded      bool   `json:"hardcoded"`
	Preset         string `json:"preset"`
	PresetDesc     string `json:"preset_desc"`
	PresetName     string `json:"preset_name"`
}
