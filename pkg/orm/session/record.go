package session

import (
	"github.com/peterouob/golang_template/pkg/orm/clause"
	"github.com/peterouob/golang_template/tools"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldName)
		recordValues = append(recordValues, table.RecordValue(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	tools.HandelError("insert error ", err)
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem() //ELEM to for slice type
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldName)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows := s.Raw(sql, vars...).QueryRows()
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldName {
			// use addr to get pointer and reflect to interface
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		err := rows.Scan(values...)
		tools.HandelError("find error ", err)
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
