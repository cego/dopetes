package model

import (
	"slices"
)

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

type DockerBuildxHistoryState struct {
	ids []string
}

func (r *DockerBuildxHistoryState) AddId(id string) {
	r.ids = append(r.ids, id)
}

func (r *DockerBuildxHistoryState) HasId(id string) bool {
	idx := slices.IndexFunc(r.ids, func(c string) bool { return c == id })
	return idx != -1
}

func NewDockerBuildxHistoryState() *DockerBuildxHistoryState {
	return &DockerBuildxHistoryState{
		ids: []string{},
	}
}

type ElasticDocument struct {
	Timestamp string `json:"@timestamp"`
	Message   string `json:"message"`
	ImageName string `json:"dopetes.image_name"`
	Type      string `json:"dopetes.type"`
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
