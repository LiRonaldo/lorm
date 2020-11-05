package session

import (
	"fmt"
	"lorm/day3-save-query/log"
	"lorm/day3-save-query/schema"
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
func (s *Session) Find(value interface{}) error {

}
