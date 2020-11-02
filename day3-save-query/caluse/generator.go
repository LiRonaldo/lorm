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
	/*generators[VALUES]=_values
	generators[SELECT]=_select
	generators[LIMIT]=_limit
	generators[WHERE]=_where
	generators[ORDERBY]=_orderby*/

}

func _orderby(values ...interface{}) (s string, i []interface{}) {
	return
}

func _where(values ...interface{}) (s string, i []interface{}) {
	return
}

func _limit(values ...interface{}) (s string, i []interface{}) {
	return
}

func _select(values ...interface{}) (s string, i []interface{}) {
	return
}

func _values(values ...interface{}) (s string, i []interface{}) {
	return
}

func _insert(values ...interface{}) (s string, i []interface{}) {
	//插入
	tableName := values[0]
	fiedls := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fiedls), []interface{}{}
}
