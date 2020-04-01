package main

import (
	"database/sql"
	"fmt"
	"github.com/IonutCraciun/databaseRestApi/databasehelper"
)

// Tag  a tag
type Tag struct {
	User string `json:"User"`
	Host string `json:"Host"`
}

func main() {
	fmt.Println("Hello, world.")
	fmt.Println(databasehelper.Test2)
	var connection *sql.DB
	connection = databasehelper.Connect("mysql", "root", "", "tcp(127.0.0.1:3306)")
	results, err := connection.Query("SELECT User, Host FROM user")

	if err != nil {
		fmt.Println(err.Error())
	}

	for results.Next() {
		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.User, &tag.Host)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		fmt.Println(tag.User)
		fmt.Println(tag.Host)
	}
}
