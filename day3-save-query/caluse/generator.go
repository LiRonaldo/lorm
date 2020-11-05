package caluse

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby

}

func _orderby(values ...interface{}) (s string, i []interface{}) {
	return fmt.Sprintf("order by %s", values[0]), []interface{}{}
	return
}

func _where(values ...interface{}) (s string, i []interface{}) {
	sql, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", sql), vars
}

func _limit(values ...interface{}) (s string, i []interface{}) {
	return "limit ?", values
}
func _select(values ...interface{}) (s string, i []interface{}) {
	//SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("select %v from %s", fields, tableName), []interface{}{}
}

func _values(values ...interface{}) (s string, i []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _insert(values ...interface{}) (s string, i []interface{}) {
	//插入	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fiedls := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fiedls), []interface{}{}
}

//拼接,
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}
