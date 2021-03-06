package session

import (
	"database/sql"
	"fmt"
	"lorm/day3-save-query/dialect"
	"testing"
)

type User struct {
	Name string
	Age  int
}

var (
	user1 = &User{"liyuxiang", 18}
	user2 = &User{"duxiufeng", 20}
	user3 = &User{"Jack", 25}
)

func (u *User) AfterQuery(s *Session) error {
	fmt.Println("1111111111111111")
	u.Name = "*******"
	return nil
}

func TestSession_Exec(t *testing.T) {
	db, _ := sql.Open("mysql", "root:123456@tcp(localhost:3306)/lorm?charset=utf8")
	TestDial, _ := dialect.GetDialect("mysql")
	_ = New(db, TestDial).Model(&User{})
	/*num, _ := s.Insert(user1, user2)
	fmt.Println(num)*/
	//var users []User
	//传切片地址，防止扩容。
	/*if err := s.Limit(1).Find(&users); err != nil {
		t.Fatal("failed to query all")
	}*/
	/*user := &User{}
	s.First(user)
	fmt.Println(user.Name)*/

}
func Test_Init(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../../gweb.db")
	TestDial, _ := dialect.GetDialect("sqlite3")
	s := New(db, TestDial).Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()

}
