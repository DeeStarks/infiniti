package db

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	var createdIds string
	var tests = []struct {
		table_id	int
		method		string
	} {
		{1, PermissionsEnum().Get},
		{1, PermissionsEnum().Create},
		{1, PermissionsEnum().Update},
		{1, PermissionsEnum().Delete},
	}

	dbAdapter, err := NewDBAdapter("postgres", "postgresql://postgres:infiniti@localhost:5432/infiniti?sslmode=disable")
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
		createdIds += fmt.Sprintf("%d, ", permObj.Id)
		createdIds = createdIds[:len(createdIds)-2]

		if permObj.TableId != test.table_id && permObj.Method != test.method {
			t.Errorf("Expected:\n    Table Id: %d\n    Method: %s\nGot:\n    Table Id: %d\n    Method: %s\n", 
				test.table_id, test.method, permObj.TableId, permObj.Method)
		}
	}

	_, err = dbAdapter.db.Exec("DELETE FROM permissions WHERE id IN ($1)", createdIds)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	var tests = []struct {
		table_id	int
		method		string
	} {
		{1, PermissionsEnum().Get},
		{1, PermissionsEnum().Create},
		{1, PermissionsEnum().Update},
		{1, PermissionsEnum().Delete},
	}

	dbAdapter, err := NewDBAdapter("postgres", "postgresql://postgres:infiniti@localhost:5432/infiniti?sslmode=disable")
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