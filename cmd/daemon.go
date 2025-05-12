package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cego/dopetes/handlers"
	"github.com/cego/dopetes/model"
	"github.com/cego/dopetes/routines"
	"github.com/cego/go-lib"
	"github.com/docker/docker/client"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
)

func DaemonRun(cmd *cobra.Command, _ []string) {
	logger := cego.NewLogger()
	m := model.New()
	r := cego.NewRenderer(slog.Default())

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	routines.StartDockerEventsRoutine(cmd.Context(), m, cli, logger)

	mux := http.NewServeMux()
	mux.Handle("POST /api/publish", handlers.NewPublishHandler(cmd.Context(), logger, m, r))
	mux.Handle("POST /api/clear", handlers.NewClearHandler(cmd.Context(), logger, m, r))
	port, err := cmd.Flags().GetInt("listen-port")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	server := graceful.WithDefaults(&http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		BaseContext:  func(_ net.Listener) context.Context { return cmd.Context() },
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	})
	logger.Info(fmt.Sprintf("Listening on port %d", port))
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		logger.Error("Failed to gracefully shutdown %s", err.Error())
		os.Exit(1)
	}
	logger.Info("Server was shutdown gracefully")
}

func InitDaemon() *cobra.Command {
	var daemon = &cobra.Command{
		Use:   "daemon",
		Short: "Starts daemon that stores docker pull events and hosts a rest api",
		Run:   DaemonRun,
	}

	daemon.Flags().Int("listen-port", 2900, "Listeport for the webserver")
	return daemon
}
