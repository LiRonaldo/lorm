package caluse

import (
	"lorm/day-6-transaction/log"
	"testing"
)

func TestCaluse_Build(t *testing.T) {
	caluse := &Caluse{}
	caluse.Set(LIMIT, 3)
	//接受一个可变参数，SELECT 后的属于可变参数。那么可变参数【0】是user ，可变参数【1】是[]string{"name"}，interface强转的话，是要根据具体的类型。
	//可变参数的第二个类型是string数组，所以可以强转成string数组。
	caluse.Set(SELECT, "User", []string{"name"})
	caluse.Set(WHERE, "Name=?", "Tom")
	caluse.Set(ORDERBY, "Age ASC")
	sql, vars := caluse.Build(SELECT, WHERE, ORDERBY, LIMIT)
	log.Info(sql)
	log.Info(vars)
}
