package routines

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
)

func FetchConfig(ctx context.Context, l cego.Logger, configEndpoint string) (*model.DopetesConfig, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", configEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching config: %s", err)
	}

	data := &model.DopetesConfig{}
	err = yaml.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	l.Debug("Fetched config successfully")

	return data, nil
}
