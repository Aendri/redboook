package mocks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) { //nolint:ireturn
	t.Helper()
	db, dbmock, err := sqlmock.New() //
	if err != nil {                  // mock sql.DB
		t.Errorf("sqlmock creation failure: %v", err)
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Errorf("Failed to open gorm v2 db, got error: %v", err)
	}
	return gdb, dbmock
}
