package caluse

import (
	"strings"
)

type Type int

const (
	INSERT = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

type Caluse struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Caluse) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}
func (c *Caluse) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
