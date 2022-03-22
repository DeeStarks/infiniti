package application

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/services"
)

type Application struct {
	DBPort	 	db.DBPort
	Services 	services.AppServicesPort
}

func NewApplication(dbPort db.DBPort) *Application {
	return &Application{
		DBPort: 	dbPort,
		Services: 	services.NewServices(dbPort),
	}
}