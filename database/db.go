package database

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB { // OOP Constructor
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	// Connect to SQL Server
	dsn := "sqlserver://amad:1234@localhost:1433?database=goDb"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error...")
		return nil
	}
	return db
}
