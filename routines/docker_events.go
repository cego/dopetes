package routines

import (
	"context"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func StartDockerEventsRoutine(ctx context.Context, m *model.Model, d *client.Client, logger cego.Logger) {
	go func() {
		messageChan, errChan := d.Events(ctx, events.ListOptions{
			Filters: filters.NewArgs(
				filters.Arg("event", string(events.ActionPull)),
				filters.Arg("event", string(events.ActionCreate)),
			),
		})

		for {
			select {
			case message := <-messageChan:
				switch message.Action {
				case events.ActionPull:
					m.AddDockerPullEvent(message.Actor.ID)
				case events.ActionCreate:
					m.AddDockerPullEvent(message.Actor.Attributes["image"])
				}
			case err := <-errChan:
				logger.Error(err.Error())
			default:
				ctx.Done()
			}
		}
	}()
}
