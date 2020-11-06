package lorm

import (
	"fmt"
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
}
