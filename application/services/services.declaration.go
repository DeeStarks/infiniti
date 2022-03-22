package services

import (
	"github.com/deestarks/infiniti/ports"
	"github.com/deestarks/infiniti/application/core"
)

type Service struct {
	dbPort 		ports.DBPort
	corePort 	ports.CoreAppPort
}

func NewServices(dbPort ports.DBPort) ports.AppServicesPort {
	return Service{
		dbPort: 	dbPort,
		corePort: 	core.NewCoreApplication(),
	}
}