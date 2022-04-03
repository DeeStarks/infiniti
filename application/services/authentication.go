package services

import (
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/application/services/constants"
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

func (auth *UserAuth) AuthenticateUser(data map[string]interface{}) (constants.ServiceStructReturnType, error) {
	requires := []string{"email", "password"}
	var missing string

	for _, field := range requires {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return nil, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required fields: '%s'", missing[2:]),
		}
	}

	// Get the user
	email := data["email"].(string)
	password := data["password"].(string)

	userAdapter := auth.dbPort.NewUserAdapter()
	user, err := userAdapter.Get("email", email)
	if err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("User not found"),
		}
	}

	if err = auth.corePort.ComparePassword(user.Password, password); err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}

	// Get user's group
	userGroup, err := auth.dbPort.NewUserGroupAdapter().Get("user_id", user.Id)
	if err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}
	group, err := auth.dbPort.NewGroupAdapter().Get("id", userGroup.GroupId)
	if err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Serialization and return
	serializedUser := constants.ServiceStructReturnType(utils.StructToMap(user))
	serializedUser["group_name"] = group.Name
	return serializedUser, nil
}