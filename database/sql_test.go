package database_test

import (
	"GoCrawl/database"
	"GoCrawl/model"
	"database/sql"
	"testing"
)

const (
	testDbUser     = "postgres"
	testDbPassword = "postgres"
	testDbName     = "postgres"
	testDbHost     = "localhost"
	testTableName  = "test_table"

	testProductName      = "any product name will do"
	testProductLink      = "any product link will do"
	testProductImageLink = "any image link will do"
	testProductPrice     = "any price will do"
)

func TestInitDB(t *testing.T) {
	_, err := database.InitDB(testDbHost, testDbName, testDbUser, testDbPassword)
	if err != nil {
		t.Log("Is your docker image up and running?")
		t.Errorf("Error while initiating connection to database: %s", err)
	}
}

func TestCreateTableIfNotExist(t *testing.T) {
	db, err := givenTestDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test database: %s", err)
	}
	if err = database.CreateTableIfNotExist(db, testTableName); err != nil {
		t.Error("Error while creating test_table: ", err)
	}
}

func TestSaveProductEntry(t *testing.T) {
	db, err := givenTestDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test database: %s", err)
	}

	pe := givenTestProductEntry()

	if err := database.SaveProductEntry(db, testTableName, pe); err != nil {
		t.Errorf("Failed to save product entry to database: %s", err)
	}

	got, err := database.GetElementCounts(db, testTableName)
	if err != nil {
		t.Fatalf("Failed to get number of elements from table %s : %s", testTableName, err)
	}
	if got != 1 {
		t.Errorf("Got %d element(s) in %s, want %d", got, testTableName, 1)
	}

	// cleanup
	if err = database.TruncateTable(db, testTableName); err != nil {
		t.Logf("Failed to clean up test table: %s", err)
	}
}

func TestTruncateTable(t *testing.T) {
	db, err := givenTestDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test database: %s", err)
	}

	pe := givenTestProductEntry()

	want := 20
	for i := 0; i < want; i++ {
		if err = database.SaveProductEntry(db, testTableName, pe); err != nil {
			t.Fatalf("Failed to save product entry to %s : %s", testTableName, err)
		}
	}

	got, err := database.GetElementCounts(db, testTableName)
	if err != nil {
		t.Fatalf("Failed to get number of elements from table %s : %s", testTableName, err)
	}
	if got != want {
		t.Errorf("Got %d element(s) in %s, want %d", got, testTableName, want)
	}

	if err = database.TruncateTable(db, testTableName); err != nil {
		t.Logf("Failed to truncate table: %s", err)
	}

	got, err = database.GetElementCounts(db, testTableName)
	if err != nil {
		t.Fatalf("Failed to get number of elements from table %s : %s", testTableName, err)
	}
	if got != 0 {
		t.Errorf("Still got %d elements after truncating table", got)
	}
}

func TestGetElementCounts(t *testing.T) {
	db, err := givenTestDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test database: %s", err)
	}

	got, err := database.GetElementCounts(db, testTableName)
	if err != nil {
		t.Fatalf("Failed to get number of elements from table %s : %s", testTableName, err)
	}
	if got != 0 {
		t.Errorf("Got %d element(s) in %s, want %d", got, testTableName, 0)
	}

	pe := givenTestProductEntry()

	want := 20
	for i := 0; i < want; i++ {
		if err = database.SaveProductEntry(db, testTableName, pe); err != nil {
			t.Fatalf("Failed to save product entry to %s : %s", testTableName, err)
		}
	}

	got, err = database.GetElementCounts(db, testTableName)
	if err != nil {
		t.Fatalf("Failed to get number of elements from table %s : %s", testTableName, err)
	}
	if got != want {
		t.Errorf("Got %d element(s) in %s, want %d", got, testTableName, want)
	}

	// cleanup
	if err = database.TruncateTable(db, testTableName); err != nil {
		t.Logf("Failed to clean up test table: %s", err)
	}
}

func givenTestProductEntry() *model.ProductEntry {
	return model.NewProductEntry(testProductName, testProductLink, testProductImageLink, testProductPrice)
}

func givenTestDatabase() (*sql.DB, error) {
	return database.InitDB(testDbHost, testDbName, testDbUser, testDbPassword)
}
