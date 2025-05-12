package model

import "time"

type DockerPullEvent struct {
	Timestamp string `json:"@timestamp"`
	Message   string `json:"message"`
	ImageName string `json:"dopetes.image_name"`
}

type Model struct {
	dockerPullEvents []*DockerPullEvent
}

func New() *Model {
	return &Model{
		dockerPullEvents: []*DockerPullEvent{},
	}
}

func (m *Model) AddDockerPullEvent(imageName string) {
	m.dockerPullEvents = append(m.dockerPullEvents, &DockerPullEvent{
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   "dopetes detected docker pull event for " + imageName,
		ImageName: imageName,
	})
}

func (m *Model) ClearDockerPullEvents() {
	m.dockerPullEvents = []*DockerPullEvent{}
}

func (m *Model) GetDockerPullEvents() []*DockerPullEvent {
	return m.dockerPullEvents
}
