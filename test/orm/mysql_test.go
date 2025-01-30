package orm

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/peterouob/golang_template/pkg/orm/clause"
	"github.com/peterouob/golang_template/pkg/orm/dialect"
	"github.com/peterouob/golang_template/pkg/orm/session"
	"github.com/peterouob/golang_template/tools"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	tools.InitLogger()
}

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

const (
	INSERT clause.Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

type User struct {
	Name string `orm:"pk"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func initMysql() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/orm?charset=utf8")
	if err != nil {
		panic(errors.New("open mysql fail" + err.Error()))
	}
	return db
}

func NewSession(db *sql.DB) *session.Session {
	return session.New(db, TestDial)
}

func TestRecord(t *testing.T) {
	t.Helper()
	db := initMysql()
	assert.NotNil(t, db)
	s := NewSession(db)
	s.Model(&User{})
	assert.NotNil(t, s)
	err := s.CreateTable()
	assert.NoError(t, err)
	_, err = s.Insert(user1)
	assert.NoError(t, err)
	var u []User
	err = s.Find(&u)
	assert.NoError(t, err)
	assert.Equal(t, u[0].Name, "Tom")
}
