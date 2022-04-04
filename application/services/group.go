package services

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
)

type Group struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type GroupLayout struct {
	Id 		int 	`json:"id"`
	Name 	string 	`json:"name"`
}

func (service *Service) NewGroupService() *Group {
	return &Group{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}
