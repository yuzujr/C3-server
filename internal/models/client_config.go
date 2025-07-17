package models

/*
	{
	  "api": {
	    "hostname": "127.0.0.1",
	    "port": 3000,
	    "base_path": "",
	    "use_ssl": false,
	    "skip_ssl_verification": false,
	    "interval_seconds": 60,
	    "max_retries": 3,
	    "retry_delay_ms": 1000,
	    "add_to_startup": false,
	    "client_id": ""
	  }
  	"lastUpload": "",
  	"clientId": "",
  	"build_info": {
    	"build_timestamp": "",
    	"hardcoded": true,
    	"preset": "development",
    	"preset_desc": "本地开发测试配置",
    	"preset_name": "开发环境"
  	}
	}
*/

type ClientConfig struct {
	Api        ApiConfig       `json:"api"`
	ClientID   string          `json:"clientId"`
	LastUpload string          `json:"lastUpload"`
	BuildInfo  BuildInfoConfig `json:"build_info"`
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
