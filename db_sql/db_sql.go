package db_sql

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)
const (
	USER = "postgres"
	PASSWORD= "postgres"
	DATABASE ="postgres"
	HOST ="localhost"
	port = 5432
	DB_NAME="carrefour"
)

type Database struct{
	DB *sql.DB
}
//function using example
func InitDB() *sql.DB{
	db, err := sql.Open(
		"postgres",fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",HOST, USER,PASSWORD,DATABASE))
	checkErr(err)
	if err = db.Ping(); err!=nil{
		panic(err)
	}
	return db
}
func function_format(){
	db, err := sql.Open(
		"postgres",fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",HOST, USER,PASSWORD,DATABASE))
	checkErr(err)
	if err = db.Ping(); err!=nil{
		panic(err)
	}
	fmt.Println("Successfully created connection to database")
	CreateTable(db,DB_NAME)
	InsertData(db, DB_NAME,"harry","link","imgLink","1500")
	QueryDB(db,`SELECT * FROM `+DB_NAME)
	DeleteAllData(db, DB_NAME)
	db.Close()
}

func checkErr(err error){
	if err != nil{
		panic(err)
	}
}

func CreatDatabase(db *sql.DB, database string){
	_, err := db.Exec("create database "+database)
	checkErr(err)
}

func CreateTable(db *sql.DB, table string){
	_,err := db.Exec("CREATE TABLE IF NOT EXISTS "+table+" (name VARCHAR(256), link VARCHAR(256), imgLink VARCHAR(256), price VARCHAR(256))")
	checkErr(err)
}
//still hardcore for now
func InsertData(db *sql.DB, table string, name string, link string, imgLink string, price string){
	sql := `
	INSERT INTO `+table+` (name, link, imgLink, price)
	VALUES ($1, $2, $3, $4)`
	_,err := db.Exec(sql,name,link,imgLink,price)
	checkErr(err)
}

func DeleteAllData(db *sql.DB, table string){
	_,err := db.Exec("TRUNCATE "+table)
	checkErr(err)
}

func QueryDB(db *sql.DB, sql string){
	rows, err := db.Query(sql)
	checkErr(err)
	defer rows.Close()
	for rows.Next(){
		var name string
		var link string
		var imgLink string
		var price string
		err = rows.Scan(&name,&link,&imgLink,&price)
		checkErr(err)
		fmt.Print(
			"product name:\t "+name+"\n",
			"product link:\t "+link+"\n",
			"image link:\t "+imgLink+"\n",
			"product price:\t"+price+"\n\n")
	}
}