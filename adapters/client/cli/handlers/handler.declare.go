package handlers

import (
	"github.com/deestarks/infiniti/application/services"
)

type Handler struct {
	appPort		services.AppServicesPort
	version		string
}

func NewCLIHandler(appPort services.AppServicesPort) *Handler {
	return &Handler{
		appPort: appPort,
		version: "0.0.1",
	}
}