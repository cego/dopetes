package cmd

import (
	"log/slog"

	"github.com/cego/dopetes/model"
	"github.com/cego/dopetes/routines"
	"github.com/cego/go-lib"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func DaemonRun(cmd *cobra.Command, args []string) {
	elasticsearchBulkUrl, _ := cmd.Flags().GetString("elasticsearch-bulk-url")
	m := model.New()
	r := cego.NewRenderer(slog.Default())

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	routines.StartDockerEventListener(cmd.Context(), m, cli)
	routines.StartRestApi(cmd.Context(), m, 2900, r, elasticsearchBulkUrl)
}

var daemon = &cobra.Command{
	Use:   "daemon",
	Short: "Start daemon",
	Long:  `Starts daemon that stores docker pull events and a rest api`,
	Run:   DaemonRun,
}
