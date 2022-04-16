package constants

import (
	"database/sql"
	"fmt"
	"strings"
)

// Custom methods
type CustomSelector struct {
	adapter 	*sql.DB
	tableName 	string
	cols 		[]string
	condition 	struct {
		col 	string
		val 	interface{}
	}
	joins 		[]struct {
		table 	string
		col 	string
		on 		struct {
			table 	string
			col 	string
		}
	}
	order 		struct {
		col 	string
		asc 	bool
	}
}

func NewCustomSelector(adapter *sql.DB, tableName string, fields []string, col string, value interface{}, order string, asc bool) *CustomSelector {
	var cs CustomSelector
	cs.adapter = adapter
	cs.tableName = tableName
	cs.cols = fields
	cs.condition.col = col
	cs.condition.val = value
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

	if c.condition.col != "" {
		condition = fmt.Sprintf(`
			WHERE %s = $1
		`, c.condition.col)
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

// adapter: Database Adapter
func (c *CustomSelector) Query() []map[string]interface{} {
	var (
		rows 		*sql.Rows
		err 		error
		result 		[]map[string]interface{}
	)
	rows, err = c.adapter.Query(c.String(), c.condition.val)
	var cols []string
	for _, col := range c.cols {
		cols = append(cols, strings.Replace(col, ".", "__", -1)) // replace "." with "__"
	}
	if err != nil {
		return nil
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
	return result
}