package store

import (
	"fmt"
)

func gormInitTables() error {
	err := gormCreateTables()
	if err != nil {
		fmt.Printf("Error creating table: %s", err)
		return err
	}
	err = gormSeedTable()
	if err != nil {
		fmt.Printf("Error seeding table: %s", err)
	}
	return err
}

func gormCreateTables() error {
	gormDbSingleton.DropTableIfExists("phone_numbers")
	gormDbSingleton.Table("phone_numbers").CreateTable(&PhoneNo{})
	return nil
}

func gormSeedTable() error {
	for _, num := range seedData {
		gormDbSingleton.Table("phone_numbers").Create(&PhoneNo{Number: num})
	}
	return nil
}

func gormGetEntries() ([]PhoneNo, error) {
	var results []PhoneNo
	gormDbSingleton.Table("phone_numbers").Find(&results)
	return results, nil
}

func gormDeleteEntry(id int) error {
	var toDelete PhoneNo
	gormDbSingleton.Table("phone_numbers").First(&toDelete, id)
	gormDbSingleton.Table("phone_numbers").Delete(&toDelete)
	return nil
}

func gormUpdateEntry(id int, newVal string) error {
	var toUpdate PhoneNo
	gormDbSingleton.Table("phone_numbers").First(&toUpdate, id)
	toUpdate.Number = newVal
	gormDbSingleton.Table("phone_numbers").Save(&toUpdate)
	return nil
}
