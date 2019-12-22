package store

import (
	"database/sql"
	"fmt"
	"strings"
)

func sqlInitTables() error {
	err := sqlCreateTables()
	if err != nil {
		fmt.Printf("Error creating table: %s", err)
		return err
	}
	err = sqlSeedTable()
	if err != nil {
		fmt.Printf("Error seeding table: %s", err)
	}
	return err
}

func sqlCreateTables() error {
	queryString := "CREATE TABLE IF NOT EXISTS phone_numbers (id SERIAL PRIMARY KEY, phone_no varchar(255) NOT NULL)"
	result, err := dbSingleton.Exec(queryString)
	fmt.Printf("SQL create table result: %s\n", result)
	return err
}

func sqlSeedTable() error {
	truncString := "TRUNCATE phone_numbers"
	result, err := dbSingleton.Exec(truncString)
	if err != nil {
		fmt.Printf("Error truncating table\n")
		return err
	}
	insertString := ""
	numVals := make([]interface{}, 0)
	for id, num := range seedData {
		numVals = append(numVals, num)
		insertString += fmt.Sprintf(`($%d),`, id+1)
	}
	insertString = strings.TrimSuffix(insertString, ",") + ";"
	queryString := fmt.Sprintf(`INSERT INTO phone_numbers(phone_no) VALUES %s`, insertString)
	result, err = dbSingleton.Exec(queryString, numVals...)
	fmt.Printf("SQL seed table result: %s\n", result)
	return err
}

func sqlGetEntries() ([]PhoneNo, error) {
	queryString := "SELECT * FROM phone_numbers"
	rows, err := dbSingleton.Query(queryString)
	if err != nil {
		return nil, err
	}
	return sqlParseRows(rows), nil
}

func sqlDeleteEntry(id int) error {
	queryString := `DELETE FROM phone_numbers WHERE phone_numbers.id = $1`
	_, err := dbSingleton.Exec(queryString, id)
	// fmt.Printf("Deletion executed: %s", result)
	return err
}

func sqlUpdateEntry(id int, newVal string) error {
	queryString := `UPDATE phone_numbers SET phone_no=$1 WHERE phone_numbers.id = $2`
	_, err := dbSingleton.Exec(queryString, newVal, id)
	// fmt.Printf("Update executed: %s", result)
	return err
}

func sqlParseRows(rows *sql.Rows) []PhoneNo {
	results := make([]PhoneNo, 0)
	for rows.Next() {
		var id int
		var num string
		rows.Scan(&id, &num)
		results = append(results, PhoneNo{id, num})
	}
	return results
}
