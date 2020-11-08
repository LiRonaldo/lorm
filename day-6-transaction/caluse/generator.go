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
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func _count(values ...interface{}) (s string, i []interface{}) {
	return _select(values[0], []string{"count(*)"})
}

func _delete(values ...interface{}) (s string, i []interface{}) {
	return fmt.Sprintf("delete from %s", values[0]), []interface{}{}
}

func _update(values ...interface{}) (s string, i []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vals []interface{}
	for key, val := range m {
		keys = append(keys, key+"?")
		values = append(vals, val)
	}
	return fmt.Sprintf("update %s set %s", tableName, strings.Join(keys, ",")), values
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
