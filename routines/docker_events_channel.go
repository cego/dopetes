package routines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func PushDockerEventsToElastic(ctx context.Context, l cego.Logger, config *model.DopetesConfig, events chan *model.DockerPullEvent) error {
	if config == nil || config.Elasticsearch == nil {
		return fmt.Errorf("missing elasticsearch config")
	}
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: config.Elasticsearch.Hosts,
		APIKey:    config.Elasticsearch.ApiKey,
		Username:  config.Elasticsearch.Username,
		Password:  config.Elasticsearch.Password,
	})
	if err != nil {
		return fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	for {
		select {
		case e := <-events:
			l.Debug(fmt.Sprintf("Detected docker pull event for %s pushing to %s for index %s", e.ImageName, config.Elasticsearch.Hosts, config.Elasticsearch.Index))

			data, _ := json.Marshal(e)
			_, err = es.Index(config.Elasticsearch.Index, bytes.NewReader(data))
			if err != nil {
				l.Error(fmt.Sprintf("Failed to sent event to elasticsearch: %v", err))
				continue
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func StartDockerEventsChannel(ctx context.Context, d *client.Client, logger cego.Logger, dockerEvents chan *model.DockerPullEvent) {
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
				dockerEvents <- dockerPullEvent
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
				dockerEvents <- dockerPullEvent
			}
		case err := <-errChan:
			logger.Error(err.Error())
			return
		}
	}
}
