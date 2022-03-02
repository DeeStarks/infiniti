package db

import (
	"fmt"
	"log"

	"github.com/deestarks/infiniti/utils"
)

type (
	PermissionsModel struct {
		Id			int 	`json:"id"`
		TableId		int 	`json:"table_id"`
		Method		string 	`json:"method"`
	}

	PermissionsAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}

	_permissionsEnum struct {
		Get, Create, Update, Delete 	string
	}
)

func PermissionsEnum() _permissionsEnum {
	return _permissionsEnum{
		Get: "get",
		Create: "post",
		Update: "put",
		Delete: "delete",
	}
}

func (adpt *DBAdapter) NewPermissionsAdapter() *PermissionsAdapter {
	return &PermissionsAdapter{
		adapter: adpt,
		tableName: "permissions",
	}
}

func (pAdpt *PermissionsAdapter) Create(data map[string]interface{}) (*PermissionsModel, error) {
	var permission PermissionsModel

	mToS := utils.MapToStructSlice(data)
	var (
		colStr		string
		valArr		[]interface{}
	)
	for i, s := range mToS {
		colStr += s.Key + ", "
		valArr = append(valArr, s.Value)
		if i == len(mToS)-1 {
			colStr = colStr[:len(colStr)-2] // remove the last ", "
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s ( %s ) VALUES ( %s )
		RETURNING id, table_id, method
	`, pAdpt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := pAdpt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.TableId, &permission.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

func (pAdpt *PermissionsAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*PermissionsModel, error) {
	var (
		permission 	PermissionsModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, table_id, method
	`, pAdpt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1) // "len(data)+1" creates a placeholder for the value of the column to be updated.

	// Add the value of the column to be updated to the end of the array of values.
	valArr = append(valArr, colValue)

	err := pAdpt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.TableId, &permission.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

func (pAdpt *PermissionsAdapter) Delete(colName string, value interface{}) (*PermissionsModel, error) {
	var (
		permission	PermissionsModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, table_id, method
	`, pAdpt.tableName, colName)
	err = pAdpt.adapter.db.QueryRow(query, value).Scan(&permission.Id, &permission.TableId, &permission.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

// "Get" returns a single permission
func (pAdpt *PermissionsAdapter) Get(colName string, value interface{}) (*PermissionsModel, error) {
	var (
		permission	PermissionsModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, table_id, method
		FROM %s
		WHERE %s = $1
	`, pAdpt.tableName, colName)
	// `, pAdpt.tableName, colName)
	err = pAdpt.adapter.db.QueryRow(query, value).Scan(
		&permission.Id, &permission.TableId, &permission.Method,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

func (pAdpt *PermissionsAdapter) Filter(colName string, value interface{}) (*[]PermissionsModel, error) {
	var (
		permissions	[]PermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, table_id, method FROM %s WHERE %s = $1", pAdpt.tableName, colName)
	rows, err := pAdpt.adapter.db.Query(query, value)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission PermissionsModel
		err := rows.Scan(
			&permission.Id, &permission.TableId, &permission.Method,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}

// "List" returns all permissions
func (pAdpt *PermissionsAdapter) List() (*[]PermissionsModel, error) {
	var (
		permissions	[]PermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, table_id, method FROM %s", pAdpt.tableName)
	rows, err := pAdpt.adapter.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission PermissionsModel
		err := rows.Scan(
			&permission.Id, &permission.TableId, &permission.Method,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}