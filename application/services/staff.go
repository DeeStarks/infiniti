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

type StaffResource struct { // Staff with foreign keys
	Id 			int 				`json:"id"`
	FirstName 	string 				`json:"first_name"`
	LastName 	string 				`json:"last_name"`
	Email 		string 				`json:"email"`
	CreatedAt 	string 				`json:"created_at"`
}

type StaffFKResource struct { // Staff with foreign keys
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

func (u *Staff) CreateStaff(data map[string]interface{}) (StaffFKResource, error) {
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
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required fields: %s", missing),
		}
	}

	// Hash the password
	encyptPass, err := u.corePort.HashPassword(data["password"].(string))
	if err != nil {
		return StaffFKResource{}, &utils.RequestError{
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
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}

	// Add the staff to the group
	group, err := u.dbPort.NewGroupAdapter().Get("name", "staff")
	if err != nil {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	_, err = u.dbPort.NewUserGroupAdapter().Create(map[string]interface{}{
		"user_id": staff.Id,
		"group_id": group.Id,
	})
	if err != nil {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Combine and return
	var (
		userRes 	StaffFKResource
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

func (u *Staff) GetStaff(key string, value interface{}) (StaffResource, error) {
	userAdapter := u.dbPort.NewUserAdapter()
	staff, err := userAdapter.Get(key, value)
	if err != nil {
		return StaffResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("staff not found"),
		}
	}
	var staffRes StaffResource
	staffJson, _ := json.Marshal(staff)
	json.Unmarshal(staffJson, &staffRes)
	return staffRes, nil
}

func (u *Staff) ListStaff() ([]StaffFKResource, error) {
	userAdapter := u.dbPort.NewUserAdapter()
	selector := userAdapter.NewUserCustomSelector("groups.name", "staff", "users.id", true).
		Join("user_groups", "user_id", "users", "id", []string{"user_id", "group_id"}).
		Join("groups", "id", "user_groups", "group_id", []string{"id", "name"})
	data := selector.Query()

	var res []StaffFKResource
	for _, user := range data {
		// Prepare user data
		userData := map[string]interface{}{
			"id": user["users__id"],
			"email": user["users__email"],
			"password": user["users__password"],
			"first_name": user["users__first_name"],
			"last_name": user["users__last_name"],
			"created_at": user["users__created_at"],
		}

		// Prepare group data
		groupData := map[string]interface{}{
			"id": user["groups__id"],
			"name": user["groups__name"],
		}

		// Combine and return
		var staffRes StaffFKResource
		userJson, _ := json.Marshal(userData)
		json.Unmarshal(userJson, &staffRes)

		var groupRes GroupResource
		groupJson, _ := json.Marshal(groupData)
		json.Unmarshal(groupJson, &groupRes)

		staffRes.Group = groupRes
		res = append(res, staffRes)
	}
	return res, nil
}


func (u *Staff) UpdateStaff(key string, value interface{}, data map[string]interface{}) (StaffResource, error) {
	for _, v := range []string{"id", "created_at"} { // These fields cannot be updated
		delete(data, v)
	}
	userAdapter := u.dbPort.NewUserAdapter()
	staff, err := userAdapter.Update(key, value, data)
	if err != nil {
		return StaffResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("staff not found"),
		}
	}
	var staffRes StaffResource
	staffJson, _ := json.Marshal(staff)
	json.Unmarshal(staffJson, &staffRes)
	return staffRes, nil
}