package model

type Model struct {
	dockerPullEvents []string
}

func New() *Model {
	return &Model{
		dockerPullEvents: []string{},
	}
}

func (m *Model) AddDockerPullEvent(imageName string) {
	m.dockerPullEvents = append(m.dockerPullEvents, imageName)
}

func (m *Model) GetDockerPullEvents() []string {
	return m.dockerPullEvents
}
