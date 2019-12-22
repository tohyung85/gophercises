package store

import (
	// "database/sql" - if not using sqlx
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

var dbSingleton *sqlx.DB
var gormDbSingleton *gorm.DB

type DbStyle int

const (
	SQL  DbStyle = 0
	SQLX DbStyle = 1
	GORM DbStyle = 2
)

var seedData [8]string = [8]string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

type PhoneStore struct {
}

type PhoneNo struct {
	Id     int    `gorm:"AUTO_INCREMENT"`
	Number string `gorm:"type:varchar(100);column:phone_no" db:"phone_no"`
}

func NewStore(db *sqlx.DB) (*PhoneStore, error) {
	dbSingleton = db
	phoneStore := &PhoneStore{}
	err := initializeStoreTables(SQL)
	if err != nil {
		return nil, err
	}
	return phoneStore, err
}

func NewGormStore(db *gorm.DB) (*PhoneStore, error) {
	gormDbSingleton = db
	phoneStore := &PhoneStore{}
	err := initializeStoreTables(GORM)
	if err != nil {
		return nil, err
	}
	return phoneStore, err
}

func initializeStoreTables(dbStyle DbStyle) error {
	var err error = nil
	switch dbStyle {
	case SQL:
		err = sqlInitTables()
	case GORM:
		err = gormInitTables()
	default:
		err = sqlInitTables()
	}
	return err
}

func (ps *PhoneStore) GetEntries(dbStyle DbStyle) ([]PhoneNo, error) {
	fmt.Println("Getting entries for db style: ", dbStyle)
	var err error = nil
	var result []PhoneNo
	switch dbStyle {
	case SQLX:
		result, err = sqlxGetEntries()
	case GORM:
		result, err = gormGetEntries()
	default:
		result, err = sqlGetEntries()
	}
	return result, err
}

func (ps *PhoneStore) DeleteEntry(id int, dbStyle DbStyle) error {
	var err error = nil
	switch dbStyle {
	case SQL:
		err = sqlDeleteEntry(id)
	case GORM:
		err = gormDeleteEntry(id)
	default:
		err = sqlDeleteEntry(id)
	}
	return err
}

func (ps *PhoneStore) UpdateEntry(id int, newVal string, dbStyle DbStyle) error {
	var err error = nil
	switch dbStyle {
	case SQL:
		err = sqlUpdateEntry(id, newVal)
	case GORM:
		err = gormUpdateEntry(id, newVal)
	default:
		err = sqlUpdateEntry(id, newVal)
	}
	return err
}
