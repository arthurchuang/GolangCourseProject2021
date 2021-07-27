package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	USER     = "postgres"
	PASSWORD = "postgres"
	DATABASE = "postgres"
	HOST     = "localhost"
	port     = 5432
	DB_NAME  = "carrefour"
)

// Establish connection to PostgreSQL. It is required that the docker-compose image be brought up prior to the connection.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open(
		"postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE))
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}

func CreatDatabase(db *sql.DB, database string) error {
	_, err := db.Exec("create database " + database)
	return err
}

func CreateTable(db *sql.DB, table string) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + table + " (id serial not null primary key, name VARCHAR(256), link VARCHAR(256), imgLink VARCHAR(256), price VARCHAR(256))")
	return err
}

func InsertData(db *sql.DB, table string, name string, link string, imgLink string, price string) error {
	sql := `
	INSERT INTO ` + table + ` (name, link, imgLink, price)
	VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sql, name, link, imgLink, price)
	return err

	// To do: overwrite old entry with newer one in mysql table
}

func DeleteOldData(db *sql.DB, table string) error {
	_, err := db.Exec("TRUNCATE " + table)
	return err
}

func QueryDB(db *sql.DB, sql string) error {
	rows, err := db.Query(sql)
	defer rows.Close()
	for rows.Next() {
		var name string
		var link string
		var imgLink string
		var price string
		err = rows.Scan(&name, &link, &imgLink, &price)
		log.Printf(
			"product name:\t "+name+"\n",
			"product link:\t "+link+"\n",
			"image link:\t "+imgLink+"\n",
			"product price:\t"+price+"\n\n")
	}
	return err
}
