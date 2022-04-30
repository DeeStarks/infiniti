package constants

import (
	"database/sql"
	"fmt"
	"strings"
)

// Custom methods
type CustomSelector struct {
	adapter 	*sql.DB 		// Database Adapter
	tableName 	string			// Base table name (e.g. "users" in "SELECT * FROM users")
	cols 		[]string		// Columns to select
	condition 	[]struct {		// Condition to select (e.g. "WHERE users.id = 1")
		col 	string			// Column name (e.g. "users.id")
		val 	interface{}		// Value (e.g. 1)
	}
	joins 		[]struct {		// Joins to select (e.g. "LEFT JOIN groups ON groups.user_id = users.id")
		table 	string			// Table name (e.g. "groups")
		col 	string			// Column name (e.g. "groups.user_id")
		on 		struct {
			table 	string		// Table name to check (e.g. "users")
			col 	string		// Column name to check (e.g. "users.id")
		}
	}
	order 		struct {		// Order of the result (e.g. "ORDER BY users.id ASC")
		col 	string			// Column name (e.g. "users.id")
		asc 	bool			// Ascending (true) or descending (false)
	}
}

func NewCustomSelector(adapter *sql.DB, tableName string, fields []string, conditions map[string]interface{}, order string, asc bool) *CustomSelector {
	var cs CustomSelector
	cs.adapter = adapter
	cs.tableName = tableName
	cs.cols = fields
	for col, val := range conditions {
		cs.condition = append(cs.condition, struct {
			col 	string
			val 	interface{}
		}{col, val})
	}
	cs.order.col = order
	cs.order.asc = asc
	return &cs
}

func (c *CustomSelector) String() string {
	var (
		cols 		string
		joins 		string
		condition 	string
		order 		string
	)
	for _, col := range c.cols {
		cols += col + ", "
	}
	cols = cols[:len(cols)-2] // remove the last ", "

	for _, join := range c.joins {
		joins += fmt.Sprintf(`
			LEFT JOIN %s ON %s.%s = %s.%s
		`, join.table, join.table, join.col, join.on.table, join.on.col)
	}

	if len(c.condition) > 0 {
		condition = "WHERE "
		for i, cond := range c.condition {
			condition += fmt.Sprintf("%s = $%d AND ", cond.col, i+1)
		}
		condition = condition[:len(condition)-5] // remove the last " AND "
	}

	if c.order.col != "" {
		orderString := "ASC"
		if !c.order.asc {
			orderString = "DESC"
		}
		order = fmt.Sprintf(`
			ORDER BY %s %s
		`, c.order.col, orderString)
	}

	return fmt.Sprintf(`
		SELECT %s FROM %s %s %s %s
	`, cols, c.tableName, joins, condition, order)
}

// table: table name
// col: column name to check
// onTable: table to check
// onCol: column name to check on "onTable", to match with col
// fields: columns to select
func (c *CustomSelector) Join(table string, col string, onTable string, onCol string, fields []string) *CustomSelector {
	c.joins = append(c.joins, struct {
		table 	string
		col 	string
		on 		struct {
			table 	string
			col 	string
		}
	}{
		table: table,
		col: col,
		on: struct {
			table 	string
			col 	string
		}{
			table: onTable,
			col: onCol,
		},
	})
	
	for _, field := range fields {
		c.cols = append(c.cols, fmt.Sprintf("%s.%s", table, field))
	}
	return c
}

func (c *CustomSelector) Query() ([]map[string]interface{}, error) {
	var (
		rows 	*sql.Rows
		err 	error
		result 	[]map[string]interface{}
	)
	conditionValues := make([]interface{}, len(c.condition))
	for i, cond := range c.condition {
		conditionValues[i] = cond.val
	}
	rows, err = c.adapter.Query(c.String(), conditionValues...)
	var cols []string
	for _, col := range c.cols {
		cols = append(cols, strings.Replace(col, ".", "__", -1)) // replace "." with "__"
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			data		= map[string]interface{}{}
			fields 		= make([]interface{}, len(cols))
			fieldPtrs	= make([]interface{}, len(fields))
		)
		for i := range fields {
			fieldPtrs[i] = &fields[i]
		}
		rows.Scan(fieldPtrs...)

		for i, colName := range cols {
			data[colName] = fields[i]
		}
		result = append(result, data)
	}
	return result, nil
}