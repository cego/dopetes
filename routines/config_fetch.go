package routines

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
)

func FetchConfig(ctx context.Context, m *model.Model, logger cego.Logger, httpClient *http.Client, configEndpoint string) {
	req, err := http.NewRequestWithContext(ctx, "GET", configEndpoint, nil)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	data := &model.DopetesConfig{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("Fetched config successfully")
	m.SetElasticsearchConfig(data.Elasticsearch)
}
