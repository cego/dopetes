package cmd

import (
	"io"
	"net/http"
	"os"

	"github.com/cego/go-lib"
	"github.com/spf13/cobra"
)

func PublishClear(cmd *cobra.Command, _ []string) {
	logger := cego.NewLogger()
	logger.Debug("Sending POST to /api/clear")

	daemonUrl, err := cmd.Flags().GetString("daemon-url")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Create a HTTP post request
	r, err := http.NewRequest("POST", daemonUrl+"/api/clear", nil)
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

func InitClear() *cobra.Command {
	var c = &cobra.Command{
		Use:   "clear",
		Short: "Send POST to /api/clear",
		Run:   PublishClear,
	}
	c.Flags().String("daemon-url", "http://127.0.0.1:2900", "Dopetes daemon webserver url")
	return c
}
