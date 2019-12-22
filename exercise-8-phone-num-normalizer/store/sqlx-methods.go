package store

import (
	"github.com/jmoiron/sqlx"
)

func sqlxGetEntries() ([]PhoneNo, error) {
	queryString := "SELECT * FROM phone_numbers"
	rows, err := dbSingleton.Queryx(queryString)
	if err != nil {
		return nil, err
	}
	return sqlxParseRows(rows), nil
}

func sqlxParseRows(rows *sqlx.Rows) []PhoneNo {
	results := make([]PhoneNo, 0)
	for rows.Next() {
		var pn PhoneNo
		err := rows.StructScan(&pn)
		if err != nil {
			continue
		}
		results = append(results, pn)
	}
	return results
}
