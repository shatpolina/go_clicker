package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    _ "github.com/lib/pq"
    "database/sql"
)

type DB_inf struct {
    User string `json:"user"`
    Password string `json:"password"`
    DBname string `json:"dbname"`
    SSLmode string `json:"sslmode"`
}

type VK_id struct {
    ClientID string `json:"ClientID"`
    ClientSecret string `json"ClientSecret"`
}

var db *sql.DB

var obj_vk VK_id

func readJSON(file string, container interface{}) {
    rawDataIn, err := ioutil.ReadFile(file)
    if err != nil {
        fmt.Println("Cannot load " + file)
    }
    err = json.Unmarshal(rawDataIn, &container)
    if err != nil {
        fmt.Println("Invalid format " + file, err)
    }
}

func main() {
    var obj DB_inf
    readJSON("./data.json", &obj)
    connStr := "user=" + obj.User + " password=" + obj.Password + " dbname=" + obj.DBname + " sslmode=" + obj.SSLmode
    local_db, err := sql.Open("postgres", connStr)
    db = local_db
    if err != nil {
        panic(err)
    } 
    defer db.Close()
    

    readJSON("./vk_id.json", &obj_vk)
    fmt.Println(obj_vk)

    db.Exec("update users set Connected = $1", false)
    
    set_routes()
    
    http.ListenAndServe(":8080", nil)
}
