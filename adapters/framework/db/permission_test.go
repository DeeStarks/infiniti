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
			t.Errorf("Expected:\n\tTable Id: %d\n\tMethod: %s\nGot:\n\tTable Id: %d\n\tMethod: %s\n", 
				test.table_id, test.method, permObj.TableId, permObj.Method)
		}
	}
	delQuery := fmt.Sprintf("DELETE FROM permissions WHERE id IN (%s)", utils.CreatePlaceholder(len(createdIds)))
	_, err = dbAdapter.db.Exec(delQuery, createdIds...)
	if err != nil {
		t.Error(err)
	}
}

func TestPermissionUpdate(t *testing.T) {
	// Load the environment variables
	config.LoadEnv("../../../.env")

	var createdIds []interface{}
	var testCases = [4][4]struct {
		table_id	int
		method		string
	} {
		{
			{1, PermissionsEnum().Get},
			{1, PermissionsEnum().Create},
			{1, PermissionsEnum().Update},
			{1, PermissionsEnum().Delete},
		},
		{
			{2, PermissionsEnum().Get},
			{4, PermissionsEnum().Update},
			{3, PermissionsEnum().Delete},
			{5, PermissionsEnum().Create},
		},
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
	for i := 0; i < len(testCases[0]); i++ {
		createdObj, err := pAdapt.Create(map[string]interface{}{
			"table_id": testCases[0][i].table_id,
			"method": testCases[0][i].method,
		})
		if err != nil {
			t.Error(err)
		}

		updatedObj, err := pAdapt.Update("id", createdObj.Id, map[string]interface{}{
			"table_id": testCases[1][i].table_id,
			"method": testCases[1][i].method,
		})
		if err != nil {
			t.Error(err)
		}
		// Add the id to the list of created ids; this will be used to delete the created permissions
		createdIds = append(createdIds, updatedObj.Id)

		if updatedObj.TableId != testCases[1][i].table_id && updatedObj.Method != testCases[1][i].method {
			t.Errorf("Expected:\n\tTable Id: %d\n\tMethod: %s\nGot:\n\tTable Id: %d\n\tMethod: %s\n", 
				testCases[1][i].table_id, testCases[1][i].method, updatedObj.TableId, updatedObj.Method)
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
			t.Errorf("Expected to delete:\n\tId: %d\nGot:\n\tId: %d\n", createdObj.Id, obj.Id)
		}
	}
}