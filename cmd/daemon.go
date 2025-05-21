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
	logger := cego.NewLogger()
	ctx := cmd.Context()
	elasticDocumentChan := make(chan *model.ElasticDocument, 50)

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

	go routines.StartDockerEventsChannel(ctx, dockerClient, logger, elasticDocumentChan)

	config, err := routines.FetchConfig(ctx, logger, configEndpoint)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = routines.PushDockerEventsToElastic(ctx, logger, config, elasticDocumentChan)
	if err != nil {
		logger.Error(err.Error())
	}

	<-cmd.Context().Done()

	err = routines.PushDockerBuildxHistoryToElastic(config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func InitDaemon() *cobra.Command {
	var daemon = &cobra.Command{
		Use:   "daemon",
		Short: "Starts daemon that stores docker pull events and hosts a rest api",
		Run:   DaemonRun,
	}

	daemon.Flags().String("config-endpoint", "http://localhost:8000/dopetes.yaml", "Endpoint to fetch config from")
	return daemon
}
