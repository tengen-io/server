/*
Server implementation of the board game Go.
*/
package main

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Server struct {
	db *sql.DB
}

func (s *Server) Start() {
	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
		return
	} else {
		log.Println("Connected to postgres")
	}

	s.db = db
	log.Println("Listening on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", s))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		w.Write([]byte("Hello"))
	}
}

func connectToDB() (*sql.DB, error) {
	connStr := "user=postgres password=postgres dbname=go_stop_go sslmode=verify-full"
	return sql.Open("postgres", connStr)
}

func main() {
	s := Server{}
	s.Start()
}
