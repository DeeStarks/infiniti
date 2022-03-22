package handlers

import (
	"github.com/deestarks/infiniti/application/services"
)

type Handler struct {
	appPort		services.AppServicesPort
}

func NewHandler(appPort services.AppServicesPort) *Handler {
	return &Handler{
		appPort: appPort,
	}
}