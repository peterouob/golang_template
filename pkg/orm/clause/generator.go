package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators = map[Type]generator{}

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO [tableName] [values]
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ",")
	query := fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields)
	return query, []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// Values [value1], [value2], ...

	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		if bindStr == "" {
			bindStr = genBindVars(len(value.([]interface{})))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i != len(values)-1 {
			sql.WriteString(", ")
		}
		vars = append(vars, value.([]interface{})...)
	}

	return sql.String(), vars
}
func _select(values ...interface{}) (string, []interface{}) {
	// SEELCT [Values] FROM [TableName]
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ",")
	query := fmt.Sprintf("SELECT %v FROM %s", fields, tableName)
	return query, []interface{}{}
}
func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT [values]
	return "LIMIT ?", values
}
func _where(values ...interface{}) (string, []interface{}) {
	// WHERE [TableName]
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %v", desc), vars
}
func _orderby(values ...interface{}) (string, []interface{}) {
	// ORDER BY [TableName]
	return fmt.Sprintf("ORDER BY %v", values[0]), []interface{}{}
}
