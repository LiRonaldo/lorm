package session

import (
	"lorm/day3-save-query/caluse"
)

//
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.caluse.Set(caluse.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	//没加... 加了...之后相当于是个数组了，每个数组元素可能又是不一样的类型的值。比如我穿一个切片过去，那么接受方a[0]就是个切片，
	//如果我没有加... 那么接受方a[0].b[0]才是切片。
	s.caluse.Set(caluse.VALUES, recordValues...)
	sql, vars := s.caluse.Build(caluse.INSERT, caluse.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
