package userhandler

import (
	"PRService/internal/app"
	"PRService/pkg/logger"
)

type Handler struct {
	*app.Services
	logger *logger.Logger
}

// NewHandler returns new Handler.
func NewHandler(app *app.Services, logger *logger.Logger) *Handler {
	return &Handler{app, logger}
}
