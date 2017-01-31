package main

import (
	"github.com/trusch/jamesd2/db"
	"github.com/trusch/jamesd2/http"
)

func main() {
	db, _ := db.New("mongodb://localhost/test-db")
	http.ListenAndServe(db, ":8080")
}
