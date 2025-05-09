package routines

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
	"github.com/ory/graceful"
)

type PostToElastic struct {
	ElasticsearchUrl    string `json:"elasticsearch_url" validate:"required"`
	ElasticsearchApiKey string `json:"elasticsearch_api_key" validate:"required"`
}

func StartRestApi(ctx context.Context, m *model.Model, listenPort int, renderer *cego.Renderer, elasticsearcBulkUrl string) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/post-to-elastic", func(w http.ResponseWriter, r *http.Request) {
		body := &PostToElastic{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			slog.Error(err.Error())
			renderer.Text(w, 400, err.Error())
			return
		}

		if body.ElasticsearchApiKey == "" {
			renderer.Text(w, 400, "elasticsearch_api_key is required")
			return
		}

		if body.ElasticsearchUrl == "" {
			renderer.Text(w, 400, "elasticsearch_url is required")
			return
		}

		renderer.Text(w, 200, fmt.Sprintf("%d docker pull events have been send to %s successfully\n", len(m.GetDockerPullEvents()), elasticsearcBulkUrl))
	})

	server := graceful.WithDefaults(&http.Server{
		Addr:         fmt.Sprintf(":%d", listenPort),
		Handler:      mux,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	})

	slog.Info(fmt.Sprintf("Listening on port %v", listenPort))
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		slog.Error("Failed to gracefully shutdown")
		os.Exit(1)
	}
	slog.Info("Server was shutdown gracefully")

}
