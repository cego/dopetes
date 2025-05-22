package cmd

import (
	"fmt"
	"os"

	"github.com/cego/dopetes/model"
	"github.com/cego/dopetes/routines"
	"github.com/cego/go-lib"
	"github.com/docker/docker/client"
	"github.com/elastic/go-elasticsearch/v9"
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

	elasticClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: config.Elasticsearch.Hosts,
		APIKey:    config.Elasticsearch.ApiKey,
		Username:  config.Elasticsearch.Username,
		Password:  config.Elasticsearch.Password,
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	go routines.PushDockerEventsToElastic(ctx, logger, config, elasticClient, elasticDocumentChan)

	dockerBuildxHistoryState := model.NewDockerBuildxHistoryState()
	routines.StartDockerBuildxHistoryInterval(logger, config, elasticClient, dockerBuildxHistoryState)

	<-cmd.Context().Done()

	err = routines.PushDockerBuildxHistoryToElastic(logger, config, elasticClient, dockerBuildxHistoryState)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	fmt.Println("I made it")
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
