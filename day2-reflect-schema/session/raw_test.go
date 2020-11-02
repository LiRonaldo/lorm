package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"lorm/day2-reflect-schema/dialect"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_Exec(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../../gweb.db")
	TestDial, _ := dialect.GetDialect("sqlite3")
	s := New(db, TestDial).Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()

}
