package serveGame

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init(){

	db,_:=sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/cow?charset=utf8")
	if err := db.Ping(); err != nil {
		fmt.Println("数据库连接失败")
		panic(err)
	}
	db.SetMaxIdleConns(2000)
	db.SetMaxOpenConns(1000)

	DB = db
	fmt.Println("数据库连接成功")
}


