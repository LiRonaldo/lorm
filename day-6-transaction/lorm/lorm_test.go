package lorm

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Name string
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func TestEngine_NewSession(t *testing.T) {
	e, _ := NewEngine("mysql", "root:123456@tcp(localhost:3306)/lorm?charset=utf8")
	s := e.NewSession()
	s.Model(&User{})
	var users []User
	s.Limit(1).Find(&users)
	fmt.Println(users)
	s.Where("Name = ?", "liyuxiang1").Delete()
	s.Where("Name = ?", "duxiufeng").Update("Age = ", 50)
}

//定义一个tracsaction方法。当需要事务的时候，就调用这个方法，传入一个定义好的方法，在方法里对数据进行需要事务的操作。
//此方法，先进行begin 然后defer func 抓取错误，回滚，没有错误就提交
//如果不需要事务的方法，比如find ，直接就调用find方法就行。
func Test_transaction(t *testing.T) {
	e, _ := NewEngine("mysql", "root:123456@tcp(localhost:3306)/lorm?charset=utf8")
	s := e.NewSession()
	defer e.Close()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text , XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()

	e.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
