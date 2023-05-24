package main

import (
	_ "github.com/go-sql-driver/mysql"
)

// func main() {
// 	db, err := sql.Open("mysql", "root:123456@tcp(172.17.0.2:3306)/ginchat")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	fmt.Println("aaaaaaaaaaaaaaaaa")

// 	rows, err := db.Query("SELECT * FROM user")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for rows.Next() {
// 		var id int
// 		var name string
// 		err = rows.Scan(&id, &name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("id: %d, name: %s\n", id, name)
// 	}
// }
