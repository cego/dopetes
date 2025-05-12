package handlers

import (
	"context"
	"net/http"

	"github.com/cego/dopetes/model"
	"github.com/cego/go-lib"
)

type ClearHandler struct {
	logger   cego.Logger
	model    *model.Model
	renderer *cego.Renderer
	ctx      context.Context
}

func NewClearHandler(ctx context.Context, logger cego.Logger, model *model.Model, renderer *cego.Renderer) *ClearHandler {
	return &ClearHandler{
		ctx:      ctx,
		logger:   logger,
		model:    model,
		renderer: renderer,
	}
}
func (h *ClearHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.logger.Debug("/api/clear been hit")
	h.model.ClearDockerPullEvents()
	h.renderer.Text(w, 200, "all docker pull events have been clered")
}
