package database

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

// github.com/denisenkom/go-mssqldb

//db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

func ConnectSqlDB(DB_SERVER string, DB_NAME string, DB_USER string, DB_PASS string) {
	var dsn = "sqlserver://" + DB_USER +":" + DB_PASS + "@" + DB_SERVER + "?database=" + DB_NAME + "&connection+timeout=30&encrypt=disable"
	var err error
	DBConn, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		//QueryFields: true,
	})
	//DBConn, err = gorm.Open("mssql", "sqlserver://IberoUser1:1B3r0*5tar@origos.no-ip.com:1433?database=OrigosVCIberostar")
	// sqlDB, err := DBConn.DB()
	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	// sqlDB.SetConnMaxLifetime(time.Hour)
	// sqlDB.Close()

	if err != nil {
		panic("Failed to connect sqldatabase")
	}
	fmt.Println("Open Connection in SQLDatabase")
	//database.DBConn.AutoMigrate(&book.Book{})
	//fmt.Println("Database Migrated")
}
