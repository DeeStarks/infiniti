package services

import (
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type UserAuth struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

func (service *Service) NewUserAuthService() *UserAuth {
	return &UserAuth{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (auth *UserAuth) AuthenticateUser(email string, password string) (bool, error) {
	userAdapter := auth.dbPort.NewUserAdapter()
	user, err := userAdapter.Get("email", email)
	if err != nil {
		return false, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("User not found"),
		}
	}

	if err = auth.corePort.ComparePassword(user.Password, password); err != nil {
		return false, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}
	return true, nil
}