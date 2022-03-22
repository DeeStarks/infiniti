package application

import (
	"github.com/deestarks/infiniti/ports"
	"github.com/deestarks/infiniti/application/services"
)

type Application struct {
	db	 		ports.DBPort
	services 	ports.AppServicesPort
}

func NewApplication(dbPort ports.DBPort) *Application {
	return &Application{
		db: 		dbPort,
		services: 	services.NewServices(dbPort),
	}
}