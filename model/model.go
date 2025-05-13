package model

type DockerPullEvent struct {
	Timestamp string `json:"@timestamp"`
	Message   string `json:"message"`
	ImageName string `json:"dopetes.image_name"`
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

type Model struct {
	elasticSearchConfig ElasticSearchConfig
}

func New() *DopetesConfig {
	return &DopetesConfig{
		Elasticsearch: &ElasticSearchConfig{
			Hosts:    []string{"http://localhost:9200"},
			Username: "",
			Password: "",
			ApiKey:   "",
			Index:    "filebeat-prod-8.18.1",
		},
	}
}

func (m *Model) SetElasticsearchConfig(elasticsearchConfig ElasticSearchConfig) {
	m.elasticSearchConfig = elasticsearchConfig
}

func (m *Model) GetElasticsearchConfig() ElasticSearchConfig {
	return m.elasticSearchConfig
}
