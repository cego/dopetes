package model

type DockerBuildxHistoryLs struct {
	CreatedAt string `json:"created_at"`
	Ref       string `json:"ref"`
}

type DockerBuildxHistoryInspect struct {
	Materials []struct {
		URI     string   `json:"URI"`
		Digests []string `json:"Digests"`
	} `json:"Materials"`
}

type ElasticDocument struct {
	Timestamp string `json:"@timestamp"`
	Message   string `json:"message"`
	ImageName string `json:"dopetes.image_name"`
	EventRaw  string `json:"dopetes.event.raw"`
}

type DopetesConfig struct {
	Elasticsearch *ElasticSearchConfig `yaml:"elasticsearch"`
}

type ElasticSearchConfig struct {
	Hosts    []string `yaml:"hosts"`
	Username string   `yaml:"username,omitempty"`
	Password string   `yaml:"password,omitempty"`
	ApiKey   string   `yaml:"api_key,omitempty"`
	Index    string   `yaml:"index,omitempty"`
}
