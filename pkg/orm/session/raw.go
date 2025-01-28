package session

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/peterouob/golang_template/pkg/orm/clause"
	"github.com/peterouob/golang_template/pkg/orm/dialect"
	"github.com/peterouob/golang_template/pkg/orm/schema"
	"github.com/peterouob/golang_template/tools"
	"reflect"
	"strings"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Reset() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}
func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, args ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, args...)
	return s
}

func (s *Session) Exec() (sql.Result, error) {
	defer s.Reset()
	tools.Log(fmt.Sprintf("%s:%v", s.sql.String(), s.sqlVars))
	result, err := s.db.Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		return nil, errors.New(fmt.Sprint("session exec err:", err))
	}
	return result, nil
}

func (s *Session) QueryRow() (*sql.Row, error) {
	defer s.Reset()
	tools.Log(fmt.Sprintf("%s:%v", s.sql.String(), s.sqlVars))
	return s.db.QueryRow(s.sql.String(), s.sqlVars...), nil
}

func (s *Session) QueryRows() *sql.Rows {
	defer s.Reset()
	tools.Log(fmt.Sprintf("%s:%v", s.sql.String(), s.sqlVars))
	result, err := s.db.Query(s.sql.String(), s.sqlVars...)
	tools.HandelError("query sql error", err)
	return result
}

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		tools.HandelError("", errors.New("model is not set"))
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	if _, err := s.Raw(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", table.Name, desc)).Exec(); err != nil {
		return errors.New(fmt.Sprint("create table err:", err))
	}
	return nil
}

func (s *Session) DropTable() error {
	if _, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec(); err != nil {
		return errors.New(fmt.Sprint("drop table err:", err))
	}
	return nil
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExists(s.RefTable().Name)
	row, err := s.Raw(sql, values...).QueryRow()
	tools.HandelError("query sql error", err)
	var tmp string
	err = row.Scan(&tmp)
	tools.HandelError("row scan error", err)
	return tmp == s.RefTable().Name
}
