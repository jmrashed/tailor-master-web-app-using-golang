package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Open a connection to the database
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/database_name")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Execute a query
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate over the rows and print the results
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(id, name)
	}
}
