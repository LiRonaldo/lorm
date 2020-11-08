package lorm

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"lorm/day5-hooks/dialect"
	"lorm/day5-hooks/log"
	"lorm/day5-hooks/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() (s *session.Session) {
	return session.New(e.db, e.dialect)
}
