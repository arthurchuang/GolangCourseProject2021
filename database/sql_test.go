package database_test

import (
	"GoCrawl/database"
	"database/sql"
	"os"
	"testing"
)

var DB *sql.DB

func TestInitDB(t *testing.T) {
	db, err := database.InitDB()
	DB = db
	if err != nil {
		t.Error("Error while initiating connection to database: ", err)
		t.Log("Please check whether docker image is running")
		os.Exit(1)
	}
}

func TestCreateTable(t *testing.T) {
	if err := database.CreateTableIfNotExist(DB, "test_table"); err != nil {
		t.Error("Error while creating test_table: ", err)
	}
}

func TestInsertData(t *testing.T) {
	if err := database.InsertData(DB, "test_table", "test_name", "test_link", "test_imgLink", "test_price"); err != nil {
		t.Error("Error while inserting test data entry into table: ", err)
	}
}
