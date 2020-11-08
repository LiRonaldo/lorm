package schema

import (
	"lorm/day5-hooks/dialect"
	"testing"
)

type User struct {
	Name string
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestSchema_GetField(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
