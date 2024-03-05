package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)


var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:Li258924.@tcp(localhost:33060)/fileserver?charset=utf8")

	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:", err)
		os.Exit(1)
	}
}

// DBConn : return mysql connection
func DBConn() *sql.DB {
	return db
}

