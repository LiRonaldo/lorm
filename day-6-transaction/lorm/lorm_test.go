package lorm

import (
	"errors"
	"fmt"
	"lorm/day-6-transaction/session"
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

func Test_transaction(t *testing.T) {
	e, _ := NewEngine("mysql", "root:123456@tcp(localhost:3306)/lorm?charset=utf8")
	e.Transaction(func(session *session.Session) (i interface{}, err error) {
		session.Insert(user2)
		return nil, errors.New("插入报错！")
	})
}
