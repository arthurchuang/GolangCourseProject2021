package database_test

import (
	"GoCrawl/database"
	"GoCrawl/model"
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

func TestSaveProductEntry(t *testing.T) {
	pe := &model.ProductEntry{
		Name:    "any product name will do",
		Link:    "any product link will do",
		ImgLink: "any image link will do",
		Price:   "any price will do",
	}

	if err := database.SaveProductEntry(DB, "test_table", pe); err != nil {
		t.Errorf("Failed to save product entry to database: %s", err)
	}
}
