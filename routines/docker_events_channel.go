package routines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/elastic/go-elasticsearch/v9"
)

func PushDockerEventsToElastic(dockerPullEvent *model.DockerPullEvent, m *model.Model, logger cego.Logger) {
	elasticsearchConfig := m.GetElasticsearchConfig()
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: elasticsearchConfig.Hosts,
		APIKey:    elasticsearchConfig.ApiKey,
		Username:  elasticsearchConfig.Username,
		Password:  elasticsearchConfig.Password,
	})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug(fmt.Sprintf("Detected docker pull event for %s pushing to %s for index %s", dockerPullEvent.ImageName, elasticsearchConfig.Hosts, elasticsearchConfig.Index))
	data, _ := json.Marshal(dockerPullEvent)
	_, err = es.Index(elasticsearchConfig.Index, bytes.NewReader(data))
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func StartDockerEventsChannel(ctx context.Context, m *model.Model, d *client.Client, logger cego.Logger) {
	go func() {
		messageChan, errChan := d.Events(ctx, events.ListOptions{
			Filters: filters.NewArgs(
				filters.Arg("event", string(events.ActionPull)),
				filters.Arg("event", string(events.ActionCreate)),
			),
		})

		logger.Debug("Listening for docker events...")

		for {
			select {
			case message := <-messageChan:
				switch message.Action {
				case events.ActionPull:
					imageName := message.Actor.ID
					dockerPullEvent := &model.DockerPullEvent{
						Timestamp: time.Now().Format(time.RFC3339),
						Message:   "dopetes detected docker pull event for " + imageName,
						ImageName: imageName,
					}
					PushDockerEventsToElastic(dockerPullEvent, m, logger)
				case events.ActionCreate:
					imageName := message.Actor.Attributes["image"]
					if !strings.Contains(imageName, ":") {
						imageName = imageName + ":latest"
					}
					dockerPullEvent := &model.DockerPullEvent{
						Timestamp: time.Now().Format(time.RFC3339),
						Message:   "dopetes detected docker pull event for " + imageName,
						ImageName: imageName,
					}
					PushDockerEventsToElastic(dockerPullEvent, m, logger)
				}
			case err := <-errChan:
				logger.Error(err.Error())
			default:
				ctx.Done()
			}
		}
	}()
}
