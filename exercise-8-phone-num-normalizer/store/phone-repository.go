package store

import (
	"database/sql"
)

var dbSingleton *sql.DB

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
	Id     int
	Number string
}

func NewStore(db *sql.DB) (*PhoneStore, error) {
	dbSingleton = db
	phoneStore := &PhoneStore{}
	err := initializeStoreTables(SQL)
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
	default:
		err = sqlInitTables()
	}
	return err
}

func parseRows(rows *sql.Rows) []PhoneNo {
	results := make([]PhoneNo, 0)
	for rows.Next() {
		var id int
		var num string
		rows.Scan(&id, &num)
		results = append(results, PhoneNo{id, num})
	}
	return results
}

func (ps *PhoneStore) GetEntries(dbStyle DbStyle) ([]PhoneNo, error) {
	var err error = nil
	var result []PhoneNo
	switch dbStyle {
	case SQL:
		result, err = sqlGetEntries()
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
	default:
		err = sqlUpdateEntry(id, newVal)
	}
	return err
}
