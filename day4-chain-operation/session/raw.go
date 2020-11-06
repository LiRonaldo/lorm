package session

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"lorm/day4-chain-operation/caluse"
	"lorm/day4-chain-operation/dialect"
	"lorm/day4-chain-operation/log"
	"lorm/day4-chain-operation/schema"
	"strings"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	sqlVars  []interface{}
	dialect  dialect.Dialect
	refTable *schema.Schema
	caluse   caluse.Caluse
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.caluse = caluse.Caluse{}
}

func (s *Session) DB() *sql.DB {
	return s.db
}
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

//执行,此处return 没有写返回具体的值，但是会根据Exec 方法的返回值，同名的去返回。
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

//查询一条记录
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

//多条记录
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
