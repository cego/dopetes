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

func PushDockerEventsToElastic(ctx context.Context, logger cego.Logger, config *model.DopetesConfig, elasticClient *elasticsearch.Client, elasticDocumentChan chan *model.ElasticDocument) {
	for {
		select {
		case e := <-elasticDocumentChan:
			logger.Debug(fmt.Sprintf("Detected docker pull event for %s pushing to %s for index %s", e.ImageName, config.Elasticsearch.Hosts, config.Elasticsearch.Index))

			data, _ := json.Marshal(e)
			_, err := elasticClient.Index(config.Elasticsearch.Index, bytes.NewReader(data))
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to sent event to elasticsearch: %v", err))
				continue
			}
		case <-ctx.Done():
		}
	}
}

func StartDockerEventsChannel(ctx context.Context, dockerClient *client.Client, logger cego.Logger, elasticDocumentChan chan *model.ElasticDocument) {
	messageChan, errChan := dockerClient.Events(ctx, events.ListOptions{
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
				res, _ := json.Marshal(message)
				imageName := message.Actor.ID
				dockerPullEvent := &model.ElasticDocument{
					Timestamp: time.Now().Format(time.RFC3339),
					Message:   "dopetes detected docker pull event for " + imageName,
					ImageName: imageName,
					Type:      "event",
					EventRaw:  string(res),
				}
				elasticDocumentChan <- dockerPullEvent
			case events.ActionCreate:
				if message.Type != events.ContainerEventType {
					continue
				}
				imageName := message.Actor.Attributes["image"]
				if strings.HasPrefix(imageName, "sha256:") {
					continue
				}
				if !strings.Contains(imageName, ":") {
					imageName = imageName + ":latest"
				}
				res, _ := json.Marshal(message)
				dockerPullEvent := &model.ElasticDocument{
					Timestamp: time.Now().Format(time.RFC3339),
					Message:   fmt.Sprintf("dopetes detected docker create event of type %s for %s", message.Type, imageName),
					ImageName: imageName,
					Type:      "event",
					EventRaw:  string(res),
				}
				elasticDocumentChan <- dockerPullEvent
			}
		case err := <-errChan:
			logger.Error(err.Error())
			return
		}
	}
}
