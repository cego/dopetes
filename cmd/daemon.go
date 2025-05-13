package cmd

import (
	"os"

	"github.com/cego/dopetes/model"
	"github.com/cego/dopetes/routines"
	"github.com/cego/go-lib"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func DaemonRun(cmd *cobra.Command, _ []string) {
	l := cego.NewLogger()
	ctx := cmd.Context()
	dockerEvents := make(chan *model.DockerPullEvent, 50)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	configEndpoint, err := cmd.Flags().GetString("config-endpoint")
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	go routines.StartDockerEventsChannel(ctx, dockerClient, l, dockerEvents)

	config, err := routines.FetchConfig(ctx, l, configEndpoint)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	err = routines.PushDockerEventsToElastic(ctx, l, config, dockerEvents)
	if err != nil {
		l.Error(err.Error())
	}
}

func InitDaemon() *cobra.Command {
	var daemon = &cobra.Command{
		Use:   "daemon",
		Short: "Starts daemon that stores docker pull events and hosts a rest api",
		Run:   DaemonRun,
	}

	daemon.Flags().String("config-endpoint", "http://localhost:8080/dopetes.yaml", "Endpoint to fetch config from")
	return daemon
}
