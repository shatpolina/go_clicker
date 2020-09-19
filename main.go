package main

import (
    "net/http"
    _ "github.com/lib/pq"
    "database/sql"
)

var db *sql.DB

func main() {
    connStr := "user=postgres password=postgres dbname=httphello sslmode=disable"
    local_db, err := sql.Open("postgres", connStr)
    db = local_db
    if err != nil {
        panic(err)
    } 
    defer db.Close()
    
    db.Exec("update users set Connected = $1", false)
    
    set_routes()
    
    http.ListenAndServe(":8080", nil)
}
