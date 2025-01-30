package sqlite

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peterouob/golang_template/pkg/orm/clause"
	"github.com/peterouob/golang_template/pkg/orm/dialect"
	"github.com/peterouob/golang_template/pkg/orm/engine"
	"github.com/peterouob/golang_template/pkg/orm/schema"
	"github.com/peterouob/golang_template/pkg/orm/session"
	"github.com/peterouob/golang_template/tools"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func init() {
	tools.InitLogger()
}

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

const (
	INSERT clause.Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

func TestConnect(t *testing.T) {
	eg := engine.NewEngine("sqlite3", "./main.db")
	assert.NotEqual(t, nil, eg)
	defer eg.Close()
	s := eg.NewSession()
	_, err := s.Raw("DROP TABLE IF EXISTS User;").Exec()
	assert.NoError(t, err)
	_, err = s.Raw("CREATE TABLE User(Name text);").Exec()
	assert.NoError(t, err)
	_, err = s.Raw("CREATE TABLE User(Name text);").Exec()
	assert.Error(t, err)
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	t.Logf("Exec success, %d rows affected", count)
}

func NewSession() *session.Session {
	return session.New(TestDB, TestDial)
}

type User struct {
	Name string `orm:"primary key"`
	Age  int
}

func TestParse(t *testing.T) {
	sc := schema.Parse(&User{}, TestDial)
	assert.NotNil(t, sc)
	if sc.Name != "User" || len(sc.Fields) != 2 {
		t.Fatalf("parse error")
	}

	if sc.GetField("Name").Tag != "primary key" {
		t.Fatalf("parse error in tag")
	}
}

func TestSession_Create(t *testing.T) {
	s := NewSession().Model(&User{})
	err := s.DropTable()
	assert.NoError(t, err)
	err = s.CreateTable()
	assert.NoError(t, err)
	exist := s.HasTable()
	assert.True(t, exist)
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}

func testSelect(t *testing.T) {
	var clause1 clause.Clause
	clause1.Set(LIMIT, 3)
	clause1.Set(SELECT, "User", []string{"*"})
	clause1.Set(WHERE, "Name = ?", "Tom")
	clause1.Set(ORDERBY, "Age ASC")
	sql_, vars := clause1.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql_, vars)
	if sql_ != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}
