package db

import (
	"fmt"
	"log"
)

type DBModel struct {
	adapter		*DBAdapter
}

func (adpt *DBAdapter) NewDBModel() *DBModel {
	return &DBModel{adapter: adpt}
}

func (model *DBModel) CreateTable(tableName string, columns map[string]string) error {
	log.Printf("Creating table \"%s\"\n", tableName)
	columnsStr := ""
	for col, colType := range columns {
		columnsStr += fmt.Sprintf("%s %s, ", col, colType)
	}
	columnsStr = columnsStr[:len(columnsStr)-2]
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id BIGSERIAL PRIMARY KEY, %s)", tableName, columnsStr)
	_, err := model.adapter.db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = model.adapter.db.Exec(fmt.Sprintf("INSERT INTO db_tables (table_name) VALUES ('%s')", tableName))
	log.Printf("Created table \"%s\"\n", tableName)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (model *DBModel) DropTable(tableName string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err := model.adapter.db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = model.adapter.db.Exec(fmt.Sprintf("DELETE FROM db_tables WHERE table_name = '%s'", tableName))
	log.Printf("Dropped table \"%s\"\n", tableName)
	if err != nil {
		log.Fatal(err)
	}
	return err
}