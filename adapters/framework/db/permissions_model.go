package db

import (
	"fmt"
	"log"
	"github.com/deestarks/infiniti/pkg"
)

type (
	PermissionsModel struct {
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

func (adpt *DBAdapter) NewPermissionsModel() *PermissionsModel {
	return &PermissionsModel{
		adapter: adpt,
		tableName: "permissions",
	}
}

func (model *PermissionsModel) Create(data map[string]interface{}) error {
	colStr := ""
	var valArr []interface{}

	for col, val := range data {
		colStr += col + ", "
		valArr = append(valArr, val)
	}
	colStr = colStr[:len(colStr)-2] // remove the last ", "
	
	_, err := model.adapter.db.Exec(fmt.Sprintf(`
		INSERT INTO %s ( %s ) VALUES ( %s )
	`, model.tableName, colStr, pkg.CreatePlaceholder(len(valArr))), valArr...)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}