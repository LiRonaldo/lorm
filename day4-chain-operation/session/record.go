package session

import (
	"errors"
	"lorm/day4-chain-operation/caluse"
	"reflect"
)

//插入
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.caluse.Set(caluse.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	//没加... 会多一层层级
	//加了...之后相当于是个数组了，每个数组元素可能又是不一样的类型的值。比如我穿一个切片过去，那么接受方a[0]就是个切片，
	//如果我没有加... 我传的是个切片。本身就是个数组。所以数组加可变参数（本身是数组）所以就多了一层
	s.caluse.Set(caluse.VALUES, recordValues...)
	sql, vars := s.caluse.Build(caluse.INSERT, caluse.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//更新
func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.caluse.Set(caluse.UPDATE, s.RefTable().Name, m)
	sql, vars := s.caluse.Build(caluse.UPDATE, caluse.UPDATE)
	result, err := s.Raw(sql, vars).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//删除
func (s *Session) Delete() (int64, error) {
	s.caluse.Set(caluse.DELETE, s.RefTable().Name)
	sql, vars := s.caluse.Build(caluse.DELETE, caluse.WHERE)
	result, err := s.Raw(sql, vars).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//count
func (s *Session) Count() (int64, error) {
	s.caluse.Set(caluse.COUNT, s.RefTable().Name)
	sql, vars := s.caluse.Build(caluse.COUNT, caluse.WHERE)
	row := s.Raw(sql, vars).QueryRow()
	var num int64
	if err := row.Scan(&num); err != nil {
		return 0, err
	}
	return num, nil
}

//limit
func (s *Session) Limit(num interface{}) *Session {
	s.caluse.Set(caluse.LIMIT, num)
	return s
}

//where
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.caluse.Set(caluse.WHERE, append(append(vars, desc), args))
	return s
}

//order by
func (s *Session) OrderBy(desc string) *Session {
	s.caluse.Set(caluse.ORDERBY, desc)
	return s
}

//查询一条记录
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
