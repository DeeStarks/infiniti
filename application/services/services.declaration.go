package services

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
)

type Service struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

func NewServices(dbPort db.DBPort) AppServicesPort {
	return &Service{
		dbPort: 	dbPort,
		corePort: 	core.NewCoreApplication(),
	}
}