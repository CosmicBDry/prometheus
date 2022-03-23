package models

type ConfigJobs struct {
	Global        interface{} `yaml:"global" json:"global"`
	Alerting      interface{} `yaml:"alerting" json:"alerting"`
	RuleFiles     interface{} `yaml:"rule_files" json:"rule_files"`
	ScrapeConfigs []*struct {
		JobName     string `yaml:"job_name" json:"job_name"`
		MetricsPath string `yaml:"metrics_path" json:"metrics_path"`
		Scheme      string `yaml:"scheme"  json:"scheme"`
		BasicAuth   struct {
			Username string ` yaml:"username" json:"username"`
			Password string ` yaml:"password" json:"password"`
		} ` yaml:"basic_auth" json:"basic_auth"`
		FileSdConfigs []*struct {
			Files           []string ` yaml:"files" json:"files"`
			RefreshInterval string   `yaml:"refresh_interval" json:"refresh_interval"`
		} `yaml:"file_sd_configs" json:"file_sd_configs"`
	} `yaml:"scrape_configs" json:"scrape_configs"`
}

type ResponseJson struct {
	Code          int    `json:"code"`
	Config        string `json:"config"`
	ConfigVersion string `json:"configversion"`
	Text          string `json:"text"`
	Error         string `json:"error"`
}
