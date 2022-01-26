package db

import (
	"fmt"
	"log"
	"github.com/deestarks/infiniti/pkg"
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
	var (
		modelObj	 	PermissionsModel
		colStr			string
		valArr 			[]interface{}
	)

	for col, val := range data {
		colStr += col + ", "
		valArr = append(valArr, val)
	}
	colStr = colStr[:len(colStr)-2] // remove the last ", "
	
	dbTx, err := pAdpt.adapter.db.Begin()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	prepQuery, err := dbTx.Prepare(fmt.Sprintf(`
		INSERT INTO %s ( %s ) VALUES ( %s )
		RETURNING id, table_id, method
	`, pAdpt.tableName, colStr, pkg.CreatePlaceholder(len(valArr))))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer prepQuery.Close()
	prepQuery.QueryRow(valArr...).Scan(&modelObj.Id, &modelObj.TableId, &modelObj.Method)

	err = dbTx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &modelObj, nil
}