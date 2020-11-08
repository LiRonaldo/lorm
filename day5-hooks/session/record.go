package session

import (
	"errors"
	"lorm/day5-hooks/caluse"
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
	sql, vars := s.caluse.Build(caluse.UPDATE, caluse.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//删除
func (s *Session) Delete() (int64, error) {
	s.caluse.Set(caluse.DELETE, s.RefTable().Name)
	sql, vars := s.caluse.Build(caluse.DELETE, caluse.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//count
func (s *Session) Count() (int64, error) {
	s.caluse.Set(caluse.COUNT, s.RefTable().Name)
	sql, vars := s.caluse.Build(caluse.COUNT, caluse.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
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
	s.caluse.Set(caluse.WHERE, append(append(vars, desc), args...)...)
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

//查询方法，
//传进去一个切片指针，将查询后的数据保存到切片，不用返回这个切片，
//此处为啥不传切片，而是 要传切片的 指针。因为当吧切片传过去之后，
//如果发生了扩容，是不会影响外部的切片的。也就是外部打印还是发生
//扩容之前的切片。
/*1)destSlice.Type().Elem() 获取切片的单个元素的类型 destType，使用 reflect.New() 方法创建一个 destType 的实例，作为 Model() 的入参，映射出表结构 RefTable()。
2）根据表结构，使用 clause 构造出 SELECT 语句，查询到所有符合条件的记录 rows。
3）遍历每一行记录，利用反射创建 destType 的实例 dest，将 dest 的所有字段平铺开，构造切片 values。
4）调用 rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段。
5）将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中。*/
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	s.caluse.Set(caluse.SELECT, table.Name, table.FieldNames)
	sql, vars := s.caluse.Build(caluse.SELECT, caluse.WHERE, caluse.ORDERBY, caluse.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return nil
	}
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		//调用钩子。实现查询出来的内容，修改。
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
