package application

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/services"
)

type Application struct {
	db	 		db.DBPort
	services 	services.AppServicesPort
}

func NewApplication(dbPort db.DBPort) *Application {
	return &Application{
		db: 		dbPort,
		services: 	services.NewServices(dbPort),
	}
}