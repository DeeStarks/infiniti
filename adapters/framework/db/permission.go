package db

import (
	"fmt"
	"log"
	"github.com/deestarks/infiniti/lib"
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
	var modelObj PermissionsModel

	colStr, valArr := lib.SplitMap(data)
	query := fmt.Sprintf(`
		INSERT INTO %s ( %s ) VALUES ( %s )
		RETURNING id, table_id, method
	`, pAdpt.tableName, colStr, lib.CreatePlaceholder(len(valArr)))

	dbTx, err := pAdpt.adapter.db.Begin()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	prepQuery, err := dbTx.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer prepQuery.Close()
	err = prepQuery.QueryRow(valArr...).Scan(&modelObj.Id, &modelObj.TableId, &modelObj.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = dbTx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &modelObj, nil
}

func (pAdpt *PermissionsAdapter) Delete(colName string, value interface{}) (*PermissionsModel, error) {
	var (
		modelObj	PermissionsModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, table_id, method
	`, pAdpt.tableName, colName)
	err = pAdpt.adapter.db.QueryRow(query, value).Scan(&modelObj.Id, &modelObj.TableId, &modelObj.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &modelObj, nil
}

func (pAdpt *PermissionsAdapter) Get(colName string, value interface{}) (*PermissionsModel, error) {
	var (
		modelObj	PermissionsModel
		err			error
	)
	query := fmt.Sprintf(`
		SELECT id, table_id, method
		FROM %s
		WHERE %s = $1
	`, pAdpt.tableName, colName)
	err = pAdpt.adapter.db.QueryRow(query, value).Scan(&modelObj.Id, &modelObj.TableId, &modelObj.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &modelObj, nil
}