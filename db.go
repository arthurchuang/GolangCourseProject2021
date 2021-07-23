package db
import (
	"fmt"
	//"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error){
	if err != nil{
		panic(err)
	}
}

var db *sql.DB

func Init_DB() error{
	db, err := sql.Open("mysql","root:123456@/harry")
	checkErr(err)
	//_ = &Env{d＿: db} //db *sql.DB
	fmt.Println("database setting")
	return db.Ping()
}

func Insert(db *sql.DB,name, link, imgLink, price string){
	stmt, err := db.Prepare("INSERT Carrefour_Table SET product_name=?, product_link=?, product_image=?,product_price=?")
	checkErr(err)
	fmt.Println(stmt)
	db.Close()
	// res, err := stmt.Exec(name,link,imgLink,price)
	// checkErr(err)
	// id, err := res.LastInsertId()
	// checkErr(err)
	// fmt.Println("id:",id)
	//db_Query(name,link,imgLink,price)
	//db.Close()
}

func db_Query(name, link, imgLink, price string){
	db, err := sql.Open("mysql","root:123456@/root?charset=utf8")
	checkErr(err)
	rows, err := db.Query("SELECT * FROM Carrefour_Table")
	checkErr(err)
	for rows.Next(){
		err = rows.Scan(&name, &link,&imgLink,&price)
		checkErr(err)
		fmt.Printf("product name: %s\n", name)
		fmt.Printf("product link: %s\n", link)
		fmt.Printf("product image: %s\n", imgLink)
		fmt.Printf("product price: %s", price)
		fmt.Printf("---------------------")
	}
	fmt.Println("Query Completed!")
	db.Close()

}