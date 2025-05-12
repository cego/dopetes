package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
)

type PublishHandler struct {
	logger   cego.Logger
	model    *model.Model
	renderer *cego.Renderer
	ctx      context.Context
}

func NewPublishHandler(ctx context.Context, logger cego.Logger, model *model.Model, renderer *cego.Renderer) *PublishHandler {
	return &PublishHandler{
		ctx:      ctx,
		logger:   logger,
		model:    model,
		renderer: renderer,
	}
}

type PublishHandlerReqBody struct {
	ElasticsearchHost   string `json:"elasticsearch_host"`
	ElasticsearchApiKey string `json:"elasticsearch_api_key"`
}

func (h *PublishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := &PublishHandlerReqBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.logger.Error(err.Error())
		h.renderer.Text(w, 400, err.Error())
		return
	}

	if body.ElasticsearchApiKey == "" {
		h.renderer.Text(w, 400, "elasticsearch_api_key is required")
		return
	}

	if body.ElasticsearchHost == "" {
		h.renderer.Text(w, 400, "elasticsearch_host is required")
		return
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{body.ElasticsearchHost},
		APIKey:    body.ElasticsearchApiKey,
	})
	if err != nil {
		h.renderer.Text(w, 500, err.Error())
		return
	}

	var bulkErrors []error

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "filebeat-prod-8.18.1",
		Client: es,

		OnError: func(ctx context.Context, err error) {
			bulkErrors = append(bulkErrors, err)
		},
	})
	if err != nil {
		h.renderer.Text(w, 500, err.Error())
		return
	}

	for _, dockerPullEvent := range h.model.GetDockerPullEvents() {
		var b []byte
		b, err = json.Marshal(&dockerPullEvent)
		if err != nil {
			h.renderer.Text(w, 500, err.Error())
			return
		}

		err = bi.Add(h.ctx, esutil.BulkIndexerItem{
			Action: "create",
			Body:   bytes.NewReader(b),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					h.logger.Error(err.Error())
				} else {
					h.logger.Error(fmt.Sprintf("%s %s", res.Error.Type, res.Error.Reason))
				}
			},
		})
		if err != nil {
			h.renderer.Text(w, 500, err.Error())
			return
		}
	}

	err = bi.Close(h.ctx)
	if err != nil {
		h.renderer.Text(w, 500, err.Error())
		return
	}

	if len(bulkErrors) > 0 {
		errorText := ""
		for bulkError := range bulkErrors {
			errorText += bulkErrors[bulkError].Error()
		}
		h.renderer.Text(w, 500, errorText)
		return
	}

	biStats := bi.Stats()
	h.renderer.Text(w, 200, fmt.Sprintf("%d docker pull events have been indexed in %s\n", biStats.NumCreated, body.ElasticsearchHost))
}
