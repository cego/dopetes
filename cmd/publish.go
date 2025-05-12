package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/cego/dopetes/handlers"
	"github.com/cego/go-lib"
	"github.com/spf13/cobra"
)

func PublishRun(cmd *cobra.Command, _ []string) {
	logger := cego.NewLogger()
	logger.Debug("Sending POST to /api/publish")

	daemonUrl, err := cmd.Flags().GetString("daemon-url")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	elasticHost, err := cmd.Flags().GetString("elasticsearch-host")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	elasticApiKey, err := cmd.Flags().GetString("elasticsearch-api-key")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	body := handlers.PublishHandlerReqBody{
		ElasticsearchHost:   elasticHost,
		ElasticsearchApiKey: elasticApiKey,
	}
	marshalled, err := json.Marshal(body)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Create a HTTP post request
	r, err := http.NewRequest("POST", daemonUrl+"/api/publish", bytes.NewBuffer(marshalled))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	resBody, e := io.ReadAll(res.Body)
	if e != nil {
		logger.Error(e.Error())
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		logger.Error(res.Status + "\n" + string(resBody))
		os.Exit(1)
	}
	logger.Info(res.Status + "\n" + string(resBody))
}

func InitPublish() *cobra.Command {
	var publish = &cobra.Command{
		Use:   "publish",
		Short: "Send POST to /api/publish",
		Run:   PublishRun,
	}
	publish.Flags().String("daemon-url", "http://127.0.0.1:2900", "Dopetes daemon webserver url")
	publish.Flags().String("elasticsearch-host", "http://elasticsearch:9200", "Elasticsearch host")
	publish.Flags().String("elasticsearch-api-key", "", "Elasticsearch Api Key")
	return publish
}
