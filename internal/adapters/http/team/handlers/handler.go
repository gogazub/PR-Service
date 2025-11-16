package teamhandlers

import (
	"PRService/internal/app"

	"go.uber.org/zap"
)
type Handler struct {
	*app.Services
	logger *zap.SugaredLogger
}

// NewHandler returns new Handler.
func NewHandler(app *app.Services, logger *zap.SugaredLogger) *Handler {
	return &Handler{app, logger}
}