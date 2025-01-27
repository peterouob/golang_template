package engine

import (
	"database/sql"
	"fmt"
	"github.com/peterouob/golang_template/pkg/orm/dialect"
	"github.com/peterouob/golang_template/pkg/orm/session"
	"github.com/peterouob/golang_template/tools"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(drive, source string) *Engine {
	db, err := sql.Open(drive, source)
	tools.HandelError("open db error", err)
	err = db.Ping()
	tools.HandelError("ping db error", err)
	d, ok := dialect.GetDialect(drive)
	if !ok {
		tools.HandelError(fmt.Sprintf("get dialect error drive:%s", drive), err)
	}

	e := &Engine{db: db, dialect: d}
	tools.Log(fmt.Sprintf("open db success,source %s, drive %s", source, drive))
	return e
}

func (e *Engine) Close() {
	err := e.db.Close()
	tools.HandelError("close db error", err)
	tools.Log("close db success ...")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
