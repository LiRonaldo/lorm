package session

import (
	"database/sql"
	"lorm/day3-save-query/dialect"
	"lorm/day3-save-query/log"
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

func TestSession_Exec(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/lorm?charset=utf8")
	log.Error(err)
	TestDial, _ := dialect.GetDialect("mysql")
	s := New(db, TestDial).Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
}
func Test_Init(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../../gweb.db")
	TestDial, _ := dialect.GetDialect("sqlite3")
	s := New(db, TestDial).Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()

}
