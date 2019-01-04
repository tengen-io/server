package main

import (
    "database/sql"
	"fmt"
    "log"
    _ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello Go-Stop")

    connStr := "user=postgres password=postgres dbname=go_stop_go sslmode=verify-full"
    _, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Println("Connected to postgres")
    }
}
