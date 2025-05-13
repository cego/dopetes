package model

type DockerPullEvent struct {
	Timestamp string `json:"@timestamp"`
	Message   string `json:"message"`
	ImageName string `json:"dopetes.image_name"`
}

type DopetesConfig struct {
	Elasticsearch ElasticSearchConfig `json:"elasticsearch"`
}

type ElasticSearchConfig struct {
	Hosts    []string `json:"hosts"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	ApiKey   string   `json:"api_key,omitempty"`
	Index    string   `json:"index,omitempty"`
}

type Model struct {
	elasticSearchConfig ElasticSearchConfig
}

func New() *Model {
	return &Model{
		elasticSearchConfig: ElasticSearchConfig{
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
