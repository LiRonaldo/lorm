package session

import (
	"fmt"
	"lorm/day4-chain-operation/caluse"
	"lorm/day4-chain-operation/log"
	"lorm/day4-chain-operation/schema"
	"reflect"
	"strings"
)

//将表名赋给s.reftable
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

//创建表，根据传入的strcut 返回一个schema ，循环属性，利用语法进行拼接
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
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
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
