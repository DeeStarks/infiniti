package db

import (
	"fmt"
	"testing"

	"github.com/deestarks/infiniti/config"
	"github.com/deestarks/infiniti/utils"
)

func TestPermissionCreate(t *testing.T) {
	// Load the environment variables
	config.LoadEnv("../../../.env")

	var createdIds []interface{}
	var tests = []struct {
		table_id	int
		method		string
	} {
		{1, PermissionsEnum().Get},
		{1, PermissionsEnum().Create},
		{1, PermissionsEnum().Update},
		{1, PermissionsEnum().Delete},
	}

	dbAdapter, err := NewDBAdapter(
		"postgres", 
		fmt.Sprintf(
			"postgresql://%s:%s@%s:5432/%s?sslmode=disable", 
			config.GetEnv("DB_USER"), 
			config.GetEnv("DB_PASS"), 
			config.GetEnv("DB_HOST"), 
			config.GetEnv("DB_NAME"),
		),
	)
	if err != nil {
		t.Error(err)
	}
	pAdapt := dbAdapter.NewPermissionsAdapter()
	for _, test := range tests {
		permObj, err := pAdapt.Create(map[string]interface{}{
			"table_id": test.table_id,
			"method": test.method,
		})
		if err != nil {
			t.Error(err)
		}
		// Add the id to the list of created ids; this will be used to delete the created permissions
		createdIds = append(createdIds, permObj.Id)

		if permObj.TableId != test.table_id && permObj.Method != test.method {
			t.Errorf("Expected:\n    Table Id: %d\n    Method: %s\nGot:\n    Table Id: %d\n    Method: %s\n", 
				test.table_id, test.method, permObj.TableId, permObj.Method)
		}
	}
	delQuery := fmt.Sprintf("DELETE FROM permissions WHERE id IN (%s)", utils.CreatePlaceholder(len(createdIds)))
	_, err = dbAdapter.db.Exec(delQuery, createdIds...)
	if err != nil {
		t.Error(err)
	}
}

func TestPermissionDelete(t *testing.T) {
	// Load the environment variables
	config.LoadEnv("../../../.env")

	var tests = []struct {
		table_id	int
		method		string
	} {
		{1, PermissionsEnum().Get},
		{1, PermissionsEnum().Create},
		{1, PermissionsEnum().Update},
		{1, PermissionsEnum().Delete},
	}

	dbAdapter, err := NewDBAdapter(
		"postgres", 
		fmt.Sprintf(
			"postgresql://%s:%s@%s:5432/%s?sslmode=disable", 
			config.GetEnv("DB_USER"), 
			config.GetEnv("DB_PASS"), 
			config.GetEnv("DB_HOST"), 
			config.GetEnv("DB_NAME"),
		),
	)
	if err != nil {
		t.Error(err)
	}
	pAdapt := dbAdapter.NewPermissionsAdapter()
	for _, test := range tests {
		createdObj, err := pAdapt.Create(map[string]interface{}{
			"table_id": test.table_id,
			"method": test.method,
		})
		if err != nil {
			t.Error(err)
		}

		obj, err := pAdapt.Delete("id", createdObj.Id)
		if err != nil {
			t.Error(err)
		}
		if obj.Id != createdObj.Id {
			t.Errorf("Expected to delete:\n    Id: %d\nGot:\n    Id: %d\n", createdObj.Id, obj.Id)
		}
	}
}