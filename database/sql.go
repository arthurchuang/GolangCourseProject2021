package database

import (
	"GoCrawl/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// InitDB establishes a connection to PostgreSQL. It is required that the docker-compose image be brought up prior to the connection.
func InitDB(host string, dbname string, user string, password string) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// CreateTableIfNotExist creates a table for storing product entries in the given database.
func CreateTableIfNotExist(db *sql.DB, tableName string) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (id serial primary key, name VARCHAR(256), link VARCHAR(256), imgLink VARCHAR(256), price VARCHAR(256))")
	return err
}

// SaveProductEntry inserts the given ProductEntry to the given table in the specified database.
func SaveProductEntry(db *sql.DB, tableName string, pe *model.ProductEntry) error {
	sql := `
	INSERT INTO ` + tableName + ` (name, link, imgLink, price)
	VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sql, pe.Name, pe.Link, pe.ImgLink, pe.Price)
	return err
}

// TruncateTable truncates the specified table in the specified database.
func TruncateTable(db *sql.DB, tableName string) error {
	_, err := db.Exec("TRUNCATE " + tableName)
	return err
}

// GetElementCounts returns the number of elements in the specified table in the specified database.
func GetElementCounts(db *sql.DB, tableName string) (int, error) {
	var count int
	row := db.QueryRow(fmt.Sprintf("select count(1) from %s;", tableName))
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

//func QueryDB(db *sql.DB, sql string) error {
//	rows, err := db.Query(sql)
//	defer rows.Close()
//	for rows.Next() {
//		var name string
//		var link string
//		var imgLink string
//		var price string
//		err = rows.Scan(&name, &link, &imgLink, &price)
//		log.Printf(
//			"product name:\t "+name+"\n",
//			"product link:\t "+link+"\n",
//			"image link:\t "+imgLink+"\n",
//			"product price:\t"+price+"\n\n")
//	}
//	return err
//}
