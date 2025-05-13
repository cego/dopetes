package cmd

import (
	"net/http"
	"os"

	"github.com/cego/dopetes/model"
	"github.com/cego/dopetes/routines"
	"github.com/cego/go-lib"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func DaemonRun(cmd *cobra.Command, _ []string) {
	logger := cego.NewLogger()
	m := model.New()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	configEndpoint, err := cmd.Flags().GetString("config-endpoint")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	ctx := cmd.Context()
	httpClient := &http.Client{}

	routines.StartDockerEventsChannel(ctx, m, dockerClient, logger)
	routines.FetchConfig(ctx, m, logger, httpClient, configEndpoint)

	<-cmd.Context().Done()
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
