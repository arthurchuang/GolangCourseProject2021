package db_sql_test

import (
	"GoCrawl/db_sql"
	"database/sql"
	"os"
	"testing"
)

var DB *sql.DB

func TestInitDB(t *testing.T) {
	db, err := db_sql.InitDB()
	DB = db
	if err != nil {
		t.Error("Error while initiating connection to database: ", err)
		t.Log("Please check whether docker image is running")
		os.Exit(1)
	}
}

func TestCreateTable(t *testing.T) {
	if err := db_sql.CreateTable(DB, "test_table"); err != nil {
		t.Error("Error while creating test_table: ", err)
	}
}

func TestInsertData(t *testing.T) {
	if err := db_sql.InsertData(DB, "test_table", "test_name", "test_link", "test_imgLink", "test_price"); err != nil {
		t.Error("Error while inserting test data entry into table: ", err)
	}
}
