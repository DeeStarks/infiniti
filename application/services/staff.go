package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type Staff struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type StaffResource struct {
	Id 			int 				`json:"id"`
	FirstName 	string 				`json:"first_name"`
	LastName 	string 				`json:"last_name"`
	Email 		string 				`json:"email"`
	Group 		GroupResource 		`json:"group"`
	CreatedAt 	string 				`json:"created_at"`
}

// Initialize the staff service
func (service *Service) NewStaffService() *Staff {
	return &Staff{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (u *Staff) CreateStaff(data map[string]interface{}) (StaffResource, error) {
	adapter := u.dbPort.NewUserAdapter()
	// First check that all required fields are present
	required := []string{"email", "password", "first_name", "last_name"}
	var missing string

	for _, field := range required {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return StaffResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required fields: %s", missing),
		}
	}

	// Hash the password
	encyptPass, err := u.corePort.HashPassword(data["password"].(string))
	if err != nil {
		return StaffResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Prepare staff data
	staffData := map[string]interface{}{
		"email": data["email"],
		"password": encyptPass,
		"first_name": data["first_name"],
		"last_name": data["last_name"],
	}
	
	// Create the staff
	staff, err := adapter.Create(staffData)
	if err != nil {
		return StaffResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}

	// Add the staff to the group
	group, err := u.dbPort.NewGroupAdapter().Get("name", "staff")
	if err != nil {
		return StaffResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	_, err = u.dbPort.NewUserGroupAdapter().Create(map[string]interface{}{
		"user_id": staff.Id,
		"group_id": group.Id,
	})
	if err != nil {
		return StaffResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Combine and return
	var (
		userRes 	StaffResource
		groupRes 	GroupResource
	)
	// User
	userJson, _ := json.Marshal(staff)
	json.Unmarshal(userJson, &userRes)
	// Group
	groupJson, _ := json.Marshal(group)
	json.Unmarshal(groupJson, &groupRes)

	// Combine and return
	userRes.Group = groupRes
	return userRes, nil
}