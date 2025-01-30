package sqlite

import (
	"database/sql"
	"github.com/peterouob/golang_template/pkg/orm/session"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "./main.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func testRecordInit(t *testing.T) *session.Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err := s.DropTable()
	assert.NoError(t, err)
	err = s.CreateTable()
	assert.NoError(t, err)
	_, err = s.Insert(user1, user2)
	assert.NoError(t, err)
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), affected)
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}
