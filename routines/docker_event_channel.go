package routines

import (
	"context"
	"fmt"

	"github.com/cego/dopetes/model"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func StartDockerEventListener(ctx context.Context, m *model.Model, d *client.Client) {
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
				if message.Action == events.ActionPull {
					m.AddDockerPullEvent(message.Actor.ID)
				} else if message.Action == events.ActionCreate {
					m.AddDockerPullEvent(message.Actor.Attributes["image"])
				}
			case err := <-errChan:
				fmt.Println(err)
			default:
				ctx.Done()
			}
		}
	}()
}
